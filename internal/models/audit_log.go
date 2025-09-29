package models

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    *uint          `gorm:"index" json:"user_id"`
	Action    string         `gorm:"type:varchar(255);not null" json:"action"`
	Details   datatypes.JSON `gorm:"type:jsonb" json:"details"`
	CreatedAt time.Time      `json:"created_at"`

	User *User `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL" json:"user,omitempty"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}