package repository

import (
	"context"
	"errors"

	"github.com/mayye4ka/pinder/internal/errs"
	"gorm.io/gorm"
)

type Photo struct {
	UserID   uint64
	PhotoKey string
	OrderN   int
}

func (Photo) TableName() string {
	return "photos"
}

func (r *Repository) getMaxPhotoOrder(ctx context.Context, userID uint64) (int, error) {
	var maxPh Photo
	res := r.db.WithContext(ctx).Model(&Photo{}).Where("user_id = ?", userID).Order("order_n desc").First(&maxPh)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return 0, nil
		}
		r.logger.Err(res.Error).Msg("can't get max photo order")
		return 0, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get max photo order",
		}
	}
	return maxPh.OrderN, nil
}

func (r *Repository) AddPhoto(ctx context.Context, userID uint64, photoKey string) error {
	maxOrder, err := r.getMaxPhotoOrder(ctx, userID)
	if err != nil {
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get photo order",
		}
	}
	res := r.db.WithContext(ctx).Create(&Photo{
		UserID:   userID,
		PhotoKey: photoKey,
		OrderN:   maxOrder + 1,
	})
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't create photo")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't create photo",
		}
	}
	return nil
}

func (r *Repository) GetUserPhotos(ctx context.Context, userID uint64) ([]string, error) {
	var photos []Photo
	res := r.db.WithContext(ctx).Model(&Photo{}).Where("user_id = ?", userID).Order("order_n").Find(&photos)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get user photos")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get user photos",
		}
	}
	result := make([]string, len(photos))
	for i, photo := range photos {
		result[i] = photo.PhotoKey
	}
	return result, nil
}

func (r *Repository) DeleteUserPhoto(ctx context.Context, userID uint64, photoKey string) error {
	res := r.db.WithContext(ctx).Where("user_id = ? and photo_key = ?", userID, photoKey).Delete(&Photo{})
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return &errs.CodableError{
				Code:    errs.CodeNotFound,
				Message: "photo not found",
			}
		}
		r.logger.Err(res.Error).Msg("can't delete user photo")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't delete user photo",
		}
	}
	return nil
}

func (r *Repository) updatePhotoOrder(ctx context.Context, photo string, newOrder int) error {
	res := r.db.WithContext(ctx).Model(&Photo{}).Where("photo_key = ?", photo).Update("order_n", newOrder)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't update photo order")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't update photo order",
		}
	}
	return nil
}

func (r *Repository) ReorderPhotos(ctx context.Context, photos []string) error {
	for i, photo := range photos {
		err := r.updatePhotoOrder(ctx, photo, i+1)
		if err != nil {
			return err
		}
	}
	return nil
}
