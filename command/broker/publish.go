package broker

import (
	"articles/command/domain"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
)

type Publisher interface {
	Publish(channel *amqp.Channel, article domain.Article) error
}

type Broker struct {
}

func NewBroker() Publisher {
	return &Broker{}
}

func (b Broker) Publish(channel *amqp.Channel, article domain.Article) error {
	result, err := json.Marshal(article)
	if err != nil {
		return err
	}
	message := amqp.Publishing{ContentType: "application/json", Body: result}
	err = channel.Publish("", "article", false, false, message)
	if err != nil {
		return err
	}
	fmt.Println("message published")
	return nil
}
