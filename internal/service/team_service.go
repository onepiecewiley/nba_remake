package service

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	pb "nba-remake/api/proto/v1"
	"nba-remake/internal/model"
)

// GetTeam 实现 gRPC GetTeam 接口
func (s *NBAService) GetTeam(ctx context.Context, req *pb.GetTeamRequest) (*pb.TeamResponse, error) {
	// 1. 调用 DAO
	team, err := s.teamDao.GetByID(req.Id)

	// 2. 错误处理
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, status.Error(codes.NotFound, "球队不存在")
		}
		return nil, status.Error(codes.Internal, "数据库查询失败")
	}

	// 3. 转换 Model -> Proto Response
	return convertTeamModelToProto(team), nil
}

// ListTeams 实现 gRPC ListTeams 接口
func (s *NBAService) ListTeams(ctx context.Context, req *pb.ListTeamsRequest) (*pb.ListTeamsResponse, error) {
	// 1. 调用 DAO
	teams, err := s.teamDao.GetAll()
	if err != nil {
		return nil, status.Error(codes.Internal, "获取球队列表失败")
	}

	// 2. 批量转换
	var respTeams []*pb.TeamResponse
	for _, t := range teams {
		respTeams = append(respTeams, convertTeamModelToProto(t))
	}

	return &pb.ListTeamsResponse{
		Teams: respTeams,
	}, nil
}

// 辅助函数：Model 转 Proto
func convertTeamModelToProto(t *model.Team) *pb.TeamResponse {
	return &pb.TeamResponse{
		Id:           int32(t.ID),
		Name:         t.Name,
		City:         t.City,
		Abbreviation: t.Abbreviation,
		Conference:   t.Conference,
		LogoUrl:      t.LogoURL,
	}
}
