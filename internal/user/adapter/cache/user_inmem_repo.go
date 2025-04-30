package cache

import (
	"context"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/infrastructure/creational"
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

var memoryRepositoryFactory = creational.NewSingletonFactory()

type MemoryUserRepository struct {
	lock *sync.RWMutex
	data []*domain.Usr
}

func (m *MemoryUserRepository) Delete(ctx context.Context, ID uint) error {
	//TODO implement me
	panic("implement me")
}

func MemoryRepositorySupplier() (interface{}, error) {
	return &MemoryUserRepository{lock: &sync.RWMutex{}, data: make([]*domain.Usr, 0)}, nil
}

func GetMemoryRepository() (*MemoryUserRepository, error) {
	repo, err := memoryRepositoryFactory.Get("memory")
	if err != nil {
		logrus.Panic(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return nil, err
	}
	repository, ok := repo.(*MemoryUserRepository)
	if !ok {
		logrus.Panic(errors.New(errors.ErrnoCastError))
	}
	return repository, nil
}

func NewMemoryUserRepository() (*MemoryUserRepository, func()) {
	memoryRepositoryFactory.Register("memory", MemoryRepositorySupplier)
	repository, err := GetMemoryRepository()
	if err != nil {
		logrus.Panic(errors.NewWithError(errors.ErrnoInternalServerError, err))
		return nil, nil
	}
	return repository, func() {
		memoryRepositoryFactory.Clear("memory")
	}
}

func (m *MemoryUserRepository) List(_ context.Context, off, lim int) ([]*domain.Usr, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if off < 0 || lim < 0 {
		return nil, errors.New(errors.ErrnoParameterInputError)
	}
	if off >= len(m.data) {
		return []*domain.Usr{}, nil
	}
	end := off + lim
	if end > len(m.data) {
		end = len(m.data)
	}
	return m.data[off:end], nil
}
func (m *MemoryUserRepository) Create(ctx context.Context, usr *domain.Usr) (uint, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	newUUID, _ := uuid.NewUUID()
	_usr, err := domain.NewUser(uint(newUUID.ID()), usr.Username, usr.Password, usr.Email, usr.Mobile, usr.Address)
	if err != nil {
		logrus.Error(err)
		return 0, err
	}
	m.data = append(m.data, _usr)
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"insert_user":       _usr,
		"data_after_insert": m.data,
	}).Debug("user_memo_repo_create")
	return _usr.ID, nil
}

func (m *MemoryUserRepository) Get(_ context.Context, ID uint) (*domain.Usr, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, usr := range m.data {
		if usr.ID == ID {
			return usr, nil
		}
	}
	logrus.Debugf("repository %d Not Found", ID)
	return nil, errors.NewWithError(errors.ErrnoUserNotFoundError, domain.NotFountError{Id: ID})
}

func (m *MemoryUserRepository) Update(ctx context.Context, usr *domain.Usr, fun func(context.Context, *domain.Usr) (*domain.Usr, error)) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	newUsr, err := fun(ctx, usr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	var found bool
	for idx, user := range m.data {
		if user.ID == newUsr.ID {
			m.data[idx] = newUsr
			found = true
		}
	}
	if !found {
		return errors.NewWithError(errors.ErrnoUserNotFoundError, domain.NotFountError{Id: usr.ID})
	}
	return nil
}
