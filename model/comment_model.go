package model

import "github.com/google/uuid"

type Comment struct {
	ID      uuid.UUID `gorm:"type:uuid;primaryKey" json:"comment_id"`
	Content string    `json:"comment_content"`

	UserID uuid.UUID `gorm:"type:uuid"`
	User   User      `gorm:"foreignKey:UserID;references:ID"`

	ProposalID uuid.UUID `gorm:"type:uuid"`
	Proposal   Proposal  `gorm:"foreignKey:ProposalID;references:ID"`

	TimeStamp
}
