package service

import (
	"context"
	"fmt"

	"github.com/mayye4ka/pinder/models"
)

func (s *Service) UpdProfile(ctx context.Context, newProfile models.Profile) error {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return errUnauthenticated
	}
	newProfile.UserID = userId
	if newProfile.Gender != models.GenderFemale && newProfile.Gender != models.GenderMale {
		return fmt.Errorf("bad profile")
	}
	if newProfile.LocationName == "" || newProfile.Name == "" {
		return fmt.Errorf("bad profile")
	}
	if newProfile.Age == 0 || newProfile.LocationLat == 0 || newProfile.LocationLon == 0 {
		return fmt.Errorf("bad profile")
	}
	return s.repository.PutProfile(newProfile)
}

func (s *Service) GetProfile(ctx context.Context) (models.ProfileShowcase, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return models.ProfileShowcase{}, errUnauthenticated
	}
	profile, err := s.repository.GetProfile(userId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	photos, err := s.getUserPhotos(ctx, userId)
	if err != nil {
		return models.ProfileShowcase{}, err
	}
	return models.ProfileShowcase{
		Profile: profile,
		Photos:  photos,
	}, nil
}

func (s *Service) UpdPreferences(ctx context.Context, newPreferences models.Preferences) error {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return errUnauthenticated
	}
	newPreferences.UserID = userId
	return s.repository.PutPreferences(newPreferences)
}

func (s *Service) GetPreferences(ctx context.Context) (models.Preferences, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return models.Preferences{}, errUnauthenticated
	}
	preferences, err := s.repository.GetPreferences(userId)
	if err != nil {
		return models.Preferences{}, err
	}
	return preferences, nil
}
