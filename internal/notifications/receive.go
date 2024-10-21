package notifications

import (
	"context"

	notification_api "github.com/mayye4ka/pinder-api/notifications/go"
	"google.golang.org/protobuf/proto"
)

func (n *Notifier) Start(ctx context.Context) error {
	ch, err := n.rabbit.Channel()
	if err != nil {
		return err
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
		return err
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
		return err
	}

	err = ch.QueueBind(
		q.Name,
		"",
		notificationsExchangeName,
		false,
		nil,
	)
	if err != nil {
		return err
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
		return err
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
				return err
			}
			n.resultChan <- &notification
		}
	}
}

func (n *Notifier) Stop(ctx context.Context) error {
	close(n.finish)
	select {
	case <-n.finishDone:
	case <-ctx.Done():
	}
	return nil
}

func (n *Notifier) Notifications() <-chan *notification_api.UserNotification {
	return n.resultChan
}
