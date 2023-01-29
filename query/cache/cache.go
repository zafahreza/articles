package cache

import (
	"articles/query/domain"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"sort"
)

type Cache interface {
	Get(client *redis.Client) ([]domain.Article, error)
	Set(client *redis.Client, cacheChan chan domain.Article)
}

type Redis struct {
}

func NewCache() Cache {
	return &Redis{}
}

func (r Redis) Set(client *redis.Client, cacheChan chan domain.Article) {
	for {
		data := <-cacheChan
		articleByte, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}
		key := fmt.Sprintf("article%d", data.Id)
		err = client.Set(context.Background(), key, articleByte, 0).Err()
		if err != nil {
			panic(err)
		}
		fmt.Println("set data to redis")
	}
}

func (r Redis) Get(client *redis.Client) ([]domain.Article, error) {
	article := domain.Article{}
	articles := []domain.Article{}
	keys, err := client.Keys(context.Background(), "article*").Result()
	if err != nil {
		return nil, err
	}
	results, err := client.MGet(context.Background(), keys...).Result()
	if err != nil {
		return nil, err
	}
	for _, result := range results {
		err := json.Unmarshal([]byte(result.(string)), &article)
		if err != nil {
			return nil, err
		}
		articles = append(articles, article)
	}

	sort.Slice(articles[:], func(i, j int) bool {
		return articles[i].Id > articles[j].Id
	})

	if len(articles) == 0 {
		return nil, errors.New("not found")
	}

	fmt.Println("get data from redis")
	return articles, nil
}
