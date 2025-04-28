package dto

import (
	domain "github.com/SimonMorphy/go-design-pattern/internal/domain/user"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type Usr struct {
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Mobile   string `json:"mobile" validate:"omitempty,e164"`
	Address  string `json:"address" validate:"omitempty,max=255"`
}

func (u *Usr) Validate() error {
	v := validator.New()
	return v.Struct(u)
}

func (u *Usr) ToDomain() *domain.Usr {
	usr, err := domain.NewUser(u.Username, u.Password, u.Email, u.Mobile, u.Address)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return usr
}
