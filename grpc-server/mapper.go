package server

import (
	public_api "github.com/mayye4ka/pinder-api/public_api/go"
	"github.com/mayye4ka/pinder/models"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func profileShowcaseToCandidate(prof models.ProfileShowcase) *public_api.Candidate {
	return &public_api.Candidate{
		Profile: profileToProto(prof.Profile),
		Photos:  photosToLinkList(prof.Photos),
	}
}

func profileToProto(prof models.Profile) *public_api.Profile {
	return &public_api.Profile{
		Name:         prof.Name,
		Gender:       genderToProto(prof.Gender),
		Age:          int32(prof.Age),
		Bio:          prof.Bio,
		LocationLat:  prof.LocationLat,
		LocationLon:  prof.LocationLon,
		LocationName: prof.LocationName,
	}
}

func photosToLinkList(photos []models.PhotoShowcase) []string {
	res := make([]string, len(photos))
	for i, photo := range photos {
		res[i] = photo.Link
	}
	return res
}

func photosToProto(photos []models.PhotoShowcase) []*public_api.Photo {
	res := make([]*public_api.Photo, len(photos))
	for i, photo := range photos {
		res[i] = photoToProto(photo)
	}
	return res
}

func photoToProto(photo models.PhotoShowcase) *public_api.Photo {
	return &public_api.Photo{
		Key:  photo.Key,
		Link: photo.Link,
	}
}

func genderToProto(gender models.Gender) public_api.GENDER {
	switch gender {
	case models.GenderFemale:
		return public_api.GENDER_FEMALE
	default:
		return public_api.GENDER_MAlE
	}
}

func protoToGender(gender public_api.GENDER) models.Gender {
	switch gender {
	case public_api.GENDER_FEMALE:
		return models.GenderFemale
	default:
		return models.GenderMale
	}
}

func preferencesToProto(pref models.Preferences) *public_api.Preferences {
	return &public_api.Preferences{
		MaxAge:           int32(pref.MaxAge),
		MinAge:           int32(pref.MinAge),
		Gender:           genderToProto(pref.Gender),
		LocationLat:      pref.LocationLat,
		LocationLon:      pref.LocationLon,
		LocationRadiusKm: pref.LocationRadiusKm,
	}
}

func protoToPreferences(pref *public_api.Preferences) models.Preferences {
	if pref == nil {
		return models.Preferences{}
	}
	return models.Preferences{
		MaxAge:           int(pref.MaxAge),
		MinAge:           int(pref.MinAge),
		Gender:           protoToGender(pref.Gender),
		LocationLat:      pref.LocationLat,
		LocationLon:      pref.LocationLon,
		LocationRadiusKm: pref.LocationRadiusKm,
	}
}

func protoToProfile(prof *public_api.Profile) models.Profile {
	if prof == nil {
		return models.Profile{}
	}
	return models.Profile{
		Name:         prof.Name,
		Gender:       protoToGender(prof.Gender),
		Age:          int(prof.Age),
		Bio:          prof.Bio,
		LocationLat:  prof.LocationLat,
		LocationLon:  prof.LocationLon,
		LocationName: prof.LocationName,
	}
}

func protoToSwipeVerdict(sv public_api.SWIPE_VERDICT) models.SwipeVerdict {
	switch sv {
	case public_api.SWIPE_VERDICT_SWIPE_LIKE:
		return models.SwipeVerdictLike
	default:
		return models.SwipeVerdictDislike
	}
}

func protoToMsgContentType(ct public_api.MESSAGE_CONTENT_TYPE) models.MsgContentType {
	switch ct {
	case public_api.MESSAGE_CONTENT_TYPE_PHOTO:
		return models.ContentPhoto
	case public_api.MESSAGE_CONTENT_TYPE_VOICE:
		return models.ContentVoice
	default:
		return models.ContentText
	}
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

func chatsToProto(chats []models.ChatShowcase) []*public_api.Chat {
	res := make([]*public_api.Chat, len(chats))
	for i, chat := range chats {
		res[i] = chatToProto(chat)
	}
	return res
}

func chatToProto(chat models.ChatShowcase) *public_api.Chat {
	return &public_api.Chat{
		ChatId: chat.ID,
		Name:   chat.Name,
		Photo:  chat.Photo,
	}
}

func messagesToProto(messages []models.MessageShowcase) []*public_api.Message {
	res := make([]*public_api.Message, len(messages))
	for i, message := range messages {
		res[i] = messageToProto(message)
	}
	return res
}

func messageToProto(message models.MessageShowcase) *public_api.Message {
	return &public_api.Message{
		Id:          message.ID,
		SentByMe:    message.SentByMe,
		ContentType: msgContentTypeToProto(message.ContentType),
		Payload:     message.Payload,
		CreatedAt:   timestamppb.New(message.CreatedAt),
	}
}
