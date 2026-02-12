package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "nba-remake/api/proto/v1"
	"nba-remake/internal/model"
)

// ListMatches 查赛程
func (s *NBAService) ListMatches(ctx context.Context, req *pb.ListMatchesRequest) (*pb.ListMatchesResponse, error) {
	matches, err := s.matchDao.ListByDate(req.Date)
	if err != nil {
		return nil, status.Error(codes.Internal, "查询失败: "+err.Error())
	}

	var resp []*pb.MatchResponse
	for _, m := range matches {
		resp = append(resp, convertMatchToProto(m))
	}
	return &pb.ListMatchesResponse{Matches: resp}, nil
}

// GetMatch 查详情
func (s *NBAService) GetMatch(ctx context.Context, req *pb.GetMatchRequest) (*pb.MatchResponse, error) {
	cacheKey := fmt.Sprintf("match:%d", req.Id)

	// 1. 先查 Redis
	val, err := s.redisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		// 缓存命中，直接反序列化返回
		var resp pb.MatchResponse
		err := json.Unmarshal([]byte(val), &resp)
		if err != nil {
			return nil, errors.New("缓存数据反序列化失败: " + err.Error())
		}
		return &resp, nil
	}

	// 2. 缓存未命中，查 MySQL
	match, err := s.matchDao.GetByID(req.Id)
	if err != nil {
		return nil, status.Error(codes.NotFound, "比赛未找到")
	}
	return convertMatchToProto(match), nil
}

// RecordMatchEvent 写入 Kafka
func (s *NBAService) RecordMatchEvent(ctx context.Context, req *pb.RecordMatchEventRequest) (*pb.RecordMatchEventResponse, error) {
	// 1. 校验 (现在 team_id 也是必填)
	if req.MatchId == 0 || req.PlayerId == 0 || req.TeamId == 0 {
		return nil, status.Error(codes.InvalidArgument, "参数缺失: match_id, player_id, team_id 必填")
	}

	// 2. 构造完整的 Payload
	// 消费者拿到这个 JSON 后，会解析并写入 match_events 表
	payload := map[string]interface{}{
		"match_id":       req.MatchId,
		"player_id":      req.PlayerId,
		"team_id":        req.TeamId, // 新增
		"type":           req.Type,
		"sub_type":       req.SubType, // 新增
		"value":          req.Value,
		"quarter":        req.Quarter,       // 新增
		"time_remaining": req.TimeRemaining, // 新增
		"event_time":     req.EventTime,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, status.Error(codes.Internal, "JSON序列化失败")
	}

	// 3. 发送 Kafka
	// Topic: "nba_match_events"
	// Key: MatchId (确保同一场比赛的消息顺序一致)
	err = s.kafkaProducer.Send(fmt.Sprintf("%d", req.MatchId), data)
	if err != nil {
		return nil, status.Error(codes.Unavailable, "MQ 服务不可用")
	}

	return &pb.RecordMatchEventResponse{Success: true, Message: "已推送"}, nil
}

// convertMatchToProto 辅助方法
func convertMatchToProto(m *model.Match) *pb.MatchResponse {
	return &pb.MatchResponse{
		Id:            int64(m.ID),
		Date:          m.Date.Format("2006-01-02"),
		HomeTeamId:    int32(m.HomeTeamID),
		VisitorTeamId: int32(m.VisitorTeamID),
		HomeScore:     int32(m.HomeScore),
		VisitorScore:  int32(m.VisitorScore),
		Status:        int32(m.Status),
		StartTime:     m.StartTime.Format("15:04"), // 显示几点开始
		HomeTeam:      convertTeamModelToProto(&m.HomeTeam),
		VisitorTeam:   convertTeamModelToProto(&m.VisitorTeam),
	}
}
