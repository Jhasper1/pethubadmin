package models

import "time"

// AdopterAccount model (linked to existing "adopteraccount" table)
type AdopterAccount struct {
	AdopterID uint   `gorm:"primaryKey" json:"adopter_id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `json:"password"`
	Status    string `gorm:"default:'active'" json:"status"` // Add this line
	CreatedAt time.Time
	//Info      AdopterInfo `gorm:"foreignKey:AdopterID;constraint:OnDelete:CASCADE" json:"info"`
}

// TableName overrides default table name
func (AdopterAccount) TableName() string {
	return "adopteraccount"
}

// AdopterInfo model (linked to existing "adopterinfo" table)
type AdopterInfo struct {
	AdopterID     uint   `gorm:"column:adopter_id;primaryKey" json:"adopter_id"`
	FirstName     string `gorm:"column:first_name" json:"first_name"`
	LastName      string `gorm:"column:last_name" json:"last_name"`
	Age           int    `gorm:"column:age" json:"age"`
	Sex           string `gorm:"column:sex" json:"sex"`
	Address       string `gorm:"column:address" json:"address"`
	ContactNumber string `gorm:"column:contact_number" json:"contact_number"`
	Email         string `gorm:"column:email" json:"email"`
	Occupation    string `gorm:"column:occupation" json:"occupation"`
	CivilStatus   string `gorm:"column:civil_status" json:"civil_status"`
	SocialMedia   string `gorm:"column:social_media" json:"social_media"`
}

func (AdopterInfo) TableName() string {
	return "adopterinfo"
}

type AdopterMedia struct {
	AdopterID      uint   `gorm:"primaryKey;autoIncrement:false" json:"adopter_id"`
	AdopterProfile string `json:"adopter_profile"`
}

func (AdoptedPet) TableName() string {
	return "adopterpets"
}

type AdoptedPet struct {
	AdoptedID uint `gorm:"column:adopted_id;primaryKey;autoIncrement" json:"adopted_id"`

	AdopterID uint `gorm:"column:adopter_id" json:"adopter_id"`

	PetID uint `gorm:"column:pet_id" json:"pet_id"`
}
