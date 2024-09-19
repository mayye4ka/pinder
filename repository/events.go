package repository

import (
	"pinder/models"
	"time"
)

func (r *Repository) CreateEvent(PAID uint64, eventType models.PEType) error {
	pairEvent := models.PairEvent{
		PAID:      PAID,
		CreatedAt: time.Now(),
		EventType: eventType,
	}
	res := r.db.Create(&pairEvent)
	return res.Error
}

func (r *Repository) GetLastEvent(PAID uint64) (models.PairEvent, error) {
	var e models.PairEvent
	res := r.db.Model(&models.PairEvent{}).Where("paid = ?", PAID).Order("created_at desc").First(&e)
	if res.Error != nil {
		return models.PairEvent{}, res.Error
	}
	return e, nil
}
