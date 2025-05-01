package external

import (
	"context"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"github.com/google/uuid"
	"github.com/spf13/viper"
	"io"
	"os"
	"sync"

	"github.com/goccy/go-json"
	"github.com/sirupsen/logrus"

	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
)

type JsonUserRepository struct {
	lock *sync.RWMutex
	file string
}

func NewJsonUserRepository() (*JsonUserRepository, func()) {
	return &JsonUserRepository{file: viper.GetString("json.file"), lock: &sync.RWMutex{}}, func() {}
}

func (j JsonUserRepository) Update(ctx context.Context, usr *domain.Usr, fun func(context.Context, *domain.Usr) (*domain.Usr, error)) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	j.lock.Lock()
	defer j.lock.Unlock()

	data, err := j.readAllData()
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return err
	}

	_usr, err := fun(ctx, usr)
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return err
	}

	found := false
	for idx, u := range data {
		if u.ID == usr.ID {
			data[idx] = _usr
			found = true
			break
		}
	}
	if !found {
		return errors.New(errors.ErrnoUserNotFoundError)
	}

	err = j.writeAllData(data)
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return err
	}
	return nil
}

func (j JsonUserRepository) readAllData() ([]*domain.Usr, error) {
	file, err := os.OpenFile(j.file, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var result []*domain.Usr
	if err = json.NewDecoder(file).Decode(&result); err != nil {
		return nil, err
	}
	return result, nil
}

func (j JsonUserRepository) writeAllData(users []*domain.Usr) error {
	data, err := json.Marshal(users)
	if err != nil {
		return err
	}
	tempFile := j.file + ".tmp"
	if err = os.WriteFile(tempFile, data, 0644); err != nil {
		return err
	}
	return os.Rename(tempFile, j.file)
}

func (j JsonUserRepository) Create(ctx context.Context, user *domain.Usr) (uint, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}

	j.lock.Lock()
	defer j.lock.Unlock()
	user.ID = uint(uuid.New().ID())
	users, err := j.readAllData()
	if err != nil && err != io.EOF {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return 0, err
	}

	users = append(users, user)
	err = j.writeAllData(users)
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return 0, err
	}
	return user.ID, nil
}

func (j JsonUserRepository) List(_ context.Context, offset, limit int) ([]*domain.Usr, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()
	offset--
	users, err := j.readAllData()
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return nil, err
	}

	start := offset * limit
	if start >= len(users) {
		return []*domain.Usr{}, nil
	}

	end := start + limit
	if end > len(users) {
		end = len(users)
	}

	return users[start:end], nil
}

func (j JsonUserRepository) Get(_ context.Context, id uint) (*domain.Usr, error) {
	j.lock.RLock()
	defer j.lock.RUnlock()

	data, err := j.readAllData()
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return nil, err
	}

	for _, usr := range data {
		if usr.ID == id {
			return usr, nil
		}
	}
	return nil, errors.New(errors.ErrnoUserNotFoundError)
}

func (j JsonUserRepository) Delete(ctx context.Context, id uint) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	j.lock.Lock()
	defer j.lock.Unlock()

	data, err := j.readAllData()
	if err != nil {
		logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return err
	}

	for idx, usr := range data {
		if usr.ID == id {
			data = append(data[:idx], data[idx+1:]...)
			err = j.writeAllData(data)
			if err != nil {
				logrus.Error(errors.NewWithError(errors.ErrnoInternalServerError, err))
				return err
			}
			return nil
		}
	}
	return errors.New(errors.ErrnoUserNotFoundError)
}
