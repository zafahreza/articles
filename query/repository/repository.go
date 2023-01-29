package repository

import (
	"articles/query/domain"
	"articles/query/helper"
	"bytes"
	"context"
	"encoding/json"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"strconv"
	"time"
)

type Insert interface {
	Save(es *elasticsearch.Client, articleChan chan []byte, cacheChan chan domain.Article)
}

type Query interface {
	Search(es *elasticsearch.Client, params domain.Params) ([]domain.Article, error)
}

type Repository struct {
}

func Newinsert() Insert {
	return &Repository{}
}

func NewQuery() Query {
	return &Repository{}
}

func (r Repository) Save(es *elasticsearch.Client, articleChan chan []byte, cacheChan chan domain.Article) {
	article := domain.Article{}
	for {
		articleJson := <-articleChan
		err := json.Unmarshal(articleJson, &article)
		if err != nil {
			panic(err)
		}

		req := esapi.IndexRequest{
			Index:      "article",
			DocumentID: strconv.Itoa(article.Id),
			Body:       bytes.NewReader(articleJson),
			Refresh:    "true",
		}

		_, err = req.Do(context.Background(), es)
		if err != nil {
			panic(err)
		}
		cacheChan <- article
	}

}

func (r Repository) Search(es *elasticsearch.Client, params domain.Params) ([]domain.Article, error) {
	var buff bytes.Buffer
	articles := []domain.Article{}

	query := helper.WhichQuery(params)

	if err := json.NewEncoder(&buff).Encode(query); err != nil {
		return nil, err
	}

	res, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("article"),
		es.Search.WithBody(&buff),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var result map[string]interface{}

	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
		return nil, err
	}

	if result["hits"] == nil {
		return nil, err
	}

	for _, hit := range result["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		articleId, err := strconv.Atoi(hit.(map[string]interface{})["_id"].(string))
		if err != nil {
			return articles, err
		}
		date, err := time.Parse(time.RFC3339Nano, source.(map[string]interface{})["created_at"].(string))
		if err != nil {
			return articles, err
		}
		article := domain.Article{
			Id:        articleId,
			Author:    source.(map[string]interface{})["author"].(string),
			Title:     source.(map[string]interface{})["title"].(string),
			Body:      source.(map[string]interface{})["body"].(string),
			CreatedAt: date,
		}

		articles = append(articles, article)

	}

	return articles, nil

}
