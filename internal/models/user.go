package models

import (
	"time"

	"github.com/haikalvidya/go-article/internal/delivery/payload"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserModel struct {
	ID        string         `db:"id"`
	Email     string         `db:"email"`
	Name      string         `db:"name"`
	Password  string         `db:"password"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt *time.Time     `db:"updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at"`
}

// create before create gorm for adding uuid to id and created_at time
func (u *UserModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New().String()
	u.CreatedAt = time.Now()
	return
}

func (u *UserModel) PublicInfo() *payload.UserInfo {
	userpayload := &payload.UserInfo{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}

	return userpayload
}

func (UserModel) TableName() string {
	return "users"
}
