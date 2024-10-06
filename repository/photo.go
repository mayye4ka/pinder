package repository

type Photo struct {
	UserID   uint64
	PhotoKey string
}

func (Photo) TableName() string {
	return "photos"
}

func (r *Repository) AddPhoto(userID uint64, photoKey string) error {
	res := r.db.Create(&Photo{
		UserID:   userID,
		PhotoKey: photoKey,
	})
	return res.Error
}

func (r *Repository) GetUserPhotos(userID uint64) ([]string, error) {
	var photos []Photo
	res := r.db.Model(&Photo{}).Where("user_id = ?", userID).Find(&photos)
	if res.Error != nil {
		return nil, res.Error
	}
	result := make([]string, len(photos))
	for i, photo := range photos {
		result[i] = photo.PhotoKey
	}
	return result, nil
}

func (r *Repository) DeleteUserPhoto(userID uint64, photoKey string) error {
	res := r.db.Where("user_id = ? and photo_key = ?", userID, photoKey).Delete(&Photo{})
	return res.Error
}
