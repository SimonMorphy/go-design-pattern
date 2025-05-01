package user

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/SimonMorphy/go-design-pattern/internal/common/const/errors"
	"github.com/SimonMorphy/go-design-pattern/internal/infrastructure/creational"
	"github.com/goccy/go-json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"time"
)

var tokenSecret = viper.GetString("token-secret")

type Usr struct {
	ID        uint           `gorm:"primarykey"`
	CreatedAt *time.Time     `json:"created_at,omitempty" gorm:"autoCreateTime"`
	UpdatedAt *time.Time     `json:"updated_at,omitempty" gorm:"autoUpdateTime"`
	Username  string         `json:"username" validate:"required,min=3,max=20"`
	Password  string         `json:"password" validate:"required,min=6,max=32"`
	Email     string         `json:"email" validate:"required,email"`
	Mobile    string         `json:"mobile" validate:"omitempty,e164"`
	Address   string         `json:"address" validate:"omitempty,max=255"`
	Token     string         `json:"token" validate:"omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func (u *Usr) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, u)
}

func (u *Usr) MarshalBinary() (data []byte, err error) {
	return json.Marshal(u)
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
	jsonStr, err := json.Marshal(u)
	if err != nil {
		logrus.Error(err)
		return ""
	}
	return string(jsonStr)
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
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)), // 设置 token 有效期为 24 小时
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	// 创建一个新的 token 对象，指定签名方法为 HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// 使用指定的密钥对 token 进行签名
	tokenStr, err := token.SignedString([]byte(tokenSecret))
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

			return nil, errors.New(errors.ErrnoUserTokenInvalid)
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

func NewUser(id uint, username string, password string, email string, mobile string, address string) (usr *Usr, err error) {
	user := &Usr{ID: id, Username: username, Password: password, Email: email, Mobile: mobile, Address: address}
	user.generateToken()
	user.encryptPassword()
	return user, nil
}

type NotFountError struct {
	Id uint
}

func (n NotFountError) Error() string {
	return fmt.Sprintf("User %d Not Found", n.Id)
}
