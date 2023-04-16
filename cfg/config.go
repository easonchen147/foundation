package cfg

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const (
	Dev  = "dev"
	Qa   = "qa"
	Prod = "prod"
)

var AppConf *AppConfig
var gloablViper *viper.Viper

type AppConfig struct {
	File string

	Env           string `toml:"env"`
	HttpAddr      string `toml:"http_addr"`
	HttpPort      int    `toml:"http_port"`
	LogMode       string `toml:"log_mode"`
	LogFile       string `toml:"log_file"`
	LogLevel      string `toml:"log_level"`
	AccessLogFile string `toml:"access_log_file"`
	SqlLogFile    string `toml:"sql_log_file"`

	DbsConfig          map[string]*dbConfig `toml:"dbs"`
	MongoConfig        *mongoConfig         `toml:"mongo"`
	RedisConfig        *redisConfig         `toml:"redis"`
	RedisClusterConfig *redisClusterConfig  `toml:"redis_cluster"`
	KafkaConfig        *kafkaConfig         `toml:"kafka"`
	SignConfig         *signConfig          `toml:"sign"`
	TsConfig           *tsConfig            `toml:"ts"`

	Ext map[string]interface{} `toml:"ext"`
}

type dbConfig struct {
	Uri             string `toml:"uri"`
	MaxIdleConn     int    `toml:"max_idle_conn"`
	MaxOpenConn     int    `toml:"max_open_conn"`
	ConnectIdleTime int    `toml:"connect_idle_time"` //second default 300s
	ConnectLifeTime int    `toml:"connect_life_time"` //second default 600s
}

type redisConfig struct {
	Addr           string `toml:"addr"`
	Pass           string `toml:"pass"`
	Db             int    `toml:"db"`
	MinIdle        int    `toml:"min_idle"`
	PoolSize       int    `toml:"pool_size"`
	ConnectTimeout int    `toml:"connect_timeout"` //second default not set
	ReadTimeout    int    `toml:"read_timeout"`    //second default not set
	WriteTimeout   int    `toml:"write_timeout"`   //second default not set
}

type redisClusterConfig struct {
	Addrs          []string `toml:"addrs"`
	Pass           string   `toml:"pass"`
	MinIdle        int      `toml:"min_idle"`
	PoolSize       int      `toml:"pool_size"`
	ConnectTimeout int      `toml:"connect_timeout"` //second default not set
	ReadTimeout    int      `toml:"read_timeout"`    //second default not set
	WriteTimeout   int      `toml:"write_timeout"`   //second default not set
}

type kafkaConfig struct {
	Consumers map[string]*kafkaConsumerConfig `toml:"consumers"`
	Producers map[string]*kafkaProducerConfig `toml:"producers"`
}

type kafkaConsumerConfig struct {
	Broker    string `toml:"broker"`
	Topic     string `toml:"topic"`
	Group     string `toml:"group"`
	Partition int    `toml:"partition"`
}

type kafkaProducerConfig struct {
	Broker string `toml:"broker"`
	Topic  string `toml:"topic"`
}

type mongoConfig struct {
	Uri            string `toml:"uri"`
	Db             string `toml:"db"`
	ConnectTimeout uint64 `toml:"connect_timeout"`
	MaxOpenConn    uint64 `toml:"max_open_conn"`
	MaxPoolSize    uint64 `toml:"max_pool_size"`
	MinPoolSize    uint64 `toml:"min_pool_size"`
}

type signConfig struct {
	Secret string `toml:"secret"`
	Salt   string `toml:"salt"`
}

type tsConfig struct {
	Expire string `toml:"expire"`
}

func InitConfig(file string) *AppConfig {
	AppConf = &AppConfig{
		File:          file,
		Env:           Dev,
		HttpAddr:      "0.0.0.0",
		HttpPort:      8080,
		LogMode:       "console",
		LogFile:       "logs/app.log",
		LogLevel:      "debug",
		AccessLogFile: "logs/access.log",
	}
	return AppConf
}

// load 加载toml配置文件内容
func (cfg *AppConfig) load() error {
	if _, err := os.Stat(cfg.File); os.IsNotExist(err) {
		return fmt.Errorf("config file %s not existed", cfg.File)
	}

	// 全局唯一的viper
	gloablViper = viper.New()
	gloablViper.SetConfigFile(cfg.File)
	gloablViper.SetConfigType("toml")

	if err := gloablViper.ReadInConfig(); err != nil {
		return fmt.Errorf("load config file %s failed, error: %v", cfg.File, err)
	}
	if err := gloablViper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unmarshal %s to config object failed, error: %v", cfg.File, err)
	}
	return nil
}

func (cfg *AppConfig) IsDevEnv() bool {
	return cfg.Env == "dev"
}

func (cfg *AppConfig) LoadExtConfig(v interface{}) error {
	if gloablViper == nil {
		return errors.New("global viper is not initialize")
	}
	extV := gloablViper.Sub("ext")
	if extV == nil {
		return nil
	}
	return extV.Unmarshal(&v)
}

func init() {
	configFile := "app.toml"
	if envFilePath := os.Getenv("CONFIG_FILE"); envFilePath != "" {
		configFile = envFilePath
	}

	// 加载配置
	cfg := InitConfig(configFile)
	err := cfg.load()
	if err != nil {
		panic(fmt.Sprintf("load config failed, file: %s, error: %s", configFile, err))
	}
}
