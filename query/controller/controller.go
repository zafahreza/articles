package controller

import (
	"articles/query/domain"
	"articles/query/helper"
	"articles/query/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Controller struct {
	Service service.ServiceQuery
}

func NewController(s service.ServiceQuery) *Controller {
	return &Controller{
		Service: s,
	}
}

func (con *Controller) SearchArticle(c *gin.Context) {
	param := domain.Params{
		Author: c.Query("author"),
		Query:  c.Query("query"),
	}
	results, err := con.Service.Search(param)
	if err != nil {
		response := helper.ArticleToResponse("error", http.StatusInternalServerError, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	if len(results) == 0 {
		response := helper.ArticleToResponse("not found", http.StatusNotFound, nil)
		c.JSON(http.StatusNotFound, response)
		return
	}
	response := helper.ArticleToResponse("success", http.StatusOK, results)
	c.JSON(http.StatusOK, response)
}
