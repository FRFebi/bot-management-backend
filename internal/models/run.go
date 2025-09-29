package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Run struct {
	ID         uint           `gorm:"primarykey" json:"id"`
	BotID      uint           `gorm:"not null;index" json:"bot_id"`
	StartedAt  time.Time      `gorm:"default:CURRENT_TIMESTAMP" json:"started_at"`
	FinishedAt *time.Time     `json:"finished_at"`
	Success    *bool          `json:"success"`
	Log        string         `gorm:"type:text" json:"log"`
	Metrics    datatypes.JSON `gorm:"type:jsonb" json:"metrics"`
	CreatedAt  time.Time      `json:"created_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	Bot *Bot `gorm:"foreignKey:BotID;constraint:OnDelete:CASCADE" json:"bot,omitempty"`
}

func (Run) TableName() string {
	return "runs"
}