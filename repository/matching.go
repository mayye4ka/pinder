package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"sort"

	"github.com/mayye4ka/pinder/models"
)

func (r *Repository) getProfilePtr(userId uint64) (*models.Profile, error) {
	prof, err := r.GetProfile(userId)
	if err != nil {
		return nil, err
	}
	return &prof, nil
}

// TODO: should be transfered to service layer
func (r *Repository) GetWhoLikedMe(userID uint64) (*models.Profile, error) {
	var pair PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("user2 = ? and state = ?", userID, PAStatePending).
		Order("created_at").First(&pair)
	if res.Error != nil && !errors.Is(res.Error, sql.ErrNoRows) {
		return nil, res.Error
	} else if res.Error != nil {
		return nil, nil
	}
	return r.getProfilePtr(pair.User1)
}

func (r *Repository) GetHangingPartner(userID uint64) (*models.Profile, error) {
	var pas []PairAttempt
	res := r.db.Model(&PairAttempt{}).
		Where("user1 = ? and state = ?", userID, PAStatePending).
		Find(&pas)
	if res.Error != nil {
		return nil, res.Error
	}
	for _, pa := range pas {
		le, err := r.GetLastEvent(pa.ID)
		if err != nil {
			return nil, err
		}
		if le.EventType == models.PETypeSentToUser1 {
			prof, err := r.GetProfile(pa.User2)
			if err != nil {
				return nil, err
			}
			return &prof, nil
		}
	}
	return nil, nil
}

func (r *Repository) ChooseCandidateAndCreatePairAttempt(userID uint64) (*models.Profile, error) {
	prof, err := r.GetProfile(userID)
	if err != nil {
		return nil, err
	}
	pref, err := r.GetPreferences(userID)
	if err != nil {
		return nil, err
	}
	if prof.UserID == 0 || pref.UserID == 0 {
		return nil, fmt.Errorf("not ready for search")
	}
	cands, err := r.getAllCandidates()
	if err != nil {
		return nil, err
	}
	candidatePAs := []models.PairAttempt{}
	for _, candidate := range cands {
		if !pref.ProfileMatches(candidate.Profile) || !candidate.Preferences.ProfileMatches(prof) {
			continue
		}
		_, err := r.GetPendingPairAttemptByUserPair(candidate.ID, userID)
		if err == nil || !errors.Is(err, sql.ErrNoRows) {
			continue
		}
		lpa, err := r.GetLatestPairAttempt(candidate.ID, userID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			return r.getProfilePtr(candidate.ID)
		} else if err != nil {
			return nil, err
		}
		candidatePAs = append(candidatePAs, lpa)
	}
	if len(candidatePAs) == 0 {
		return nil, nil
	}
	sort.Slice(candidatePAs, func(i, j int) bool {
		return candidatePAs[i].CreatedAt.Before(candidatePAs[j].CreatedAt)
	})
	bestCandidate := candidatePAs[0].User1
	if bestCandidate == userID {
		bestCandidate = candidatePAs[0].User2
	}
	return r.getProfilePtr(bestCandidate)
}
