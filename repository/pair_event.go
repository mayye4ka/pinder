package repository

import (
	"errors"
	"time"

	"github.com/mayye4ka/pinder/errs"
	"github.com/mayye4ka/pinder/models"
	"gorm.io/gorm"
)

type PairEvent struct {
	ID        uint64
	PAID      uint64 `gorm:"column:pa_id"`
	CreatedAt time.Time
	EventType PEType
}

func (PairEvent) TableName() string {
	return "pair_events"
}

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

func (r *Repository) CreateEvent(PAID uint64, eventType models.PEType) error {
	pairEvent := PairEvent{
		PAID:      PAID,
		CreatedAt: time.Now(),
		EventType: unmapPeType(eventType),
	}
	res := r.db.Create(&pairEvent)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't create event")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't create event",
		}
	}
	return nil
}

func (r *Repository) GetLastEvent(PAID uint64) (models.PairEvent, error) {
	var e PairEvent
	res := r.db.Model(&PairEvent{}).Where("paid = ?", PAID).Order("created_at desc").First(&e)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return models.PairEvent{}, &errs.CodableError{
				Code:    errs.CodeNotFound,
				Message: "no events for this pair attempt",
			}
		}
		r.logger.Err(res.Error).Msg("can't get last event")
		return models.PairEvent{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get last event",
		}
	}
	return mapPairEvent(e), nil
}

func mapPairEvent(e PairEvent) models.PairEvent {
	return models.PairEvent{
		ID:        e.ID,
		PAID:      e.PAID,
		CreatedAt: e.CreatedAt,
		EventType: mapPeType(e.EventType),
	}
}

func mapPeType(pe PEType) models.PEType {
	switch pe {
	case PETypePACreated:
		return models.PETypePACreated
	case PETypeSentToUser1:
		return models.PETypeSentToUser1
	case PETypeUser1Liked:
		return models.PETypeUser1Liked
	case PETypeUser1Disliked:
		return models.PETypeUser1Disliked
	case PETypeSentToUser2:
		return models.PETypeSentToUser2
	case PETypeUser2Liked:
		return models.PETypeUser2Liked
	case PETypeUser2Disliked:
		return models.PETypeUser2Disliked
	case PETypePairAttemptFailed:
		return models.PETypePairAttemptFailed
	default:
		return models.PETypePairCreated
	}
}

func unmapPeType(pe models.PEType) PEType {
	switch pe {
	case models.PETypePACreated:
		return PETypePACreated
	case models.PETypeSentToUser1:
		return PETypeSentToUser1
	case models.PETypeUser1Liked:
		return PETypeUser1Liked
	case models.PETypeUser1Disliked:
		return PETypeUser1Disliked
	case models.PETypeSentToUser2:
		return PETypeSentToUser2
	case models.PETypeUser2Liked:
		return PETypeUser2Liked
	case models.PETypeUser2Disliked:
		return PETypeUser2Disliked
	case models.PETypePairAttemptFailed:
		return PETypePairAttemptFailed
	default:
		return PETypePairCreated
	}
}
