package service

import (
	"context"
	pb "nba-remake/api/proto/v1"
	"nba-remake/internal/dao"
	"nba-remake/internal/model"
	"time"
)

type NBAService struct {
	pb.UnimplementedNBAServiceServer
	playerDao *dao.PlayerDao
}

func NewNBAService(playerDao *dao.PlayerDao) *NBAService {
	return &NBAService{
		playerDao: playerDao,
	}
}

func (s *NBAService) CreatePlayer(ctx context.Context, req *pb.CreatePlayerRequest) (*pb.PlayerResponse, error) {
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		return nil, err
	}

	p := &model.Player{
		Name:         req.Name,
		JerseyNumber: uint8(req.JerseyNumber),
		Position:     req.Position,
		Height:       req.Height,
		Weight:       req.Weight,
		Birthday:     &birthday,
		Status:       req.Status,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = s.playerDao.CreatePlayer(p)
	if err != nil {
		return nil, err
	}
	resp := &pb.PlayerResponse{
		Id:           int32(p.ID),
		Name:         p.Name,
		JerseyNumber: req.JerseyNumber,
		Position:     req.Position,
		Height:       req.Height,
		Weight:       req.Weight,
		Birthday:     req.Birthday,
		Status:       req.Status,
		CreatedAt:    p.CreatedAt.Format("time.RFC3339"),
		UpdatedAt:    p.UpdatedAt.Format(time.RFC3339),
	}
	return resp, nil
}

func (s *NBAService) GetPlayer(ctx context.Context, req *pb.GetPlayerRequest) (*pb.PlayerResponse, error) {
	player, err := s.playerDao.GetPlayerByID(uint32(req.Id))
	if err != nil {
		return nil, err
	}
	resp := &pb.PlayerResponse{
		Id:           int32(player.ID),
		Name:         player.Name,
		JerseyNumber: int32(player.JerseyNumber),
		Position:     player.Position,
		Height:       player.Height,
		Weight:       player.Weight,
		Birthday:     player.Birthday.Format("2006-01-02"),
		Status:       player.Status,
		CreatedAt:    player.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    player.UpdatedAt.Format(time.RFC3339),
	}
	return resp, nil
}

func (s *NBAService) UpdatePlayer(ctx context.Context, req *pb.UpdatePlayerRequest) (*pb.PlayerResponse, error) {
	player, err := s.playerDao.GetPlayerByID(uint32(req.Id))
	if err != nil {
		return nil, err
	}
	player.Name = req.Name
	player.JerseyNumber = uint8(req.JerseyNumber)
	player.Position = req.Position
	player.Height = req.Height
	player.Weight = req.Weight
	player.Status = req.Status
	player.UpdatedAt = time.Now()
	err = s.playerDao.UpdatePlayer(player)
	if err != nil {
		return nil, err
	}
	resp := &pb.PlayerResponse{
		Id:           int32(player.ID),
		Name:         player.Name,
		JerseyNumber: int32(player.JerseyNumber),
		Position:     player.Position,
		Height:       player.Height,
		Weight:       player.Weight,
		Birthday:     player.Birthday.Format("2006-01-02"),
		Status:       player.Status,
		CreatedAt:    player.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    player.UpdatedAt.Format(time.RFC3339),
	}
	return resp, nil
}

func (s *NBAService) DeletePlayer(ctx context.Context, req *pb.DeletePlayerRequest) (*pb.DeletePlayerResponse, error) {
	err := s.playerDao.DeletePlayer(uint32(req.Id))
	if err != nil {
		return nil, err
	}
	return &pb.DeletePlayerResponse{Success: true}, nil
}

func (s *NBAService) ListPlayers(ctx context.Context, req *pb.ListPlayersRequest) (*pb.ListPlayersResponse, error) {
	page := int(req.Page)
	pageSize := int(req.PageSize)
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	players, total, err := s.playerDao.ListPlayersByFilter(req.Name, uint32(req.TeamId), req.Position, uint8(req.Status), page, pageSize)
	if err != nil {
		return nil, err
	}

	pbPlayers := make([]*pb.PlayerResponse, len(players))
	for i, player := range players {
		pbPlayers[i] = &pb.PlayerResponse{
			Id:           int32(player.ID),
			Name:         player.Name,
			JerseyNumber: int32(player.JerseyNumber),
			Position:     player.Position,
			Height:       player.Height,
			Weight:       player.Weight,
			Birthday:     player.Birthday.Format("2006-01-02"),
			Status:       player.Status,
			CreatedAt:    player.CreatedAt.Format(time.RFC3339),
			UpdatedAt:    player.UpdatedAt.Format(time.RFC3339),
		}
	}

	return &pb.ListPlayersResponse{
		Players:  pbPlayers,
		Total:    int32(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
