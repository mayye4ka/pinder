package websocket

import (
	"encoding/json"
	"log"
	"pinder/models"
	"sync"

	"github.com/gorilla/websocket"
)

type UserWsInteractor struct {
	connStore   map[uint64][]*websocket.Conn
	connStoreMu sync.RWMutex
}

func NewUserInteractor() *UserWsInteractor {
	return &UserWsInteractor{}
}

func (i *UserWsInteractor) AddUser(id uint64, conn *websocket.Conn) {
	i.connStoreMu.Lock()
	i.connStore[id] = append(i.connStore[id], conn)
	i.connStoreMu.Unlock()
	go i.keepaliveConn(id, conn)
}

func (i *UserWsInteractor) keepaliveConn(id uint64, conn *websocket.Conn) {
	for {
		// TODO: keep conns alive and close and remove from map when they die
	}
}

func (i *UserWsInteractor) sendDataPackage(ids []uint64, dpkg DataPackage) {
	i.connStoreMu.RLock()
	for _, id := range ids {
		for _, conn := range i.connStore[id] {
			err := conn.WriteJSON(dpkg)
			if err != nil {
				log.Println(err)
			}
		}
	}
	i.connStoreMu.RUnlock()
}

func (i *UserWsInteractor) SendMessage(chat models.Chat, message models.Message) error {
	senderPayload := IncomingMessageNotification{
		ChatID:      chat.ID,
		MessageID:   message.ID,
		SentByMe:    true,
		ContentType: unmapContentType(message.ContentType),
		Payload:     message.Payload,
	}
	receiverPayload := IncomingMessageNotification{
		ChatID:      chat.ID,
		MessageID:   message.ID,
		SentByMe:    false,
		ContentType: unmapContentType(message.ContentType),
		Payload:     message.Payload,
	}
	senderPayloadString, err := json.Marshal(senderPayload)
	if err != nil {
		return err
	}
	receiverPayloadString, err := json.Marshal(receiverPayload)
	if err != nil {
		return err
	}
	dpkgSender := DataPackage{
		DataType: incomingMessageNotification,
		Payload:  string(senderPayloadString),
	}
	dpkgReceiver := DataPackage{
		DataType: incomingMessageNotification,
		Payload:  string(receiverPayloadString),
	}

	sender := message.SenderID
	receiver := chat.User1
	if sender == receiver {
		receiver = chat.User2
	}
	i.sendDataPackage([]uint64{sender}, dpkgSender)
	i.sendDataPackage([]uint64{receiver}, dpkgReceiver)
	return nil
}

func (i *UserWsInteractor) NotifyLiked(userId uint64, opName, opPhoto string) error {
	payload := IncomingLikeNotification{
		OpponentName:  opName,
		OpponentPhoto: opPhoto,
	}
	payloadString, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	dpkg := DataPackage{
		DataType: incomingLikeNotification,
		Payload:  string(payloadString),
	}
	i.sendDataPackage([]uint64{userId}, dpkg)
	return nil
}

func (i *UserWsInteractor) NotifyMatch(userId uint64, opName, opPhoto string) error {
	payload := MatchNotification{
		OpponentName:  opName,
		OpponentPhoto: opPhoto,
	}
	payloadString, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	dpkg := DataPackage{
		DataType: matchNotification,
		Payload:  string(payloadString),
	}
	i.sendDataPackage([]uint64{userId}, dpkg)
	return nil
}
