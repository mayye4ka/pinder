package notifications

import (
	"context"

	notification_api "github.com/mayye4ka/pinder-api/notifications/go"
	"github.com/mayye4ka/pinder/internal/errs"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

const notificationsExchangeName = "notifications"

type NotificationReceiver struct {
	rabbit     *amqp.Connection
	logger     *zerolog.Logger
	resultChan chan *notification_api.UserNotification
	finish     chan struct{}
	finishDone chan struct{}
}

func NewNotificationReceiver(rabbit *amqp.Connection, logger *zerolog.Logger) *NotificationReceiver {
	return &NotificationReceiver{
		rabbit:     rabbit,
		logger:     logger,
		resultChan: make(chan *notification_api.UserNotification, 1024),
		finish:     make(chan struct{}),
		finishDone: make(chan struct{}),
	}
}

func (n *NotificationReceiver) Start(ctx context.Context) error {
	ch, err := n.rabbit.Channel()
	if err != nil {
		n.logger.Err(err).Msg("can't open rabbitmq channel")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't open rabbitmq channel",
		}
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		notificationsExchangeName,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		n.logger.Err(err).Msg("can't declare rabbitmq exchange")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't declare rabbitmq exchange",
		}
	}

	q, err := ch.QueueDeclare(
		"",
		false,
		true,
		true,
		false,
		nil,
	)
	if err != nil {
		n.logger.Err(err).Msg("can't declare rabbitmq queue")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't declare rabbitmq queue",
		}
	}

	err = ch.QueueBind(
		q.Name,
		"",
		notificationsExchangeName,
		false,
		nil,
	)
	if err != nil {
		n.logger.Err(err).Msg("can't bind rabbitmq queue")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't bind rabbitmq queue",
		}
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		n.logger.Err(err).Msg("can't start rabbitmq consumer")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't start rabbitmq consumer",
		}
	}
	for {
		select {
		case <-ctx.Done():
			close(n.resultChan)
			close(n.finishDone)
			return nil
		case <-n.finish:
			close(n.resultChan)
			close(n.finishDone)
			return nil
		case msg := <-msgs:
			var notification notification_api.UserNotification
			err = proto.Unmarshal(msg.Body, &notification)
			if err != nil {
				n.logger.Err(err).Msg("can't unmarshal notification")
				return &errs.CodableError{
					Code:    errs.CodeInternal,
					Message: "can't unmarshal notification",
				}
			}
			n.resultChan <- &notification
		}
	}
}

func (n *NotificationReceiver) Stop(ctx context.Context) error {
	close(n.finish)
	select {
	case <-n.finishDone:
	case <-ctx.Done():
	}
	return nil
}

func (n *NotificationReceiver) Notifications() <-chan *notification_api.UserNotification {
	return n.resultChan
}
