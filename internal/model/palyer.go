package model

import (
	nba_v "nba-remake/api/proto/v1"
	"time"
)

type Player struct {
	ID           uint32             `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	TeamID       uint32             `gorm:"not null;index;column:team_id" json:"team_id"`
	Name         string             `gorm:"type:varchar(100);not null;index;column:name" json:"name"`
	JerseyNumber uint8              `gorm:"type:tinyint unsigned;not null;column:jersey_number" json:"jersey_number"`
	Position     nba_v.Position     `gorm:"type:varchar(10);not null;column:position" json:"position"`
	Height       float64            `gorm:"type:decimal(3,2);column:height" json:"height,omitempty"`
	Weight       float64            `gorm:"type:decimal(5,2);column:weight" json:"weight,omitempty"`
	Birthday     *time.Time         `gorm:"type:date;column:birthday" json:"birthday,omitempty"`
	Status       nba_v.PlayerStatus `gorm:"type:tinyint;default:1;column:status" json:"status"`
	CreatedAt    time.Time          `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt    time.Time          `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}
