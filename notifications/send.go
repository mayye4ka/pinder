package notifications

import (
	public_api "github.com/mayye4ka/pinder-api/api/go"
	notification_api "github.com/mayye4ka/pinder-api/notifications/go"
	"github.com/mayye4ka/pinder/models"
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/protobuf/proto"
)

func (n *Notifier) SendMessage(userId uint64, notification models.MessageSend) error {
	return n.notify(
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

func (n *Notifier) NotifyLiked(userId uint64, notification models.LikeNotification) error {
	return n.notify(
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

func (n *Notifier) NotifyMatch(userId uint64, notification models.MatchNotification) error {
	return n.notify(
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

func (n *Notifier) SendTranscribedMessage(userId uint64, notification models.MessageTranscibed) error {
	return n.notify(
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

func (n *Notifier) notify(userId uint64, data *public_api.DataPackage) error {
	bytes, err := proto.Marshal(&notification_api.UserNotification{
		UserId:      userId,
		DataPackage: data,
	})
	if err != nil {
		return nil
	}
	ch, err := n.rabbit.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	err = ch.Publish(
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
		return err
	}
	return nil
}
