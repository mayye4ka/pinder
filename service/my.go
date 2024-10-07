package service

import (
	"context"
	"fmt"

	"github.com/mayye4ka/pinder/models"
	"github.com/pkg/errors"
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
	err := s.repository.PutProfile(newProfile)
	if err != nil {
		return errors.Wrap(err, "can't update profile")
	}
	return nil
}

func (s *Service) GetProfile(ctx context.Context) (models.ProfileShowcase, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return models.ProfileShowcase{}, errUnauthenticated
	}
	profile, err := s.repository.GetProfile(userId)
	if err != nil {
		return models.ProfileShowcase{}, errors.Wrap(err, "can't get profile")
	}
	photos, err := s.getUserPhotos(ctx, userId)
	if err != nil {
		return models.ProfileShowcase{}, errors.Wrap(err, "can't get user photos")
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
	err := s.repository.PutPreferences(newPreferences)
	if err != nil {
		return errors.Wrap(err, "can't update preferences")
	}
	return nil
}

func (s *Service) GetPreferences(ctx context.Context) (models.Preferences, error) {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return models.Preferences{}, errUnauthenticated
	}
	preferences, err := s.repository.GetPreferences(userId)
	if err != nil {
		return models.Preferences{}, errors.Wrap(err, "can't get preferences")
	}
	return preferences, nil
}
