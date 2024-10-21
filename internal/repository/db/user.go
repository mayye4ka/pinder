package repository

import (
	"errors"

	"github.com/mayye4ka/pinder/internal/errs"
	"github.com/mayye4ka/pinder/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Gender string

const (
	GenderMale   Gender = "male"
	GenderFemale Gender = "female"
)

type User struct {
	ID          uint64
	PhoneNumber string
	PassHash    string
}

func (User) TableName() string {
	return "users"
}

type Profile struct {
	UserID       uint64
	Name         string
	Gender       Gender
	Age          int
	Bio          string
	LocationLat  float64
	LocationLon  float64
	LocationName string
}

func (Profile) TableName() string {
	return "profiles"
}

type Preferences struct {
	UserID           uint64
	MaxAge           int
	MinAge           int
	Gender           Gender
	LocationLat      float64
	LocationLon      float64
	LocationRadiusKm float64
}

func (Preferences) TableName() string {
	return "preferences"
}

func (r *Repository) CreateUser(phoneNumber, passHash string) (models.User, error) {
	user := User{
		PhoneNumber: phoneNumber,
		PassHash:    passHash,
	}
	res := r.db.Create(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrDuplicatedKey) {
			return models.User{}, &errs.CodableError{
				Code:    errs.CodeInvalidInput,
				Message: "user with same phone number already exists",
			}
		}
		r.logger.Err(res.Error).Msg("can't create user")
		return models.User{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't create user",
		}
	}
	return mapUser(user), nil
}

func (r *Repository) GetUserByCreds(phoneNumber, passHash string) (models.User, error) {
	var user User
	res := r.db.Model(&User{}).Where("phone_number = ? and pass_hash = ?", phoneNumber, passHash).First(&user)
	if res.Error != nil {
		if errors.Is(res.Error, gorm.ErrRecordNotFound) {
			return models.User{}, &errs.CodableError{
				Code:    errs.CodeNotFound,
				Message: "invalid phone / password",
			}
		}
		r.logger.Err(res.Error).Msg("can't get user")
		return models.User{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get user",
		}
	}
	return mapUser(user), nil
}

func (r *Repository) GetProfile(userID uint64) (models.Profile, error) {
	var profile Profile
	res := r.db.Model(&Profile{}).Where("user_id=?", userID).First(&profile)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return models.Profile{}, nil
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get profile")
		return models.Profile{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get profile",
		}
	}
	return mapProfile(profile), nil
}

func (r *Repository) PutProfile(profile models.Profile) error {
	prof := unmapProfile(profile)
	res := r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&prof)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't put profile")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't update profile",
		}
	}
	return nil
}

func (r *Repository) GetPreferences(userID uint64) (models.Preferences, error) {
	var preferences Preferences
	res := r.db.Model(&Preferences{}).Where("user_id=?", userID).First(&preferences)
	if res.Error != nil && errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return models.Preferences{}, nil
	} else if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get preferences")
		return models.Preferences{}, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get preferences",
		}
	}
	return mapPreferences(preferences), nil
}

func (r *Repository) PutPreferences(preferences models.Preferences) error {
	prefs := unmapPreferences(preferences)
	res := r.db.Clauses(clause.OnConflict{UpdateAll: true}).Create(&prefs)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't put preferences")
		return &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't update preferences",
		}
	}
	return nil
}

func (r *Repository) GetAllValidUsers() ([]uint64, error) {
	var profiles []Profile
	res := r.db.Model(&Profile{}).Find(&profiles)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get all profiles")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get all profiles",
		}
	}
	var prefs []Preferences
	res = r.db.Model(&Preferences{}).Find(&prefs)
	if res.Error != nil {
		r.logger.Err(res.Error).Msg("can't get all preferences")
		return nil, &errs.CodableError{
			Code:    errs.CodeInternal,
			Message: "can't get all preferences",
		}
	}
	hasProfile := map[uint64]bool{}
	for _, p := range profiles {
		hasProfile[p.UserID] = true
	}
	ids := []uint64{}
	for _, p := range prefs {
		if hasProfile[p.UserID] {
			ids = append(ids, p.UserID)
		}
	}
	return ids, nil
}

func mapUser(user User) models.User {
	return models.User{
		ID:          user.ID,
		PhoneNumber: user.PhoneNumber,
		PassHash:    user.PassHash,
	}
}

func mapProfile(prof Profile) models.Profile {
	return models.Profile{
		UserID:       prof.UserID,
		Name:         prof.Name,
		Gender:       mapGender(prof.Gender),
		Age:          prof.Age,
		Bio:          prof.Bio,
		LocationLat:  prof.LocationLat,
		LocationLon:  prof.LocationLon,
		LocationName: prof.LocationName,
	}
}

func unmapProfile(prof models.Profile) Profile {
	return Profile{
		UserID:       prof.UserID,
		Name:         prof.Name,
		Gender:       unmapGender(prof.Gender),
		Age:          prof.Age,
		Bio:          prof.Bio,
		LocationLat:  prof.LocationLat,
		LocationLon:  prof.LocationLon,
		LocationName: prof.LocationName,
	}
}

func mapPreferences(pref Preferences) models.Preferences {
	return models.Preferences{
		UserID:           pref.UserID,
		MaxAge:           pref.MaxAge,
		MinAge:           pref.MinAge,
		Gender:           mapGender(pref.Gender),
		LocationLat:      pref.LocationLat,
		LocationLon:      pref.LocationLon,
		LocationRadiusKm: pref.LocationRadiusKm,
	}
}

func unmapPreferences(pref models.Preferences) Preferences {
	return Preferences{
		UserID:           pref.UserID,
		MaxAge:           pref.MaxAge,
		MinAge:           pref.MinAge,
		Gender:           unmapGender(pref.Gender),
		LocationLat:      pref.LocationLat,
		LocationLon:      pref.LocationLon,
		LocationRadiusKm: pref.LocationRadiusKm,
	}
}

func mapGender(g Gender) models.Gender {
	switch g {
	case GenderFemale:
		return models.GenderFemale
	default:
		return models.GenderMale
	}
}

func unmapGender(g models.Gender) Gender {
	switch g {
	case models.GenderFemale:
		return GenderFemale
	default:
		return GenderMale
	}
}
