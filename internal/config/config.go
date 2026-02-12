package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config 全局配置结构体
type Config struct {
	Server        ServerConfig        `mapstructure:"server"`
	MySQL         MySQLConfig         `mapstructure:"mysql"`
	Kafka         KafkaConfig         `mapstructure:"kafka"`
	Redis         RedisConfig         `mapstructure:"redis"`         // 新增Redis配置
	Elasticsearch ElasticsearchConfig `mapstructure:"elasticsearch"` // 新增ES配置
	MongoDB       MongoDBConfig       `mapstructure:"mongodb"`       // 新增MongoDB配置
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

type RedisConfig struct {
	Addr           string `mapstructure:"addr"`            // Redis地址（ip:port）
	Password       string `mapstructure:"password"`        // Redis密码（空则无认证）
	DB             int    `mapstructure:"db"`              // 使用的数据库编号
	MaxIdle        int    `mapstructure:"max_idle"`        // 最大空闲连接数
	MaxActive      int    `mapstructure:"max_active"`      // 最大活跃连接数
	IdleTimeout    string `mapstructure:"idle_timeout"`    // 空闲连接超时时间（如300s）
	ReadTimeout    string `mapstructure:"read_timeout"`    // 读超时（如3s）
	WriteTimeout   string `mapstructure:"write_timeout"`   // 写超时（如3s）
	ConnectTimeout string `mapstructure:"connect_timeout"` // 连接超时（如5s）
}

type ElasticsearchConfig struct {
	Addresses   []string `mapstructure:"addresses"`    // ES地址列表（支持集群）
	Username    string   `mapstructure:"username"`     // ES用户名（7.x+默认无）
	Password    string   `mapstructure:"password"`     // ES密码（7.x+默认无）
	Timeout     string   `mapstructure:"timeout"`      // 连接/读写超时时间（如10s）
	MaxRetries  int      `mapstructure:"max_retries"`  // 请求失败最大重试次数
	IndexPrefix string   `mapstructure:"index_prefix"` // 索引前缀（如nba_）
	Sniff       bool     `mapstructure:"sniff"`        // 是否自动发现集群节点
}

type MongoDBConfig struct {
	URI            string `mapstructure:"uri"`             // MongoDB连接字符串
	Database       string `mapstructure:"database"`        // 数据库名
	ConnectTimeout string `mapstructure:"connect_timeout"` // 连接超时（如10s）
	SocketTimeout  string `mapstructure:"socket_timeout"`  // 套接字超时（如30s）
	MaxPoolSize    int    `mapstructure:"max_pool_size"`   // 最大连接池大小
	MinPoolSize    int    `mapstructure:"min_pool_size"`   // 最小连接池大小
	MaxIdleTime    string `mapstructure:"max_idle_time"`   // 连接最大空闲时间（如60s）
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
