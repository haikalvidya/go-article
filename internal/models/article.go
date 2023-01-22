package models

import (
	"time"

	"github.com/haikalvidya/go-article/internal/delivery/payload"
	"gorm.io/gorm"
)

type ArticleModel struct {
	ID        int            `db:"id"`
	Title     string         `db:"title"`
	Body      string         `db:"body"`
	AuthorID  string         `db:"author_id"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt *time.Time     `db:"updated_at"`
	DeletedAt gorm.DeletedAt `db:"deleted_at"`

	Author *UserModel `gorm:"foreignKey:AuthorID"`
}

func (ArticleModel) TableName() string {
	return "articles"
}

func (a *ArticleModel) BeforeCreate(tx *gorm.DB) (err error) {
	a.CreatedAt = time.Now()
	return
}

func (a *ArticleModel) PublicInfo() *payload.ArticleInfo {
	res := &payload.ArticleInfo{
		ID:        a.ID,
		Title:     a.Title,
		Content:   a.Body,
		AuthorID:  a.AuthorID,
		CreatedAt: a.CreatedAt.Format(time.RFC3339),
		UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
	}

	if a.Author != nil {
		res.Author = a.Author.PublicInfo()
	}

	return res
}
