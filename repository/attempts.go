package repository

import (
	"pinder/models"
	"time"
)

func (r *Repository) GetLatestPairAttempt(user1, user2 uint64) (models.PairAttempt, error) {
	var pair models.PairAttempt
	res := r.db.Model(&models.PairAttempt{}).
		Where("user1 = ? and user2 = ?", user1, user2).
		Order("created_at desc").First(&pair)
	if res.Error != nil {
		return models.PairAttempt{}, res.Error
	}
	return pair, nil
}

func (r *Repository) GetPendingPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error) {
	var pair models.PairAttempt
	res := r.db.Model(&models.PairAttempt{}).
		Where("((user1 = ? and user2 = ?) or (user2 = ? and user1 = ?)) and state = ?",
			user1, user2, user1, user2, models.PAStatePending).
		First(&pair)
	if res.Error != nil {
		return models.PairAttempt{}, res.Error
	}
	return pair, nil
}

func (r *Repository) GetLatestPairAttemptByUserPair(user1, user2 uint64) (models.PairAttempt, error) {
	var pair models.PairAttempt
	res := r.db.Model(&models.PairAttempt{}).
		Where("((user1 = ? and user2 = ?) or (user2 = ? and user1 = ?))",
			user1, user2, user1, user2).
		Order("created_at desc").
		First(&pair)
	if res.Error != nil {
		return models.PairAttempt{}, res.Error
	}
	return pair, nil
}

func (r *Repository) FinishPairAttempt(PAID uint64, state models.PAState) error {
	res := r.db.Model(&models.PairAttempt{}).Where("id = ?", PAID).Update("state = ?", state)
	return res.Error
}

func (r *Repository) CreatePairAttempt(user1, user2 uint64) error {
	res := r.db.Create(&models.PairAttempt{
		User1:     user1,
		User2:     user2,
		State:     models.PAStatePending,
		CreatedAt: time.Now(),
	})
	return res.Error
}
