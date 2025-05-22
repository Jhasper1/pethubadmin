package models

import "time"

// ShelterAccount model (linked to existing "shelteraccount" table)
type ShelterAccount struct {
	ShelterID uint   `gorm:"primaryKey" json:"shelter_id"`
	Username  string `gorm:"unique;not null" json:"username"`
	Password  string `json:"password"`
	Status    string `gorm:"default:'active'" json:"status"` // Add this line
	RegStatus string `json:"reg_status"`
	CreatedAt time.Time
	//Info      ShelterInfo `gorm:"foreignKey:ShelterID;constraint:OnDelete:CASCADE" json:"info"`
}

// TableName overrides default table name
func (ShelterAccount) TableName() string {
	return "shelteraccount"
}

// ShelterInfo model (linked to existing "shelterinfo" table)

type ShelterInfo struct {
	ShelterID          uint   `gorm:"column:shelter_id;primaryKey;autoIncrement:false" json:"shelter_id"`
	ShelterName        string `gorm:"column:shelter_name" json:"shelter_name"`
	ShelterAddress     string `gorm:"column:shelter_address" json:"shelter_address"`
	ShelterLandmark    string `gorm:"column:shelter_landmark" json:"shelter_landmark"`
	ShelterContact     string `gorm:"column:shelter_contact" json:"shelter_contact"`
	ShelterEmail       string `gorm:"column:shelter_email" json:"shelter_email"`
	ShelterOwner       string `gorm:"column:shelter_owner" json:"shelter_owner"`
	ShelterDescription string `gorm:"column:shelter_description" json:"shelter_description"`
	ShelterSocial      string `gorm:"column:shelter_social" json:"shelter_social"`
}

func (ShelterInfo) TableName() string {
	return "shelterinfo"
}

type ShelterMedia struct {
	ShelterID      uint   `gorm:"primaryKey;autoIncrement:false" json:"shelter_id"`
	ShelterProfile string `json:"shelter_profile"`
	ShelterCover   string `json:"shelter_cover"`
}

func (ShelterMedia) TableName() string {
	return "sheltermedia"
}

type ShelterDonations struct {
	DonationID    uint   `gorm:"primaryKey;autoIncrement:true" json:"donation_id"`
	ShelterID     uint   `json:"shelter_id"`
	AccountNumber string `json:"account_number"`
	AccountName   string `json:"account_name"`
	QRImage       string `json:"qr_image"`
	CreatedAt     time.Time
}

func (ShelterDonations) TableName() string {
	return "shelterdonations"
}
