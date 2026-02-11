package dao

import (
	"gorm.io/gorm"
	pb "nba-remake/api/proto/v1"
	"nba-remake/internal/model"
)

type PlayerDao struct {
	db *gorm.DB
}

func NewPlayerDao(db *gorm.DB) *PlayerDao {
	return &PlayerDao{db: db}
}

// 创建球员
func (d *PlayerDao) CreatePlayer(player *model.Player) error {
	return d.db.Create(player).Error
}

// 根据id查询
func (d *PlayerDao) GetPlayerByID(id uint32) (*model.Player, error) {
	var player model.Player
	err := d.db.First(&player, id).Error
	if err != nil {
		return nil, err
	}
	return &player, nil
}

// 更新
func (d *PlayerDao) UpdatePlayer(player *model.Player) error {
	return d.db.Save(player).Error
}

// 删除
func (d *PlayerDao) DeletePlayer(id uint32) error {
	return d.db.Delete(&model.Player{}, id).Error
}

func (d *PlayerDao) ListPlayersByFilter(name string, teamID uint32, position pb.Position, status uint8, page, pageSize int) ([]*model.Player, int64, error) {
	var players []*model.Player
	var total int64

	query := d.db.Model(&model.Player{})

	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if teamID > 0 {
		query = query.Where("team_id = ?", teamID)
	}
	if position != pb.Position_POSITION_UNKNOWN {
		query = query.Where("position = ?", position)
	}
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	query.Count(&total)

	if page > 0 && pageSize > 0 {
		query = query.Offset((page - 1) * pageSize).Limit(pageSize)
	}

	if err := query.Find(&players).Error; err != nil {
		return nil, 0, err
	}

	return players, total, nil
}
