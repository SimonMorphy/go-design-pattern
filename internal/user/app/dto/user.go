package dto

import (
	domain "github.com/SimonMorphy/go-design-pattern/internal/user/domain/user"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Usr struct {
	ID       uint   `json:"id" validate:"omitempty,min=0"`
	Username string `json:"username" validate:"omitempty,min=3,max=20"`
	Password string `json:"password" validate:"omitempty,min=6,max=32"`
	Email    string `json:"email" validate:"omitempty,email"`
	Mobile   string `json:"mobile" validate:"omitempty,e164"`
	Address  string `json:"address" validate:"omitempty,max=255"`
}

func (u *Usr) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

func (u *Usr) ToDomain() *domain.Usr {
	err := u.Validate()
	if err != nil {
		logrus.Error(err)
		return nil
	}
	usr, err := domain.NewUser(u.ID, u.Username, u.Password, u.Email, u.Mobile, u.Address)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return usr
}
