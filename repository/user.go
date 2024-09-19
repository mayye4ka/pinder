package repository

import (
	"database/sql"
	"errors"
	"pinder/models"

	"gorm.io/gorm/clause"
)

func (r *Repository) CreateUser(phoneNumber, passHash string) (*models.User, error) {
	user := models.User{
		PhoneNumber: phoneNumber,
		PassHash:    passHash,
	}
	res := r.db.Create(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	// TODO: fill id for user
	return &user, nil
}

func (r *Repository) GetUserByCreds(phoneNumber, passHash string) (*models.User, error) {
	var user models.User
	res := r.db.Model(&models.User{}).Where("phone_number = ? and pass_hash = ?", phoneNumber, passHash).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (r *Repository) GetUser(id uint64) (*models.User, error) {
	var user models.User
	res := r.db.Model(&models.User{}).Where("id = ?", id).First(&user)
	if res.Error != nil {
		return nil, res.Error
	}
	return &user, nil
}

func (r *Repository) GetProfile(userID uint64) (*models.Profile, error) {
	var profile models.Profile
	res := r.db.Model(&models.Profile{}).Where("user_id=?", userID).First(&profile)
	if res.Error != nil && errors.Is(res.Error, sql.ErrNoRows) {
		return nil, nil
	} else if res.Error != nil {
		return nil, res.Error
	}
	return &profile, nil
}

func (r *Repository) PutProfile(profile models.Profile) error {
	res := r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&profile)
	return res.Error
}

func (r *Repository) GetPreferences(userID uint64) (*models.Preferences, error) {
	var preferences models.Preferences
	res := r.db.Model(&models.Preferences{}).Where("user_id=?", userID).First(&preferences)
	if res.Error != nil && errors.Is(res.Error, sql.ErrNoRows) {
		return nil, nil
	} else if res.Error != nil {
		return nil, res.Error
	}
	return &preferences, nil
}

func (r *Repository) PutPreferences(preferences models.Preferences) error {
	res := r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&preferences)
	return res.Error
}

type Candidate struct {
	ID          uint64
	Profile     models.Profile
	Preferences models.Preferences
}

func (r *Repository) GetAllCandidates() ([]Candidate, error) {
	var profiles []models.Profile
	res := r.db.Model(&models.Profile{}).Find(&profiles)
	if res.Error != nil {
		return nil, res.Error
	}
	var prefs []models.Preferences
	res = r.db.Model(&models.Preferences{}).Find(&prefs)
	if res.Error != nil {
		return nil, res.Error
	}
	profMap := map[uint64]models.Profile{}
	for _, p := range profiles {
		profMap[p.UserID] = p
	}
	cands := []Candidate{}
	for _, pref := range prefs {
		if prof, ok := profMap[pref.UserID]; ok {
			cands = append(cands, Candidate{
				ID:          pref.UserID,
				Profile:     prof,
				Preferences: pref,
			})
		}
	}
	return cands, nil
}
