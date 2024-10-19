package notifications

import (
	notification_api "github.com/mayye4ka/pinder-api/notifications/go"
	amqp "github.com/rabbitmq/amqp091-go"
)

const notificationsExchangeName = "notifications"

type Notifier struct {
	rabbit *amqp.Connection

	resultChan chan *notification_api.UserNotification
}

func NewNotifier(rabbit *amqp.Connection) *Notifier {
	return &Notifier{
		rabbit:     rabbit,
		resultChan: make(chan *notification_api.UserNotification, 1024),
	}
}
