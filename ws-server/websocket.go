package websocket

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	public_api "github.com/mayye4ka/pinder-api/public_api/go"
	"github.com/mayye4ka/pinder/models"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

const (
	tokenHeader             = "Authorization"
	authorizationTrimPrefix = "Bearer "
)

type UserWsNotifier struct {
	auth Authenticator

	connStore   map[uint64][]*websocket.Conn
	connStoreMu sync.RWMutex
}

type Authenticator interface {
	UnpackToken(ctx context.Context, token string) (uint64, error)
}

func NewUserWsNotifier(auth Authenticator) *UserWsNotifier {
	return &UserWsNotifier{
		auth: auth,
	}
}

func (i *UserWsNotifier) addUser(id uint64, conn *websocket.Conn) {
	i.connStoreMu.Lock()
	i.connStore[id] = append(i.connStore[id], conn)
	i.connStoreMu.Unlock()
	go i.serveConn(id, conn)
}

func (i *UserWsNotifier) serveConn(id uint64, conn *websocket.Conn) {
	for {
		// TODO:
	}
}

func (i *UserWsNotifier) sendBytes(id uint64, bytes []byte) {
	i.connStoreMu.RLock()
	for _, conn := range i.connStore[id] {
		err := conn.WriteMessage(websocket.BinaryMessage, bytes)
		if err != nil {
			log.Println(err)
		}
	}
	i.connStoreMu.RUnlock()
}

func (i *UserWsNotifier) notify(userId uint64, notification *public_api.DataPackage) error {
	bytes, err := proto.Marshal(notification)
	if err != nil {
		return nil
	}
	i.sendBytes(userId, bytes)
	return nil
}

func (i *UserWsNotifier) NotifyLiked(userId uint64, notification models.LikeNotification) error {
	ntfc := &public_api.DataPackage{
		Data: &public_api.DataPackage_IncomingLikeNotification{
			IncomingLikeNotification: &public_api.IncomingLikeNotification{
				OpponentName:  notification.Name,
				OpponentPhoto: notification.Photo,
			},
		},
	}
	return i.notify(userId, ntfc)
}
func (i *UserWsNotifier) NotifyMatch(userId uint64, notification models.MatchNotification) error {
	ntfc := &public_api.DataPackage{
		Data: &public_api.DataPackage_MatchNotification{
			MatchNotification: &public_api.MatchNotification{
				OpponentName:  notification.Name,
				OpponentPhoto: notification.Photo,
			},
		},
	}
	return i.notify(userId, ntfc)
}

func (i *UserWsNotifier) SendMessage(userId uint64, notification models.MessageSend) error {
	ntfc := &public_api.DataPackage{
		Data: &public_api.DataPackage_IncomingMessageNotification{
			IncomingMessageNotification: &public_api.IncomingMessageNotification{
				ChatId:      notification.ChatID,
				MessageId:   notification.MessageID,
				SentByMe:    notification.SentByMe,
				ContentType: msgContentTypeToProto(notification.ContentType),
				Payload:     notification.Payload,
			},
		},
	}
	return i.notify(userId, ntfc)
}

func (i *UserWsNotifier) SendTranscribedMessage(userId uint64, notification models.MessageTranscibed) error {
	ntfc := &public_api.DataPackage{
		Data: &public_api.DataPackage_VoiceTranscribedNotification{
			VoiceTranscribedNotification: &public_api.VoiceTranscribedNotification{
				MessageId: notification.MessageID,
				Text:      notification.Text,
			},
		},
	}
	return i.notify(userId, ntfc)
}

var upgrader = websocket.Upgrader{}

func (n *UserWsNotifier) wsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get(tokenHeader) == "" {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	token := strings.TrimPrefix(r.Header.Get("Authentication"), authorizationTrimPrefix)
	userId, err := n.auth.UnpackToken(r.Context(), token)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	n.addUser(userId, conn)
}

func (n *UserWsNotifier) Start(port int) error {
	http.HandleFunc("/ws", n.wsHandler)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
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
