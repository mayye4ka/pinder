package service

import (
	"context"

	"github.com/mayye4ka/pinder/models"
	"github.com/pkg/errors"
)

func (s *Service) getUserPhotos(ctx context.Context, userId uint64) ([]models.PhotoShowcase, error) {
	photos, err := s.repository.GetUserPhotos(userId)
	if err != nil {
		return nil, errors.Wrap(err, "can't get user photos")
	}
	res := make([]models.PhotoShowcase, len(photos))
	for i, photo := range photos {
		link, err := s.filestorage.MakeProfilePhotoLink(ctx, photo)
		if err != nil {
			return nil, errors.Wrap(err, "can't make profile photo links")
		}
		res[i] = models.PhotoShowcase{
			Key:  photo,
			Link: link,
		}
	}
	return res, nil
}

func (s *Service) AddPhoto(ctx context.Context, photo string) error {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return errUnauthenticated
	}
	key, err := s.filestorage.SaveProfilePhoto(ctx, []byte(photo))
	if err != nil {
		return errors.Wrap(err, "can't save profile photo")
	}
	return s.repository.AddPhoto(userId, key)
}

func (s *Service) DeletePhoto(ctx context.Context, photoKey string) error {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return errUnauthenticated
	}
	err := s.filestorage.DelProfilePhoto(ctx, photoKey)
	if err != nil {
		return errors.Wrap(err, "can't delete profile photo")
	}
	return s.repository.DeleteUserPhoto(userId, photoKey)
}
