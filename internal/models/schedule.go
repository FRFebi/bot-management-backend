package models

import (
	"time"

	"gorm.io/gorm"
)

type Schedule struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	BotID          uint           `gorm:"not null;index" json:"bot_id"`
	CronExpression string         `gorm:"type:varchar(100);not null" json:"cron_expression"`
	IsActive       bool           `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	Bot *Bot `gorm:"foreignKey:BotID;constraint:OnDelete:CASCADE" json:"bot,omitempty"`
}

func (Schedule) TableName() string {
	return "schedules"
}