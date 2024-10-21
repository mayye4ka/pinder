package repository

import (
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type Repository struct {
	db     *gorm.DB
	logger *zerolog.Logger
}

func New(db *gorm.DB, logger *zerolog.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: logger,
	}
}
