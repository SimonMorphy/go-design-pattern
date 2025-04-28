package user

import (
	"context"
	"github.com/SimonMorphy/go-design-pattern/internal/common/config/models"
	domain "github.com/SimonMorphy/go-design-pattern/internal/domain/user"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MysqlUserRepository struct {
	DB *gorm.DB
}

func NewMysqlUserRepository() MysqlUserRepository {
	db, err := models.GetDB()
	if err != nil {
		logrus.Panic(err)
	}
	return MysqlUserRepository{
		DB: db,
	}
}
func (m MysqlUserRepository) List(ctx context.Context, off, lim int) ([]*domain.Usr, error) {
	var users []*domain.Usr
	tx := m.DB.WithContext(ctx).Offset(off).Limit(lim).Find(&users)
	if tx.Error != nil {
		logrus.Error(tx.Error)
		return nil, tx.Error
	}
	return users, nil
}

func (m MysqlUserRepository) Create(ctx context.Context, user *domain.Usr) (uint, error) {
	tx := m.DB.WithContext(ctx).Create(user)
	if tx.Error != nil {
		return -1, tx.Error
	}
	return user.ID, nil
}

func (m MysqlUserRepository) Get(ctx context.Context, ID uint) (*domain.Usr, error) {
	usr := domain.Usr{}
	tx := m.DB.WithContext(ctx).First(&usr, ID)
	if tx.Error == nil {
		return &usr, nil
	}
	return nil, tx.Error
}

func (m MysqlUserRepository) Update(ctx context.Context, usr *domain.Usr, fun func(context.Context, *domain.Usr) (*domain.Usr, error)) error {
	newUsr, err := fun(ctx, usr)
	if err != nil {
		logrus.Error(err)
		return err
	}
	tx := m.DB.WithContext(ctx).Updates(newUsr)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
