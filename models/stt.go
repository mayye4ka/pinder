package models

type SttTask struct {
	UserID    uint64
	MessageID uint64
	Speech    string
}

type SttResult struct {
	UserID    uint64
	MessageID uint64
	Text      string
}
