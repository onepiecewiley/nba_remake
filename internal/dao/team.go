package dao

import (
	"gorm.io/gorm"
	"nba-remake/internal/model"
)

type TeamDao struct {
	db *gorm.DB
}

// NewTeamDao 构造函数
func NewTeamDao(db *gorm.DB) *TeamDao {
	return &TeamDao{db: db}
}

// GetByID 根据ID获取球队
func (d *TeamDao) GetByID(id int32) (*model.Team, error) {
	var team model.Team
	// GORM 查不到数据会返回 ErrRecordNotFound
	err := d.db.Where("id = ?", id).First(&team).Error
	return &team, err
}

// GetAll 获取所有球队
func (d *TeamDao) GetAll() ([]*model.Team, error) {
	var teams []*model.Team
	// 按照名字排序
	err := d.db.Order("name asc").Find(&teams).Error
	return teams, err
}
