package repo

import (
	"time"

	"gorm.io/gorm"
)

// TODO: separate transport, domain and persistence models
// sharing a model exposes sensitive fields and leads to tight coupling (hard to change)

type User struct {
	ID   uint   `gorm:"primarykey"`
	Name string `binding:"required,min=3,max=64" json:"name"`
}

type Trainer struct {
	// Note: I wouldn't separate user/trainer but the appointment data had user_id=1, trainer_id=1
	ID   uint   `gorm:"primarykey"`
	Name string `binding:"required,min=3,max=64" json:"name"`
}

// TODO enforce uniqueness of appointments at storage layer (user+trainer+start+end)
type Appointment struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	StartsAt  time.Time `json:"starts_at"`
	EndsAt    time.Time `json:"ends_at"`
	UserID    uint      `json:"user_id"`
	TrainerID uint      `json:"trainer_id"`

	User    User    `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:UserID;references:ID" json:"user"`
	Trainer Trainer `gorm:"constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;foreignKey:TrainerID;references:ID" json:"trainer"`
}

func (u *Appointment) BeforeSave(tx *gorm.DB) error {
	u.StartsAt = u.StartsAt.UTC()
	u.EndsAt = u.EndsAt.UTC()
	return nil
}
