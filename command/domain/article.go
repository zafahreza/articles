package domain

import "time"

type Article struct {
	Id        int `json:"id"`
	Author    string `json:"author"`
	Title     string `json:"title"`
	Body      string `json:"body"`
	CreatedAt time.Time `json:"created_at"`
}

type RequestArticle struct {
	Author string `json:"author" binding:"required"`
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body" binding:"required"`
}

type JsonResponse struct {
	Status string `json:"status"`
	Code int `json:"code"`
	Data interface{} `json:"data"`
}