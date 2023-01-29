package helper

import "articles/query/domain"

func ArticleToResponse(status string, code int, data interface{}) domain.JsonResponse {
	reponse := domain.JsonResponse{
		Status: status,
		Code:   code,
		Data:   data,
	}

	return reponse
}

func WhichQuery(params domain.Params) map[string]interface{} {
	type keyValue map[string]interface{}
	if params.Query == "" && params.Author == "" {
		query := keyValue{
			"sort": []keyValue{
				{
					"id": keyValue{
						"order": "desc",
					},
				},
			},
		}

		return query
	} else if params.Query == "" && params.Author != "" {
		query := keyValue{
			"query": keyValue{
				"bool": keyValue{
					"must": keyValue{
						"match": keyValue{
							"author": params.Author,
						},
					},
				},
			},
			"sort": []keyValue{
				{
					"id": keyValue{
						"order": "desc",
					},
				},
			},
		}

		return query
	} else if params.Query != "" && params.Author == "" {
		query := keyValue{
			"query": keyValue{
				"bool": keyValue{
					"should": []keyValue{
						{
							"match": map[string]interface{}{
								"title": params.Query,
							},
						},
						{
							"match": map[string]interface{}{
								"body": params.Query,
							},
						},
					},
				},
			},
			"sort": []keyValue{
				{
					"id": keyValue{
						"order": "desc",
					},
				},
			},
		}

		return query
	} else {
		query := keyValue{
			"query": keyValue{
				"bool": keyValue{
					"must": []keyValue{
						{
							"term": keyValue{
								"author": params.Author,
							},
						},
						{
							"bool": keyValue{
								"must": keyValue{
									"bool": keyValue{
										"should": []keyValue{
											{
												"term": keyValue{
													"title": params.Query,
												},
											},
											{
												"term": keyValue{
													"body": params.Query,
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"sort": []keyValue{
				{
					"id": keyValue{
						"order": "desc",
					},
				},
			},
		}
		return query
	}
}
