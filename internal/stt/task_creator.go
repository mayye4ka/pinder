package stt

import (
	"context"

	"github.com/mayye4ka/pinder/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
)

type SttTaskCreator struct {
	rabbit *amqp.Connection
}

func NewTaskCreator(rabbit *amqp.Connection) (*SttTaskCreator, error) {
	return &SttTaskCreator{
		rabbit: rabbit,
	}, nil
}

func (s *SttTaskCreator) PutTask(ctx context.Context, task models.SttTask) error {
	// TODO: put to rabbit
	return nil
}
