package user

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	errors "github.com/SimonMorphy/go-design-pattern/internal/common/const"
	"github.com/SimonMorphy/go-design-pattern/internal/common/infrastructure/creational"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

var tokenSecret = viper.GetString("token-secret")

type Usr struct {
	gorm.Model
	Username string `json:"username" validate:"required,min=3,max=20"`
	Password string `json:"password" validate:"required,min=6,max=32"`
	Email    string `json:"email" validate:"required,email"`
	Mobile   string `json:"mobile" validate:"omitempty,e164"`
	Address  string `json:"address" validate:"omitempty,max=255"`
	Token    string `json:"token" validate:"omitempty"`
}

func (u *Usr) Clone() creational.Cloneable {
	jsonStr, err := json.Marshal(u)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	var usr Usr
	err = json.Unmarshal(jsonStr, &usr)
	if err != nil {
		logrus.Error(err)
		return nil
	}
	return &usr
}

func (u *Usr) Print() string {
	//TODO implement me
	panic("implement me")
}

func (u *Usr) TableName() string {
	return "usr"
}

type UsrClaims struct {
	UserId   uint   `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (u *Usr) encryptPassword() {
	hash := md5.New()
	hash.Write([]byte(u.Password))
	hashedBytes := hash.Sum(nil)
	u.Password = hex.EncodeToString(hashedBytes)
}

func (u *Usr) generateToken() {
	claims := UsrClaims{
		UserId:   u.ID,
		Username: u.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), // 设置 token 有效期为 24 小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// 创建一个新的 token 对象，指定签名方法为 HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的密钥对 token 进行签名
	tokenStr, err := token.SignedString(tokenSecret)
	if err != nil {
		logrus.Error(err)
		return
	}
	u.Token = tokenStr
}

func (u *Usr) IsExpired() error {
	token, err := jwt.ParseWithClaims(u.Token, &UsrClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

			return nil, errors.NewWithMsgf(603, "unexpected signing method: %v", token.Header["alg"])
		}
		return tokenSecret, nil
	})
	if err != nil {
		return err
	}
	// 验证 token 是否有效
	if _, ok := token.Claims.(*UsrClaims); ok && token.Valid {
		return nil
	}
	return errors.New(errors.ErrnoUserTokenInvalid)
}

func NewUser(username string, password string, email string, mobile string, address string) (usr *Usr, err error) {
	user := &Usr{Username: username, Password: password, Email: email, Mobile: mobile, Address: address}
	user.generateToken()
	user.encryptPassword()
	va := validator.New()
	err = va.Struct(usr)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	return user, nil
}

type NotFountError struct {
	Id uint
}

func (n NotFountError) Error() string {
	return fmt.Sprintf("User %d Not Found", n.Id)
}
