package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	MySQL  MySQLConfig  `mapstructure:"mysql"`
	Kafka  KafkaConfig  `mapstructure:"kafka"`
}

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Port string `mapstructure:"port"`
}

type MySQLConfig struct {
	DSN     string `mapstructure:"dsn"`
	MaxIdle int    `mapstructure:"max_idle"`
	MaxOpen int    `mapstructure:"max_open"`
}

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
	Topic   string   `mapstructure:"topic"`
	GroupID string   `mapstructure:"group_id"`
}

// LoadConfig 读取配置文件
func LoadConfig() *Config {
	viper.SetConfigName("config")    // 配置文件名
	viper.SetConfigType("yaml")      // 文件格式
	viper.AddConfigPath(".")         // 搜索路径
	viper.AddConfigPath("./configs") // 也可以搜 configs 目录

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("读取配置文件失败: %v", err)
	}

	var config Config
	// 将读取到的配置映射到结构体中
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("配置解析失败: %v", err)
	}

	log.Println("配置加载成功")
	return &config
}
