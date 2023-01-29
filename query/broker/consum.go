package broker

import (
	"articles/query/domain"
	"encoding/json"
	"github.com/streadway/amqp"
)

type Consumer interface {
	Consum(ch *amqp.Channel, articleChan chan []byte)
}

type Broker struct {
}

func NewConsumer() Consumer {
	return &Broker{}
}

func (b Broker) Consum(ch *amqp.Channel, articleChan chan []byte) {
	msgs, err := ch.Consume("article", "", true, false, false, false, nil)
	if err != nil {
		panic(err)
	}

	msg := domain.Article{}

	newChan := make(chan bool)
	go func() {
		for m := range msgs {
			err = json.Unmarshal(m.Body, &msg)
			if err != nil {
				panic(err)
			}

			articleChan <- m.Body

		}
	}()
	<-newChan
}
