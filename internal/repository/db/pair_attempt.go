package repository

import (
	"errors"
	"time"

	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
	"gorm.io/gorm"
)

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

func (r *Repository) GetLatestPairAttempt(user1, user2 uint64) (models.PairAttempt, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("user1 = ? and user2 = ?", user1, user2).
		Order("created_at desc").First(&pair)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't latest pair attempt")
		return models.PairAttempt{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get latest pair attempt",
		}
	}
	return mapPairAttempt(pair), nil
}

func (r *Repository) GetPendingPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("((user1 = ? and user2 = ?) or (user2 = ? and user1 = ?)) and state = ?",
			user1, user2, user1, user2, models.PAStatePending).
		First(&pair)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get pending pa for this pair")
		return models.PairAttempt{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get pending pa for this pair",
		}
	}
	return mapPairAttempt(pair), nil
}

func (r *Repository) GetLatestPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("((user1 = ? and user2 = ?) or (user2 = ? and user1 = ?))",
			user1, user2, user1, user2).
		Order("created_at desc").
		First(&pair)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get latest pa for this pair")
		return models.PairAttempt{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get latest pa for this pair",
		}
	}
	return mapPairAttempt(pair), nil
}

func (r *Repository) FinishPairAttempt(PAID uint64, state models.PAState) error {
	res := r.db.Model(&PairAttempt{}).Where("id = ?", PAID).Update("state", unmapPaState(state))
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't finish pair attempt")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't finish pair attempt",
		}
	}
	return nil
}

func (r *Repository) CreatePairAttempt(user1, user2 uint64) (models.PairAttempt, error) {
	pa := PairAttempt{
		User1:     user1,
		User2:     user2,
		State:     PAStatePending,
		CreatedAt: time.Now(),
	}
	res := r.db.Create(&pa)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't create pair attempt")
		return models.PairAttempt{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't create pair attempt",
		}
	}
	return mapPairAttempt(pa), nil
}

func (r *Repository) GetWhoLikedMe(userID uint64) (uint64, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("user2 = ? and state = ?", userID, PAStatePending).
		Order("created_at").First(&pair)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		r.logger.Err(res.Error).Msg("can't get who likes you")
		return 0, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get who likes you",
		}
	}
	return pair.User1, nil
}

func (r *Repository) GetPendingPairAttempts(user1ID uint64) ([]models.PairAttempt, error) {
	var pas []PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("user1 = ? and state = ?", user1ID, PAStatePending).
		Find(&pas)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get pending pair attempts")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get pending pair attempts",
		}
	}
	return mapPairAttempts(pas), nil
}

func mapPairAttempts(pas []PairAttempt) []models.PairAttempt {
	res := make([]models.PairAttempt, len(pas))
	for i, pa := range pas {
		res[i] = mapPairAttempt(pa)
	}
	return res
}

func mapPairAttempt(pa PairAttempt) models.PairAttempt {
	return models.PairAttempt{
		ID:        pa.ID,
		User1:     pa.User1,
		User2:     pa.User2,
		State:     mapPaState(pa.State),
		CreatedAt: pa.CreatedAt,
	}
}

func mapPaState(s PAState) models.PAState {
	switch s {
	case PAStateMatch:
		return models.PAStateMatch
	case PAStateMismatch:
		return models.PAStateMismatch
	default:
		return models.PAStatePending
	}
}

func unmapPaState(s models.PAState) PAState {
	switch s {
	case models.PAStateMatch:
		return PAStateMatch
	case models.PAStateMismatch:
		return PAStateMismatch
	default:
		return PAStatePending
	}
}
