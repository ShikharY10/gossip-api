package handler

import (
	"github.com/ShikharY10/gbAPI/logger"
	"github.com/streadway/amqp"
)

type QueueHandler struct {
	channel *amqp.Channel
	logger  *logger.Logger
}

func CreateQueueHandler(channel *amqp.Channel, logger *logger.Logger) *QueueHandler {
	queueHandler := &QueueHandler{
		channel: channel,
		logger:  logger,
	}
	return queueHandler
}

func (r *QueueHandler) Produce(name string, job []byte) error {
	err := r.channel.Publish(
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
