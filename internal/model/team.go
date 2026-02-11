package model

import "time"

// Team 对应数据库中的 teams 表
type Team struct {
	ID           uint32 `gorm:"primaryKey"`
	Name         string `gorm:"type:varchar(50);not null"`
	Abbreviation string `gorm:"type:char(3);not null;uniqueIndex"` // LAL, GSW
	City         string `gorm:"type:varchar(50);not null"`
	Conference   string `gorm:"type:enum('East','West');not null"`
	LogoURL      string `gorm:"column:logo_url;type:varchar(255)"`
	HomeArena    string `gorm:"column:home_arena;type:varchar(100)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
