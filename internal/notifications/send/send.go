package notifications

import (
	"context"

	public_api "github.com/mayye4ka/pinder-api/api/go"
	notification_api "github.com/mayye4ka/pinder-api/notifications/go"
	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
)

const notificationsExchangeName = "notifications"

type NotificationSender struct {
	rabbit *amqp.Connection
	logger *zerolog.Logger
}

func NewNotificationSender(rabbit *amqp.Connection, logger *zerolog.Logger) (*NotificationSender, error) {
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
		notificationsExchangeName,
		"fanout",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		logger.Err(err).Msg("can't declare rabbitmq exchange")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't declare rabbitmq exchange",
		}
	}
	return &NotificationSender{
		rabbit: rabbit,
		logger: logger,
	}, nil
}

func (n *NotificationSender) SendMessage(ctx context.Context, userId uint64, notification models.MessageSend) error {
	return n.notify(
		ctx,
		userId,
		&public_api.DataPackage{
			Data: &public_api.DataPackage_IncomingMessageNotification{
				IncomingMessageNotification: &public_api.IncomingMessageNotification{
					ChatId:      notification.ChatID,
					MessageId:   notification.MessageID,
					SentByMe:    notification.SentByMe,
					ContentType: msgContentTypeToProto(notification.ContentType),
					Payload:     notification.Payload,
				},
			},
		},
	)
}

func (n *NotificationSender) NotifyLiked(ctx context.Context, userId uint64, notification models.LikeNotification) error {
	return n.notify(
		ctx,
		userId,
		&public_api.DataPackage{
			Data: &public_api.DataPackage_IncomingLikeNotification{
				IncomingLikeNotification: &public_api.IncomingLikeNotification{
					OpponentName:  notification.Name,
					OpponentPhoto: notification.Photo,
				},
			},
		},
	)
}

func (n *NotificationSender) NotifyMatch(ctx context.Context, userId uint64, notification models.MatchNotification) error {
	return n.notify(
		ctx,
		userId,
		&public_api.DataPackage{
			Data: &public_api.DataPackage_MatchNotification{
				MatchNotification: &public_api.MatchNotification{
					OpponentName:  notification.Name,
					OpponentPhoto: notification.Photo,
				},
			},
		},
	)
}

func (n *NotificationSender) SendTranscribedMessage(ctx context.Context, userId uint64, notification models.MessageTranscibed) error {
	return n.notify(
		ctx,
		userId,
		&public_api.DataPackage{
			Data: &public_api.DataPackage_VoiceTranscribedNotification{
				VoiceTranscribedNotification: &public_api.VoiceTranscribedNotification{
					MessageId: notification.MessageID,
					Text:      notification.Text,
				},
			},
		},
	)
}

func (n *NotificationSender) notify(ctx context.Context, userId uint64, data *public_api.DataPackage) error {
	bytes, err := proto.Marshal(&notification_api.UserNotification{
		UserId:      userId,
		DataPackage: data,
	})
	if err != nil {
		n.logger.Err(err).Msg("can't unmarshal notification")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't unmarshal notification",
		}
	}
	ch, err := n.rabbit.Channel()
	if err != nil {
		n.logger.Err(err).Msg("can't open rabbitmq channel")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't open rabbitmq channel",
		}
	}
	defer ch.Close()
	err = ch.PublishWithContext(
		ctx,
		notificationsExchangeName,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(bytes),
		},
	)
	if err != nil {
		n.logger.Err(err).Msg("can't send notification")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't send notification",
		}
	}
	return nil
}

func msgContentTypeToProto(ct models.MsgContentType) public_api.MESSAGE_CONTENT_TYPE {
	switch ct {
	case models.ContentPhoto:
		return public_api.MESSAGE_CONTENT_TYPE_PHOTO
	case models.ContentVoice:
		return public_api.MESSAGE_CONTENT_TYPE_VOICE
	default:
		return public_api.MESSAGE_CONTENT_TYPE_TEXT
	}
}
