package entity

import (
	"time"

	"gitlab.com/stoqu/stoqu-be/internal/config"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserEntity struct {
	Name         string `json:"name" gorm:"size:50;not null"`
	Email        string `json:"email" gorm:"index:idx_user_email;unique;size:150;not null"`
	PasswordHash string `json:"-"`
	Password     string `json:"password" gorm:"-"`

	// fk
	RoleID string `json:"role_id" gorm:"not null"`
}

type UserModel struct {
	Entity
	UserEntity

	// relations
	UserProfile *UserProfileModel `json:"-" gorm:"foreignKey:UserID"`
}

func (UserModel) TableName() string {
	return "users"
}

func (m *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	err = m.Entity.BeforeCreate(tx)
	if err != nil {
		return
	}

	m.HashPassword()
	m.Password = ""
	return
}

func (m *UserModel) HashPassword() {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(m.Password), bcrypt.DefaultCost)
	m.PasswordHash = string(bytes)
}

func (m *UserModel) GenerateToken() (string, error) {
	var (
		jwtKey = config.Config.JWT.Secret
	)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    m.ID,
		"email": m.Email,
		"name":  m.Name,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString([]byte(jwtKey))
	return tokenString, err
}
