package config

import (
	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Channel *amqp.Channel
}

func RabbitInit(env *ENV) (*RabbitMQ, error) {
	var rabbitMQ RabbitMQ
	var address string = "amqp://" + env.RabbitMQUsername + ":" + env.RabbitMQPassword + "@" + env.RabbitMQHost + ":" + env.RabbitMQPort + "/"
	conn, err := amqp.Dial(address)
	if err != nil {
		return nil, err
	}
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	rabbitMQ.Channel = channel
	return &rabbitMQ, nil
}
