package models

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Bot struct {
	ID          uint           `gorm:"primarykey" json:"id"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Version     string         `gorm:"type:varchar(50)" json:"version"`
	Config      datatypes.JSON `gorm:"type:jsonb" json:"config"`
	Status      string         `gorm:"type:varchar(20);default:'stopped'" json:"status"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	Schedules []Schedule `gorm:"foreignKey:BotID" json:"schedules,omitempty"`
	Runs      []Run      `gorm:"foreignKey:BotID" json:"runs,omitempty"`
}

func (Bot) TableName() string {
	return "bots"
}