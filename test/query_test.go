package test

import (
	"articles/query/app"
	"articles/query/broker"
	"articles/query/cache"
	"articles/query/controller"
	"articles/query/domain"
	"articles/query/repository"
	"articles/query/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupRouter() http.Handler {
	ch := app.InitBroker()
	es := app.InitESClient()
	redisClient := app.InitCache()

	articleChan := make(chan []byte)
	cacheChan := make(chan domain.Article)

	msgBroker := broker.NewConsumer()
	insert := repository.Newinsert()
	redis := cache.NewCache()
	query := repository.NewQuery()
	services := service.NewServiceQuery(query, es, redis, redisClient, cacheChan)
	controllers := controller.NewController(services)

	go msgBroker.Consum(ch, articleChan)
	go insert.Save(es, articleChan, cacheChan)
	go redis.Set(redisClient, cacheChan)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.GET("/articles", controllers.SearchArticle)

	return router
}

func TestGetEmptyParamsSuccess(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/articles", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "success", keyValue["status"])
	assert.NotNil(t, keyValue["data"])
}

func TestGetQueryParamsSuccess(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/articles?query=test", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "success", keyValue["status"])
	assert.NotNil(t, keyValue["data"])
}

func TestGetQueryParamsFailed(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/articles?query=itShouldBeNothingHere", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 404, result.StatusCode)
	assert.Equal(t, "not found", keyValue["status"])
	assert.Nil(t, keyValue["data"])
}

func TestGetAuthorParamsSuccess(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/articles?author=test", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "success", keyValue["status"])
	assert.NotNil(t, keyValue["data"])
}

func TestGetAuthorParamsFailed(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/articles?author=itShouldBeNothingHere", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 404, result.StatusCode)
	assert.Equal(t, "not found", keyValue["status"])
	assert.Nil(t, keyValue["data"])
}

func TestGetAuthorAndQueryParamsSuccess(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/articles?author=test&query=test", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "success", keyValue["status"])
	assert.NotNil(t, keyValue["data"])
}

func TestGetAuthorAndQueryParamsFailed(t *testing.T) {
	router := setupRouter()

	request := httptest.NewRequest(http.MethodGet, "http://localhost:3001/articles?author=itShouldBeNothingHere&query=itShouldBeNothingHere", nil)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 404, result.StatusCode)
	assert.Equal(t, "not found", keyValue["status"])
	assert.Nil(t, keyValue["data"])
}
