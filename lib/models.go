package lib

// type Model struct {
// 	ID         uuid.UUID `gorm:"type:varchar(36);primary_key;default:(uuid());not null" json:"id"`
// 	gorm.Model `json:"-"`
// }

// func (model *Model) BeforeCreate(tx *gorm.DB) error {
// 	id, err := uuid.NewRandom()
// 	model.ID = id
// 	return err
// }

type SaveConditions struct {
	Root bool
}
