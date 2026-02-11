package dao

import (
	"gorm.io/gorm"
	"nba-remake/internal/model"
)

type MatchDao struct {
	db *gorm.DB
}

func NewMatchDao(db *gorm.DB) *MatchDao {
	return &MatchDao{db: db}
}

// GetByID 查单场 (Preload 球队)
func (d *MatchDao) GetByID(id int64) (*model.Match, error) {
	var match model.Match
	err := d.db.Preload("HomeTeam").Preload("VisitorTeam").
		Where("id = ?", id).First(&match).Error
	return &match, err
}

// ListByDate 查列表 (按时间排序)
func (d *MatchDao) ListByDate(date string) ([]*model.Match, error) {
	var matches []*model.Match
	err := d.db.Preload("HomeTeam").Preload("VisitorTeam").
		Where("date = ?", date).
		Order("start_time asc").
		Find(&matches).Error
	return matches, err
}
