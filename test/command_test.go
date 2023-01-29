package test

import (
	"articles/command/app"
	"articles/command/broker"
	"articles/command/controller"
	"articles/command/repository"
	"articles/command/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func setupRouterCommand() http.Handler {
	db := app.InitDB()
	ch := app.InitBroker()

	repo := repository.NewRepository()
	messageBroker := broker.NewBroker()
	services := service.NewService(repo, messageBroker, db, ch)
	controlService := controller.NewController(services)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()

	router.POST("/articles", controlService.CreateArticle)

	return router
}

func TestPostSuccess(t *testing.T) {
	router := setupRouterCommand()

	requestBody := strings.NewReader(`{
    	"author":"test",
    	"title":"ini test title",
    	"body":"ini adalah test body article"
		}`,
	)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/articles", requestBody)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "success", keyValue["status"])
	assert.NotNil(t, keyValue["data"].(map[string]interface{})["id"])
	assert.Equal(t, "test", keyValue["data"].(map[string]interface{})["author"])
	assert.Equal(t, "ini test title", keyValue["data"].(map[string]interface{})["title"])
	assert.Equal(t, "ini adalah test body article", keyValue["data"].(map[string]interface{})["body"])
	assert.NotNil(t, keyValue["data"].(map[string]interface{})["created_at"])
}

func TestPostSuccess2(t *testing.T) {
	router := setupRouterCommand()

	requestBody := strings.NewReader(`{
    	"author":"someone",
    	"title":"this is title test",
    	"body":"this is test for the body"
		}`,
	)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/articles", requestBody)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 200, result.StatusCode)
	assert.Equal(t, "success", keyValue["status"])
	assert.NotNil(t, keyValue["data"].(map[string]interface{})["id"])
	assert.Equal(t, "someone", keyValue["data"].(map[string]interface{})["author"])
	assert.Equal(t, "this is title test", keyValue["data"].(map[string]interface{})["title"])
	assert.Equal(t, "this is test for the body", keyValue["data"].(map[string]interface{})["body"])
	assert.NotNil(t, keyValue["data"].(map[string]interface{})["created_at"])
}

func TestPostEmptyBody(t *testing.T) {
	router := setupRouterCommand()

	requestBody := strings.NewReader(``)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/articles", requestBody)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 400, result.StatusCode)
	assert.Equal(t, "error", keyValue["status"])
	assert.Nil(t, keyValue["data"])
}

func TestPostEmptyField(t *testing.T) {
	router := setupRouterCommand()

	requestBody := strings.NewReader(`{}`)
	request := httptest.NewRequest(http.MethodPost, "http://localhost:3000/articles", requestBody)

	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, request)

	result := recorder.Result()
	var keyValue map[string]interface{}
	body, _ := io.ReadAll(result.Body)
	json.Unmarshal(body, &keyValue)

	assert.Equal(t, 400, result.StatusCode)
	assert.Equal(t, "error", keyValue["status"])
}
