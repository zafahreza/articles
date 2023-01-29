package app

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"os"
)

func InitBroker() *amqp.Channel {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	host := os.Getenv("RABBITMQ_HOST")

	url := fmt.Sprintf("amqp://guest:guest@%s:5672/", host)
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}

	return ch

}
