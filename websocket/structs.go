package websocket

type MsgContentType string

const (
	ContentText  MsgContentType = "text"
	ContentPhoto MsgContentType = "photo"
	ContentVoice MsgContentType = "voice"
)

type DataType string

const (
	matchNotification           DataType = "match_notification"
	incomingLikeNotification    DataType = "incoming_like_notification"
	incomingMessageNotification DataType = "incoming_message_notification"
)

type DataPackage struct {
	DataType DataType
	Payload  string
}

type MatchNotification struct {
	OpponentName  string
	OpponentPhoto string
}

type IncomingLikeNotification struct {
	OpponentName  string
	OpponentPhoto string
}

type IncomingMessageNotification struct {
	ChatID      uint64
	MessageID   uint64
	SentByMe    bool
	ContentType MsgContentType
	Payload     string
}
