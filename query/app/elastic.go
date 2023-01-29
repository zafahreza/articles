package app

import (
	"fmt"
	"github.com/elastic/go-elasticsearch"
)

func InitESClient() *elasticsearch.Client {

	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}
	fmt.Println("ES initialized...")

	return client

}
