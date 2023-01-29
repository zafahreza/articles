package service

import (
	"articles/query/cache"
	"articles/query/domain"
	"articles/query/repository"
	"github.com/elastic/go-elasticsearch"
	"github.com/redis/go-redis/v9"
)

type ServiceQuery interface {
	Search(params domain.Params) ([]domain.Article, error)
}

type Service struct {
	Repo      repository.Query
	EsClient  *elasticsearch.Client
	Cache     cache.Cache
	RedClient *redis.Client
	CacheChan chan domain.Article
}

func NewServiceQuery(repo repository.Query, es *elasticsearch.Client, cache2 cache.Cache, client *redis.Client, cacheChan chan domain.Article) ServiceQuery {
	return &Service{
		Repo:      repo,
		EsClient:  es,
		Cache:     cache2,
		RedClient: client,
		CacheChan: cacheChan,
	}
}

func (s *Service) Search(params domain.Params) ([]domain.Article, error) {
	if params.Author == "" && params.Query == "" {
		results, err := s.Cache.Get(s.RedClient)
		if err != nil {
			results, err = s.Repo.Search(s.EsClient, params)
			if err != nil {
				return nil, err
			}
			for _, result := range results {
				s.CacheChan <- result
			}
		}
		return results, nil
	}

	results, err := s.Repo.Search(s.EsClient, params)
	if err != nil {
		return nil, err
	}
	return results, nil
}
