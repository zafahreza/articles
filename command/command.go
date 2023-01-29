package main

import (
	"articles/command/app"
	"articles/command/broker"
	"articles/command/controller"
	"articles/command/repository"
	"articles/command/service"
	"github.com/gin-gonic/gin"
)

func main() {
	db := app.InitDB()
	ch := app.InitBroker()

	repo := repository.NewRepository()
	messageBroker := broker.NewBroker()
	services := service.NewService(repo, messageBroker, db, ch)
	controlService := controller.NewController(services)

	router := gin.Default()

	router.POST("/articles", controlService.CreateArticle)
	router.Run(":3000")
}
