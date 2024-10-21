package stt

import (
	"context"

	stt_api "github.com/mayye4ka/pinder-api/stt/go"
	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

const (
	sttTaskExchangeName = "stt_tasks"
	sttTaskRoutingKey   = "stt_task"
)

type SttTaskCreator struct {
	rabbit *amqp.Connection
	logger *zerolog.Logger
}

func NewTaskCreator(rabbit *amqp.Connection, logger *zerolog.Logger) (*SttTaskCreator, error) {
	ch, err := rabbit.Channel()
	if err != nil {
		logger.Err(err).Msg("can't open rabbitmq channel")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't open rabbitmq channel",
		}
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		sttTaskExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Err(err).Msg("can't declare stt tasks exchange")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't declare stt tasks exchange",
		}
	}

	return &SttTaskCreator{
		rabbit: rabbit,
		logger: logger,
	}, nil
}

func (s *SttTaskCreator) PutTask(ctx context.Context, task models.SttTask) error {
	t := &stt_api.SttTask{
		UserId:    task.UserID,
		MessageId: task.MessageID,
		Speech:    []byte(task.Speech),
	}
	body, err := proto.Marshal(t)
	if err != nil {
		s.logger.Err(err).Msg("can't marshal stt task")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't marshal stt task",
		}
	}
	ch, err := s.rabbit.Channel()
	if err != nil {
		s.logger.Err(err).Msg("can't open rabbitmq channel")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't open rabbitmq channel",
		}
	}
	defer ch.Close()

	err = ch.PublishWithContext(ctx, sttTaskExchangeName, sttTaskRoutingKey, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         []byte(body),
	})
	if err != nil {
		s.logger.Err(err).Msg("can't publish stt task to rabbitmq")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't publish stt task to rabbitmq",
		}
	}

	return nil
}
