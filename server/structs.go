package server

import (
	"time"
)

// auth

type LoginRequest struct {
	PhoneNumber string
	Password    string
}

type LoginResponse struct {
	Token string
}

type RegisterRequest struct {
	PhoneNumber string
	Password    string
}

type RegisterResponse struct {
	Token string
}

// my

type GetProfileResponse struct {
	Profile Profile
}

type UpdProfileRequest struct {
	Token      string
	NewProfile Profile
}

type GetPreferencesResponse struct {
	Preferences Preferences
}

type UpdPreferencesRequest struct {
	Token          string
	NewPreferences Preferences
}

type DelPhotoRequest struct {
	Token    string
	PhotoKey string
}

// main process

type NextPartnerResponse struct {
	Partner Profile
}

type SwipeRequest struct {
	Token        string
	CandidateID  uint64
	SwipeVerdict SwipeVerdict
}

// other

type Empty struct{}

type RequestWithToken struct {
	Token string
}

type Profile struct {
	Name         string
	Gender       Gender
	Age          int
	Bio          string
	Photos       []string
	LocationLat  float64
	LocationLon  float64
	LocationName string
}

type Preferences struct {
	MaxAge           int
	MinAge           int
	Gender           Gender
	LocationLat      float64
	LocationLon      float64
	LocationRadiusKm float64
}

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type SwipeVerdict string

const (
	SwipeLike    SwipeVerdict = "like"
	SwipeDislike SwipeVerdict = "dislike"
)

//chat

type Chat struct {
	ChatID uint64
	Name   string
	Photo  string
}

type MsgContentType string

const (
	ContentText  MsgContentType = "text"
	ContentPhoto MsgContentType = "photo"
	ContentVoice MsgContentType = "voice"
)

type Message struct {
	ID          uint64
	SentByMe    bool
	ContentType MsgContentType
	Payload     string
	CreatedAt   time.Time
}

type ListChatsResponse struct {
	Chats []Chat
}

type ListMessagesRequest struct {
	Token  string
	ChatId uint64
}

type ListMessagesResponse struct {
	Messages []Message
}
