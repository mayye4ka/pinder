package models

import (
	"log"
	"time"

	"github.com/jftuga/geodist"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type User struct {
	ID          uint64
	PhoneNumber string
	PassHash    string
}

func (User) TableName() string {
	return "users"
}

type Profile struct {
	UserID       uint64
	Name         string
	Gender       Gender
	Age          int
	Bio          string
	LocationLat  float64
	LocationLon  float64
	LocationName string
}

func (Profile) TableName() string {
	return "profiles"
}

type Photo struct {
	UserID   uint64
	PhotoKey string
}

func (Photo) TableName() string {
	return "photos"
}

type Preferences struct {
	UserID           uint64
	MaxAge           int
	MinAge           int
	Gender           Gender
	LocationLat      float64
	LocationLon      float64
	LocationRadiusKm float64
}

func (Preferences) TableName() string {
	return "preferences"
}

type PairAttempt struct {
	ID        uint64
	User1     uint64
	User2     uint64
	State     PAState
	CreatedAt time.Time
}

func (PairAttempt) TableName() string {
	return "pair_attempts"
}

type PAState string

const (
	PAStatePending  PAState = "pending"
	PAStateMatch    PAState = "match"
	PAStateMismatch PAState = "mismatch"
)

type PEType string

const (
	PETypePACreated         PEType = "pa_created"
	PETypeSentToUser1       PEType = "sent_to_user_1"
	PETypeUser1Liked        PEType = "user_1_liked"
	PETypeUser1Disliked     PEType = "user_1_disliked"
	PETypeSentToUser2       PEType = "sent_to_user_2"
	PETypeUser2Liked        PEType = "user_2_liked"
	PETypeUser2Disliked     PEType = "user_2_disliked"
	PETypePairAttemptFailed PEType = "pair_attempt_failed"
	PETypePairCreated       PEType = "pair_created"
)

type PairEvent struct {
	ID        uint64
	PAID      uint64
	CreatedAt time.Time
	EventType PEType
}

func (PairEvent) TableName() string {
	return "pair_events"
}

func (p *Preferences) ProfileMatches(profile Profile) bool {
	if p.Gender != "" && p.Gender != profile.Gender {
		return false
	}
	if p.MinAge != 0 && p.MinAge > profile.Age {
		return false
	}
	if p.MaxAge != 0 && p.MaxAge < profile.Age {
		return false
	}
	_, dst, err := geodist.VincentyDistance(geodist.Coord{
		Lat: p.LocationLat,
		Lon: p.LocationLon,
	}, geodist.Coord{
		Lat: profile.LocationLat,
		Lon: profile.LocationLon,
	})
	if err != nil {
		log.Println(err)
		return false
	}
	if dst > p.LocationRadiusKm {
		return false
	}
	return true
}

type Chat struct {
	ID    uint64
	User1 uint64
	User2 uint64
}

func (Chat) TableName() string {
	return "chats"
}

type MsgContentType string

const (
	ContentText  MsgContentType = "text"
	ContentPhoto MsgContentType = "photo"
	ContentVoice MsgContentType = "voice"
)

type Message struct {
	ID          uint64
	ChatID      uint64
	SenderID    uint64
	ContentType MsgContentType
	Payload     string
	CreatedAt   time.Time
}

func (Message) TableName() string {
	return "messages"
}
