package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	pb "nba-remake/api/proto/v1"
	"nba-remake/internal/config" // 引入 config
	"nba-remake/internal/dao"
	"nba-remake/internal/mq"
	"nba-remake/internal/service"
)

func main() {
	// 1. 加载配置 (第一件事)
	conf := config.LoadConfig()

	// 2. 数据库连接 (使用配置文件的 DSN)
	db, err := gorm.Open(mysql.Open(conf.MySQL.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("DB Error:", err)
	}

	// 3. 初始化 Kafka (传入配置部分的结构体)
	kafkaProducer, err := mq.NewProducer(conf.Kafka)
	if err != nil {
		log.Fatal("Kafka Error:", err)
	}
	defer kafkaProducer.Close()

	// 4. 初始化 DAO & Service
	playerDAO := dao.NewPlayerDao(db)
	teamDAO := dao.NewTeamDao(db)
	matchDAO := dao.NewMatchDao(db)

	nbaService := service.NewNBAService(playerDAO, teamDAO, matchDAO, kafkaProducer)

	// 5. 启动 gRPC (使用配置文件的 Port)
	lis, err := net.Listen("tcp", conf.Server.Port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterNBAServiceServer(s, nbaService)

	log.Printf("%s running on %s", conf.Server.Name, conf.Server.Port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
