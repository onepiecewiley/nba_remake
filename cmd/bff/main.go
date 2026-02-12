package main

import (
	"context"
	"log"
	myErrors "nba-remake/errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "nba-remake/api/proto/v1"
)

func main() {
	// 1. 连接后端的 gRPC 服务
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("无法连接 gRPC 服务: %v", err)
	}
	defer conn.Close()

	client := pb.NewNBAServiceClient(conn)

	// 2. 初始化 Gin
	r := gin.Default()

	// 3. 配置 CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 球员相关路由
	r.GET("/api/players", func(c *gin.Context) {
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
		name := c.Query("name")
		teamID, _ := strconv.Atoi(c.Query("team_id"))
		position := c.Query("position")
		status, _ := strconv.Atoi(c.Query("status"))

		resp, err := client.ListPlayers(context.Background(), &pb.ListPlayersRequest{
			Page:     int32(page),
			PageSize: int32(pageSize),
			Name:     name,
			TeamId:   int32(teamID),
			Position: pb.Position(pb.Position_value[position]),
			Status:   pb.PlayerStatus(status),
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/api/players/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		resp, err := client.GetPlayer(context.Background(), &pb.GetPlayerRequest{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "球员未找到"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.POST("/api/players", func(c *gin.Context) {
		var req struct {
			Name         string  `json:"name"`
			TeamId       int32   `json:"team_id"`
			JerseyNumber int32   `json:"jersey_number"`
			Position     string  `json:"position"`
			Height       float64 `json:"height"`
			Weight       float64 `json:"weight"`
			Birthday     string  `json:"birthday"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}

		resp, err := client.CreatePlayer(context.Background(), &pb.CreatePlayerRequest{
			Name:         req.Name,
			TeamId:       req.TeamId,
			JerseyNumber: req.JerseyNumber,
			Position:     pb.Position(pb.Position_value[req.Position]),
			Height:       req.Height,
			Weight:       req.Weight,
			Birthday:     req.Birthday,
		})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusCreated, resp)
	})

	// 球队相关路由
	r.GET("/api/teams", func(c *gin.Context) {
		resp, err := client.ListTeams(context.Background(), &pb.ListTeamsRequest{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	r.GET("/api/teams/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		resp, err := client.GetTeam(context.Background(), &pb.GetTeamRequest{Id: int32(id)})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": myErrors.NewError(myErrors.CodeTeamNotFound, "球队未找到", "")})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	// 比赛相关路由（你原有的代码）
	r.GET("/api/matches", func(c *gin.Context) {
		date := c.Query("date")
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		resp, err := client.ListMatches(context.Background(), &pb.ListMatchesRequest{Date: date})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp.Matches)
	})

	r.GET("/api/matches/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, _ := strconv.ParseInt(idStr, 10, 64)

		resp, err := client.GetMatch(context.Background(), &pb.GetMatchRequest{Id: id})
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "比赛未找到"})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	// 事件路由
	r.POST("/api/matches/events", func(c *gin.Context) {
		var req struct {
			MatchID       int64  `json:"match_id"`
			PlayerID      int32  `json:"player_id"`
			TeamID        int32  `json:"team_id"`
			Type          int32  `json:"type"`
			SubType       string `json:"sub_type"`
			Value         int32  `json:"value"`
			Quarter       int32  `json:"quarter"`
			TimeRemaining string `json:"time_remaining"`
		}

		if err := c.BindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
			return
		}

		_, err := client.RecordMatchEvent(context.Background(), &pb.RecordMatchEventRequest{
			MatchId:       req.MatchID,
			PlayerId:      req.PlayerID,
			TeamId:        req.TeamID,
			Type:          req.Type,
			SubType:       req.SubType,
			Value:         req.Value,
			Quarter:       req.Quarter,
			TimeRemaining: req.TimeRemaining,
			EventTime:     time.Now().Format(time.RFC3339),
		})

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "发送失败"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "queued"})
	})

	// 启动 BFF
	log.Println("BFF Server 运行在 :8080")
	r.Run(":8080")
}
