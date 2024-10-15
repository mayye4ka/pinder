package repository

import (
	"errors"

	"github.com/mayye4ka/pinder/errs"
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

func (r *Repository) getMaxPhotoOrder(userID uint64) (int, error) {
	var max int
	res := r.db.Model(&Photo{}).Select("max(order_n)").Where("user_id = ?", userID).First(max)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get max photo order")
		return 0, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get max photo order",
		}
	}
	return max, nil
}

func (r *Repository) AddPhoto(userID uint64, photoKey string) error {
	maxOrder, err := r.getMaxPhotoOrder(userID)
	if err != nil {
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get photo order",
		}
	}
	res := r.db.Create(&Photo{
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

func (r *Repository) GetUserPhotos(userID uint64) ([]string, error) {
	var photos []Photo
	res := r.db.Model(&Photo{}).Where("user_id = ?", userID).Order("order_n").Find(&photos)
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

func (r *Repository) DeleteUserPhoto(userID uint64, photoKey string) error {
	res := r.db.Where("user_id = ? and photo_key = ?", userID, photoKey).Delete(&Photo{})
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

func (r *Repository) updatePhotoOrder(photo string, newOrder int) error {
	res := r.db.Model(&Photo{}).Where("photo_key = ?", photo).Update("order_n", newOrder)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't update photo order")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't update photo order",
		}
	}
	return nil
}

func (r *Repository) ReorderPhotos(photos []string) error {
	for i, photo := range photos {
		err := r.updatePhotoOrder(photo, i+1)
		if err != nil {
			return err
		}
	}
	return nil
}
