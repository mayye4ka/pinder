package models

type MatchNotification struct {
	Name  string
	Photo string
}

type LikeNotification struct {
	Name  string
	Photo string
}

type MessageSend struct {
	ChatID      uint64
	MessageID   uint64
	SentByMe    bool
	ContentType MsgContentType
	Payload     string
}

type MessageTranscibed struct {
	MessageID uint64
	Text      string
}
