package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	pb "nba-remake/api/proto/v1"
	"nba-remake/internal/config"
	"nba-remake/internal/dao"
	"nba-remake/internal/mq"
	"nba-remake/internal/processor"
	"nba-remake/internal/service"
)

func main() {
	// 1. 加载配置
	conf := config.LoadConfig()

	// 2. 初始化共享资源 (DB)
	db, err := gorm.Open(mysql.Open(conf.MySQL.DSN), &gorm.Config{})
	if err != nil {
		log.Fatal("DB连接失败:", err)
	}

	// 初始化 Kafka Producer
	kafkaProducer, err := mq.NewProducer(conf.Kafka)
	if err != nil {
		log.Fatal("Kafka Producer 失败:", err)
	}
	defer kafkaProducer.Close()

	// 初始化 DAO & Service
	playerDAO := dao.NewPlayerDao(db)
	teamDAO := dao.NewTeamDao(db)
	matchDAO := dao.NewMatchDao(db)
	nbaService := service.NewNBAService(playerDAO, teamDAO, matchDAO, kafkaProducer)

	// 初始化 gRPC Server
	server := grpc.NewServer()
	pb.RegisterNBAServiceServer(server, nbaService)

	lis, err := net.Listen("tcp", conf.Server.Port)
	if err != nil {
		log.Fatal("端口监听失败:", err)
	}

	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	consumerGroup, err := sarama.NewConsumerGroup(conf.Kafka.Brokers, conf.Kafka.GroupID, saramaConfig)
	if err != nil {
		log.Fatal("Kafka Consumer Group 失败:", err)
	}
	defer consumerGroup.Close()

	statsHandler := processor.NewStatsHandler(db)
	ctx, cancel := context.WithCancel(context.Background())

	// 1. 启动 gRPC 服务
	go func() {
		log.Printf("gRPC 服务启动于 %s", conf.Server.Port)
		if err := server.Serve(lis); err != nil {
			log.Printf("gRPC 停止: %v", err)
		}
	}()

	// 2. 启动 消费者 服务
	go func() {
		log.Println("👷 消费者 Worker 已启动")
		for {
			if err := consumerGroup.Consume(ctx, []string{conf.Kafka.Topic}, statsHandler); err != nil {
				log.Printf("消费错误: %v", err)
				time.Sleep(time.Second * 2) // 失败休眠一下
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	// 阻塞主线程，直到收到 Ctrl+C
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("正在停止服务...")

	// 停止 gRPC
	server.GracefulStop()

	// 停止 Consumer
	cancel()

	log.Println("服务已全部停止")
}
