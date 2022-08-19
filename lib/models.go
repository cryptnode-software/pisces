package lib

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Model struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:uuid_generate_v4();" json:"id"`
	gorm.Model `json:"-"`
}

func (model *Model) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	model.ID = id
	return err
}
