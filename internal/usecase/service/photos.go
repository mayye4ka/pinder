package service

import (
	"context"
	"reflect"

	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"
	"github.com/pkg/errors"
)

func (s *Service) getUserPhotos(ctx context.Context, userId uint64) ([]models.PhotoShowcase, error) {
	photos, err := s.repository.GetUserPhotos(ctx, userId)
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
	return s.repository.AddPhoto(ctx, userId, key)
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
	return s.repository.DeleteUserPhoto(ctx, userId, photoKey)
}

func (s *Service) ReorderPhotos(ctx context.Context, newOrder []string) error {
	userId := ctx.Value(userIdContextKey).(uint64)
	if userId == 0 {
		return errUnauthenticated
	}
	photos, err := s.repository.GetUserPhotos(ctx, userId)
	if err != nil {
		return errors.Wrap(err, "can't reorder photos")
	}

	exPhotoMap := map[string]bool{}
	newPhotoMap := map[string]bool{}
	for _, p := range photos {
		newPhotoMap[p] = true
	}
	for _, p := range newOrder {
		exPhotoMap[p] = true
	}
	if !reflect.DeepEqual(newPhotoMap, exPhotoMap) {
		return &errs.CodableError{
			Code:    errs.CodeInvalidInput,
			Message: "should specify all photos to reorder",
		}
	}
	err = s.repository.ReorderPhotos(ctx, newOrder)
	if err != nil {
		return errors.Wrap(err, "can't reorder photos")
	}
	return nil
}
