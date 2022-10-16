package config

import (
	"fmt"
	"math/rand"

	"github.com/ShikharY10/gbAPI/models"
	"github.com/streadway/amqp"
)

type RMQ struct {
	RedisDB *models.Redis
	Msgs    <-chan amqp.Delivery
	ch      *amqp.Channel
}

func RabbitInit(RMQIP string, username string, password string) *RMQ {
	var r RMQ
	var address string = "amqp://" + username + ":" + password + "@" + RMQIP + ":5672/"
	conn, err := amqp.Dial(address)
	if err != nil {
		fmt.Println("[ERROR] : ", err.Error())
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("[ERROR] : ", err.Error())
	}
	r.ch = ch
	fmt.Println("RabbitMQ Connected")
	return &r
}

func (r *RMQ) GetEngineChannel() string {
	names := r.RedisDB.GetEngineName()
	fmt.Println("LEN: ", len(names))
	randomIndex := rand.Intn(len(names))
	pick := names[randomIndex]
	return pick
}

func (r *RMQ) Produce(name string, job []byte) error {
	err := r.ch.Publish(
		"",
		name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        job,
		},
	)
	return err
}
