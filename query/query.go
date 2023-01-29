package main

import (
	"articles/query/app"
	"articles/query/broker"
	"articles/query/cache"
	"articles/query/controller"
	"articles/query/domain"
	"articles/query/repository"
	"articles/query/service"
	"github.com/gin-gonic/gin"
)

func main() {
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

	router := gin.Default()

	router.GET("/articles", controllers.SearchArticle)

	router.Run(":3001")

}
