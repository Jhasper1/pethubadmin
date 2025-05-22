package models

import "time"

type SubmittedReport struct {
	ID          uint      `gorm:"primaryKey;column:id" json:"id"`
	ShelterID   uint      `gorm:"column:shelter_id" json:"shelter_id"`
	AdopterID   uint      `gorm:"column:adopter_id" json:"adopter_id"`
	Reason      string    `gorm:"type:text;column:reason" json:"reason"`
	Description string    `gorm:"type:text;column:description" json:"description"`
	Status      string    `gorm:"type:text;column:status;default:'pending'" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;autoCreateTime" json:"created_at"`

	Shelter ShelterInfo `gorm:"foreignKey:ShelterID;references:ShelterID" json:"shelter"`
	Adopter AdopterInfo `gorm:"foreignKey:AdopterID;references:AdopterID" json:"adopter"`
}

func (SubmittedReport) TableName() string {
	return "submittedreports"
}
