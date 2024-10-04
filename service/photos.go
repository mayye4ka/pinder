package service

import (
	"context"

	"github.com/mayye4ka/pinder/models"
)

func (s *Service) getUserPhotos(ctx context.Context, userId uint64) ([]models.PhotoShowcase, error) {
	photos, err := s.repository.GetUserPhotos(userId)
	if err != nil {
		return nil, err
	}
	res := make([]models.PhotoShowcase, len(photos))
	for i, photo := range photos {
		link, err := s.filestorage.MakeProfilePhotoLink(ctx, photo)
		if err != nil {
			return nil, err
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
		return err
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
		return err
	}
	return s.repository.DeleteUserPhoto(userId, photoKey)
}
