package model

import "time"

// Match 比赛主表
type Match struct {
	ID            uint64    `gorm:"primaryKey"`
	Date          time.Time `gorm:"type:date;index"`
	Season        string    `gorm:"type:varchar(10);default:'2023-24'"` // SQL里有这个
	HomeTeamID    uint      `gorm:"not null"`
	VisitorTeamID uint      `gorm:"not null"`
	HomeScore     int       `gorm:"default:0"`
	VisitorScore  int       `gorm:"default:0"`
	Status        int       `gorm:"default:0"`
	StartTime     time.Time

	HomeTeam    Team `gorm:"foreignKey:HomeTeamID"`
	VisitorTeam Team `gorm:"foreignKey:VisitorTeamID"`
}

// MatchEvent 比赛事件流水表
// 对应数据库: match_events
type MatchEvent struct {
	ID            uint64    `gorm:"primaryKey;autoIncrement"`
	MatchID       uint64    `gorm:"column:match_id;not null;index"`
	PlayerID      uint32    `gorm:"column:player_id;not null;index"`
	TeamID        uint32    `gorm:"column:team_id;not null"`                // 发生时属于哪个队(冗余)
	Type          int8      `gorm:"column:type;not null"`                   // 事件类型 (TINYINT)
	SubType       string    `gorm:"column:sub_type;type:varchar(20)"`       // 子类型: 3pt, dunk, layup
	Value         int       `gorm:"column:value;not null;default:0"`        // 分值
	Quarter       int8      `gorm:"column:quarter;default:1"`               // 第几节
	TimeRemaining string    `gorm:"column:time_remaining;type:varchar(10)"` // 剩余时间 e.g. "10:23"
	EventTime     time.Time `gorm:"column:event_time;autoCreateTime"`       // 物理写入时间
}
