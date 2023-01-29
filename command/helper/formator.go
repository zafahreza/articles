package helper

import "articles/command/domain"

func ArticleToResponse(status string, code int, data interface{}) domain.JsonResponse {
	reponse := domain.JsonResponse{
		Status: status,
		Code: code,
		Data: data,
	}

	return reponse
}

