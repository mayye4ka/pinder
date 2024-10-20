package stt

import (
	"context"

	"github.com/mayye4ka/pinder/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SttTaskCreator struct {
	rabbit *amqp.Connection
}

func NewTaskCreator(rabbit *amqp.Connection) *SttTaskCreator {
	return &SttTaskCreator{
		rabbit: rabbit,
	}
}

func (s *SttTaskCreator) PutTask(task models.SttTask) error {
	// TODO: put to rabbit
	return nil
}

type SttResultHandler interface {
	HandleSttResult(res models.SttResult) error
}

type SttResultReceiver struct {
	rabbit  *amqp.Connection
	handler SttResultHandler
}

func NewResultReceiver(rabbit *amqp.Connection, handler SttResultHandler) *SttResultReceiver {
	return &SttResultReceiver{
		rabbit:  rabbit,
		handler: handler,
	}
}

func (s *SttResultReceiver) Start(ctx context.Context) error {
	<-ctx.Done()
	// TODO: receive from rabbit and call handler
	return nil
}
