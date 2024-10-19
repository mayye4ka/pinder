package notifications

import (
	notification_api "github.com/mayye4ka/pinder-api/notifications/go"
	"google.golang.org/protobuf/proto"
)

func (n *Notifier) Start() error {
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
	for msg := range msgs {
		var notification *notification_api.UserNotification
		err = proto.Unmarshal(msg.Body, notification)
		if err != nil {
			return err
		}
		n.resultChan <- notification
	}
	return nil
}

func (n *Notifier) Notifications() <-chan *notification_api.UserNotification {
	return n.resultChan
}
