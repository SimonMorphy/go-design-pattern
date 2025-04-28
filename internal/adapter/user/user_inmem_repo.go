package user

import (
	"context"
	"fmt"
	domain "github.com/SimonMorphy/go-design-pattern/internal/domain/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"sync"
)

type MemoryUserRepository struct {
	lock *sync.RWMutex
	data []*domain.Usr
}

func NewMemoryUserRepository() *MemoryUserRepository {
	return &MemoryUserRepository{lock: &sync.RWMutex{}, data: make([]*domain.Usr, 0)}
}

func (m *MemoryUserRepository) List(ctx context.Context, off, lim int) ([]*domain.Usr, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	if off < 0 || lim < 0 {
		return nil, fmt.Errorf("invalid pagination parameters: offset and limit must be non-negative")
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
func (m *MemoryUserRepository) Create(ctx context.Context, user *domain.Usr) (uint, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	newUUID, _ := uuid.NewUUID()
	usr, err := domain.NewUser(user.Username, user.Password, user.Email, user.Mobile, user.Address)
	if err != nil {
		logrus.Error(err)
		return -1, err
	}
	usr.ID = uint(newUUID.ID())
	m.data = append(m.data, usr)
	logrus.WithContext(ctx).WithFields(logrus.Fields{
		"insert_user":       usr,
		"data_after_insert": m.data,
	}).Debug("user_memo_repo_create")
	return usr.ID, nil
}

func (m *MemoryUserRepository) Get(_ context.Context, ID uint) (*domain.Usr, error) {
	m.lock.RLock()
	defer m.lock.RUnlock()
	for _, usr := range m.data {
		if usr.ID == ID {
			return usr, nil
		}
	}
	logrus.Debug("user %d Not Found", ID)
	return nil, domain.NotFountError{Id: ID}
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
		return domain.NotFountError{Id: usr.ID}
	}
	return nil
}
