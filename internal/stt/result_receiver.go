package stt

import (
	"context"

	"github.com/mayye4ka/pinder/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SttResultHandler interface {
	HandleSttResult(ctx context.Context, res models.SttResult) error
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

func (s *SttResultReceiver) Stop(ctx context.Context) error {
	return nil
}
