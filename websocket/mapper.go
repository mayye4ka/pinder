package websocket

import (
	"pinder/models"
)

func unmapContentType(content models.MsgContentType) MsgContentType {
	if content == models.ContentPhoto {
		return ContentPhoto
	}
	if content == models.ContentText {
		return ContentText
	}
	return ContentVoice
}
