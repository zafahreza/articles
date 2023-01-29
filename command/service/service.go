package service

import (
	"articles/command/broker"
	"articles/command/domain"
	"articles/command/helper"
	"articles/command/repository"
	"context"
	"github.com/streadway/amqp"
	"gorm.io/gorm"
)

type SaveService interface {
	CreateArticle(ctx context.Context, req domain.RequestArticle) (domain.Article, error)
}

type Service struct {
	Repo   repository.Saver
	Broker broker.Publisher
	Db     *gorm.DB
	Ch     *amqp.Channel
}

func NewService(repo repository.Saver, broker broker.Publisher, db *gorm.DB, ch *amqp.Channel) SaveService {
	return &Service{
		Repo:   repo,
		Broker: broker,
		Db:     db,
		Ch:     ch,
	}
}

func (s *Service) CreateArticle(ctx context.Context, req domain.RequestArticle) (domain.Article, error) {

	tx := s.Db.Begin()
	defer helper.CommitOrRollback(tx)

	article, err := s.Repo.Save(ctx, tx, req)
	if err != nil {
		return article, err
	}

	err = s.Broker.Publish(s.Ch, article)
	if err != nil {
		return article, err
	}

	return article, nil
}
