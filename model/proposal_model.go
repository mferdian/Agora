package model

import "github.com/google/uuid"

type Proposal struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey" json:"proposal_id"`
	Title       string    `json:"proposal_name"`
	Description string    `json:"proposal_description"`

	UserID uuid.UUID `gorm:"type:uuid"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	TimeStamp
}
