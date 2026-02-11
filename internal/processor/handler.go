package processor

import (
	"encoding/json"
	"log"
	"time"

	"github.com/IBM/sarama"
	"gorm.io/gorm"

	"nba-remake/internal/model"
)

// EventDTO 用于接收 Kafka 消息的数据结构
// 保持与 Producer 发送的 JSON 字段一致
type EventDTO struct {
	MatchID       uint64 `json:"match_id"`
	PlayerID      uint32 `json:"player_id"`
	TeamID        uint32 `json:"team_id"`
	Type          int8   `json:"type"`
	SubType       string `json:"sub_type"`
	Value         int    `json:"value"`
	Quarter       int8   `json:"quarter"`
	TimeRemaining string `json:"time_remaining"`
	EventTime     string `json:"event_time"`
}

type StatsHandler struct {
	db *gorm.DB
}

func NewStatsHandler(db *gorm.DB) *StatsHandler {
	return &StatsHandler{db: db}
}

// Setup 在消费者组会话开始前执行
func (h *StatsHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// Cleanup 在消费者组会话结束后执行
func (h *StatsHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim 消费循环核心逻辑
func (h *StatsHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		// 业务处理
		if err := h.processEvent(msg.Value); err != nil {
			// 生产环境建议接入报警系统或写入死信队列(DLQ)
			log.Printf("[Consumer Error] partition=%d offset=%d err=%v", msg.Partition, msg.Offset, err)
		}

		// 标记消息为已处理
		session.MarkMessage(msg, "")
	}
	return nil
}

// processEvent 处理单条消息，使用数据库事务保证一致性
func (h *StatsHandler) processEvent(data []byte) error {
	var event EventDTO
	if err := json.Unmarshal(data, &event); err != nil {
		return err
	}

	// 开启事务
	return h.db.Transaction(func(tx *gorm.DB) error {
		// 1. 写入流水表
		newRecord := model.MatchEvent{
			MatchID:       event.MatchID,
			PlayerID:      event.PlayerID,
			TeamID:        event.TeamID,
			Type:          event.Type,
			SubType:       event.SubType,
			Value:         event.Value,
			Quarter:       event.Quarter,
			TimeRemaining: event.TimeRemaining,
			// EventTime 还是取当前写入时间较为准确，也可解析 event.EventTime
			EventTime: time.Now(),
		}

		if err := tx.Create(&newRecord).Error; err != nil {
			return err
		}

		// 2. 如果是得分事件(Type=1)，更新比赛主表比分
		// 使用 gorm.Expr 进行原子递增，防止并发覆盖
		if event.Type == 1 && event.Value > 0 {
			var match model.Match
			// 仅查询 TeamID 字段用于判断主客队，减少开销
			if err := tx.Select("id", "home_team_id", "visitor_team_id").First(&match, event.MatchID).Error; err != nil {
				return err
			}

			if uint32(match.HomeTeamID) == event.TeamID {
				// 更新主队得分
				if err := tx.Model(&model.Match{}).Where("id = ?", event.MatchID).
					UpdateColumn("home_score", gorm.Expr("home_score + ?", event.Value)).Error; err != nil {
					return err
				}
			} else if uint32(match.VisitorTeamID) == event.TeamID {
				// 更新客队得分
				if err := tx.Model(&model.Match{}).Where("id = ?", event.MatchID).
					UpdateColumn("visitor_score", gorm.Expr("visitor_score + ?", event.Value)).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
