package stt

import (
	"context"
	"time"

	stt_api "github.com/mayye4ka/pinder-api/stt/go"
	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

const (
	sttResultExchangeName = "stt_results"
	sttResultRoutingKey   = "stt_result"
	handleResultTimeout   = time.Minute
)

type SttResultHandler interface {
	HandleSttResult(ctx context.Context, res models.SttResult) error
}

type SttResultReceiver struct {
	rabbit     *amqp.Connection
	handler    SttResultHandler
	logger     *zerolog.Logger
	finish     chan struct{}
	finishDone chan struct{}
}

func NewResultReceiver(rabbit *amqp.Connection, handler SttResultHandler, logger *zerolog.Logger) *SttResultReceiver {
	return &SttResultReceiver{
		rabbit:     rabbit,
		handler:    handler,
		logger:     logger,
		finish:     make(chan struct{}),
		finishDone: make(chan struct{}),
	}
}

func (s *SttResultReceiver) Start(ctx context.Context) error {
	ch, err := s.rabbit.Channel()
	if err != nil {
		s.logger.Err(err).Msg("can't open rabbitmq channel")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't open rabbitmq channel",
		}
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		sttResultExchangeName,
		"direct",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Err(err).Msg("can't declare stt results exchange")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't declare stt results exchange",
		}
	}

	q, err := ch.QueueDeclare(
		"",
		true,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		s.logger.Err(err).Msg("can't declare stt results queue")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't declare stt results queue",
		}
	}

	err = ch.QueueBind(
		q.Name,
		sttResultRoutingKey,
		sttResultExchangeName,
		false,
		nil,
	)
	if err != nil {
		s.logger.Err(err).Msg("can't bind stt results queue")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't bind stt results queue",
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		false,
		true,
		false,
		false,
		nil,
	)
	if err != nil {
		s.logger.Err(err).Msg("can't start rabbitmq consumer")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't start rabbitmq consumer",
		}
	}

	for {
		select {
		case <-ctx.Done():
			close(s.finishDone)
			return nil
		case <-s.finish:
			close(s.finishDone)
			return nil
		case msg := <-msgs:
			var result stt_api.SttResult
			err = proto.Unmarshal(msg.Body, &result)
			if err != nil {
				s.logger.Err(err).Msg("can't unmarshal stt result")
				return &errs.CodableError{
					Code:    errs.CodeInternal,
					Message: "can't unmarshal stt result",
				}
			}
			ctxTo, cancel := context.WithTimeout(ctx, handleResultTimeout)
			err := s.handler.HandleSttResult(ctxTo, models.SttResult{
				UserID:    result.UserId,
				MessageID: result.MessageId,
				Text:      result.Text,
			})
			cancel()
			if err != nil {
				s.logger.Err(err).Msg("can't handle stt result")
				return &errs.CodableError{
					Code:    errs.CodeInternal,
					Message: "can't handle stt result",
				}
			} else {
				err = ch.Ack(msg.DeliveryTag, false)
				if err != nil {
					s.logger.Err(err).Msg("can't ack stt result")
					return &errs.CodableError{
						Code:    errs.CodeInternal,
						Message: "can't ack stt result",
					}
				}
			}
		}
	}
}

func (s *SttResultReceiver) Stop(ctx context.Context) error {
	close(s.finish)
	select {
	case <-s.finishDone:
	case <-ctx.Done():
	}
	return nil
}
