package repository

import (
	"articles/command/domain"
	"context"
	"gorm.io/gorm"
)

type Saver interface {
	Save(ctx context.Context, tx *gorm.DB, article domain.RequestArticle) (domain.Article, error)
}

type Repository struct {
}

func NewRepository() Saver {
	return &Repository{}
}

func (r Repository) Save(ctx context.Context, tx *gorm.DB, article domain.RequestArticle) (domain.Article, error) {
	result := domain.Article{
		Title:  article.Title,
		Author: article.Author,
		Body:   article.Body,
	}

	err := tx.WithContext(ctx).Create(&result).Error
	if err != nil {
		return result, err
	}

	return result, nil
}
