package service

import (
	"context"
	"math/rand"
	"sort"
	"time"

	"github.com/mayye4ka/pinder/errs"
	"github.com/mayye4ka/pinder/models"
)

func (s *Service) NextPartner(ctx context.Context) (models.ProfileShowcase, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return models.ProfileShowcase{}, errUnauthenticated
	}

	myProfile, err := s.repository.GetProfile(userId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	myPrefs, err := s.repository.GetPreferences(userId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	if myProfile.UserID == 0 || myPrefs.UserID == 0 {
		return models.ProfileShowcase{}, &errs.CodableError{
			Code:    errs.CodeInvalidInput,
			Message: "incomplete profile",
		}
	}

	partner, err := s.submitHangingPartner(ctx, userId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	if partner != nil {
		return *partner, nil
	}

	partner, err = s.submitWhoLikedMe(ctx, userId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	if partner != nil {
		return *partner, nil
	}

	partner, err = s.chooseCandidateAndCreateNewPair(ctx, userId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	if partner == nil {
		return models.ProfileShowcase{}, &errs.CodableError{
			Code:    errs.CodeNotFound,
			Message: "lower your expectations to zero",
		}
	}
	return *partner, nil
}

func (s *Service) submitHangingPartner(ctx context.Context, userID uint64) (*models.ProfileShowcase, error) {
	hp, err := s.getHangingPartner(userID)
	if err != nil {
		return nil, err
	}
	if hp == 0 {
		return nil, nil
	}
	prof, err := s.createProfileShowcase(ctx, hp)
	if err != nil {
		return nil, err
	}
	return &prof, nil
}

func (s *Service) submitWhoLikedMe(ctx context.Context, userID uint64) (*models.ProfileShowcase, error) {
	liker, err := s.repository.GetWhoLikedMe(userID)
	if err != nil {
		return nil, err
	}
	if liker == 0 {
		return nil, nil
	}
	pa, err := s.repository.GetLatestPairAttempt(liker, userID)
	if err != nil {
		return nil, err
	}
	err = s.repository.CreateEvent(pa.ID, models.PETypeSentToUser2)
	if err != nil {
		return nil, err
	}
	prof, err := s.createProfileShowcase(ctx, liker)
	if err != nil {
		return nil, err
	}
	return &prof, nil
}

func (s *Service) createProfileShowcase(ctx context.Context, candidateId uint64) (models.ProfileShowcase, error) {
	profile, err := s.repository.GetProfile(candidateId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	photos, err := s.getUserPhotos(ctx, candidateId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	return models.ProfileShowcase{
		Profile: profile,
		Photos:  photos,
	}, nil
}

func (s *Service) getHangingPartner(userID uint64) (uint64, error) {
	pas, err := s.repository.GetPendingPairAttempts(userID)
	if err != nil {
		return 0, err
	}
	for _, pa := range pas {
		le, err := s.repository.GetLastEvent(pa.ID)
		if err != nil {
			return 0, err
		}
		if le.EventType == models.PETypeSentToUser1 {
			return pa.User2, nil
		}
	}
	return 0, nil
}

func (s *Service) chooseCandidateAndCreateNewPair(ctx context.Context, userId uint64) (*models.ProfileShowcase, error) {
	candidate, err := s.chooseCandidate(userId)
	if err != nil {
		return nil, err
	}
	if candidate == 0 {
		return nil, nil
	}
	pa, err := s.repository.CreatePairAttempt(userId, candidate)
	if err != nil {
		return nil, err
	}
	err = s.repository.CreateEvent(pa.ID, models.PETypePACreated)
	if err != nil {
		return nil, err
	}
	err = s.repository.CreateEvent(pa.ID, models.PETypeSentToUser1)
	if err != nil {
		return nil, err
	}
	prof, err := s.createProfileShowcase(ctx, candidate)
	if err != nil {
		return nil, err
	}
	return &prof, nil
}

func (s *Service) chooseCandidate(userId uint64) (uint64, error) {
	ids, err := s.repository.GetAllValidUsers()
	if err != nil {
		return 0, err
	}
	myPref, err := s.repository.GetPreferences(userId)
	if err != nil {
		return 0, err
	}
	myProf, err := s.repository.GetProfile(userId)
	if err != nil {
		return 0, err
	}
	candidates := []uint64{}
	for _, id := range ids {
		pref, err := s.repository.GetPreferences(id)
		if err != nil {
			return 0, err
		}
		prof, err := s.repository.GetProfile(id)
		if err != nil {
			return 0, err
		}
		if !myPref.ProfileMatches(prof) || !pref.ProfileMatches(myProf) {
			continue
		}
		pa, _ := s.repository.GetPendingPairAttemptByUserPair(userId, id)
		if pa.ID != 0 {
			continue
		}
		candidates = append(candidates, id)
	}
	noShows := []uint64{}
	withLike := []uint64{}
	withDislike := []uint64{}
	latestPaTime := map[uint64]time.Time{}
	for _, id := range candidates {
		pa, _ := s.repository.GetLatestPairAttemptByUserPair(userId, id)
		if pa.ID == 0 {
			noShows = append(noShows, id)
			continue
		}
		if pa.State == models.PAStateMatch {
			withLike = append(withLike, id)
		} else {
			withDislike = append(withDislike, id)
		}
		latestPaTime[getWhoIsNotMe(pa.User1, pa.User2, userId)] = pa.CreatedAt
	}
	if len(noShows) > 0 {
		return noShows[rand.Intn(len(noShows))], nil
	}
	if len(withDislike) > 0 {
		sort.Slice(withDislike, func(i, j int) bool {
			return latestPaTime[withDislike[i]].Before(latestPaTime[withDislike[j]])
		})
		return withDislike[0], nil
	}
	if len(withLike) > 0 {
		sort.Slice(withLike, func(i, j int) bool {
			return latestPaTime[withLike[i]].Before(latestPaTime[withLike[j]])
		})
		return withLike[0], nil
	}
	return 0, nil
}
