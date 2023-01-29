package controller

import (
	"articles/command/domain"
	"articles/command/helper"
	"articles/command/service"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type Controller struct {
	service service.SaveService
}

func NewController(saveService service.SaveService) *Controller {
	return &Controller{
		service: saveService,
	}
}

func (con *Controller) CreateArticle(c *gin.Context) {
	req := domain.RequestArticle{}

	err := c.ShouldBindJSON(&req)

	if err == io.EOF {
		response := helper.ArticleToResponse("error", http.StatusBadRequest, nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	if err != nil {
		errors := helper.ValidationError(err)
		errorM := gin.H{"errors": errors}

		response := helper.ArticleToResponse("error", http.StatusBadRequest, errorM)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	result, err := con.service.CreateArticle(c, req)
	if err != nil {
		response := helper.ArticleToResponse("error", http.StatusInternalServerError, nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := helper.ArticleToResponse("success", http.StatusOK, result)
	c.JSON(http.StatusOK, response)
}
