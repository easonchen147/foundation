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

	Env           string `mapstructure:"env"`
	HttpAddr      string `mapstructure:"http_addr"`
	HttpPort      int    `mapstructure:"http_port"`
	LogMode       string `mapstructure:"log_mode"`
	LogFile       string `mapstructure:"log_file"`
	LogLevel      string `mapstructure:"log_level"`
	AccessLogFile string `mapstructure:"access_log_file"`
	SqlLogFile    string `mapstructure:"sql_log_file"`

	DbsConfig          map[string]*dbConfig `mapstructure:"dbs"`
	MongoConfig        *mongoConfig         `mapstructure:"mongo"`
	RedisConfig        *redisConfig         `mapstructure:"redis"`
	RedisClusterConfig *redisClusterConfig  `mapstructure:"redis_cluster"`
	KafkaConfig        *kafkaConfig         `mapstructure:"kafka"`
	SignConfig         *signConfig          `mapstructure:"sign"`
	TsConfig           *tsConfig            `mapstructure:"ts"`

	Ext map[string]interface{} `mapstructure:"ext"`
}

type dbConfig struct {
	Uri             string `mapstructure:"uri"`
	MaxIdleConn     int    `mapstructure:"max_idle_conn"`
	MaxOpenConn     int    `mapstructure:"max_open_conn"`
	ConnectIdleTime int    `mapstructure:"connect_idle_time"` //second default 300s
	ConnectLifeTime int    `mapstructure:"connect_life_time"` //second default 600s
}

type redisConfig struct {
	Addr           string `mapstructure:"addr"`
	Pass           string `mapstructure:"pass"`
	Db             int    `mapstructure:"db"`
	MinIdle        int    `mapstructure:"min_idle"`
	PoolSize       int    `mapstructure:"pool_size"`
	ConnectTimeout int    `mapstructure:"connect_timeout"` //second default not set
	ReadTimeout    int    `mapstructure:"read_timeout"`    //second default not set
	WriteTimeout   int    `mapstructure:"write_timeout"`   //second default not set
}

type redisClusterConfig struct {
	Addrs          []string `mapstructure:"addrs"`
	Pass           string   `mapstructure:"pass"`
	MinIdle        int      `mapstructure:"min_idle"`
	PoolSize       int      `mapstructure:"pool_size"`
	ConnectTimeout int      `mapstructure:"connect_timeout"` //second default not set
	ReadTimeout    int      `mapstructure:"read_timeout"`    //second default not set
	WriteTimeout   int      `mapstructure:"write_timeout"`   //second default not set
}

type kafkaConfig struct {
	Consumers map[string]*kafkaConsumerConfig `mapstructure:"consumers"`
	Producers map[string]*kafkaProducerConfig `mapstructure:"producers"`
}

type kafkaConsumerConfig struct {
	Broker    string `mapstructure:"broker"`
	Topic     string `mapstructure:"topic"`
	Group     string `mapstructure:"group"`
	Partition int    `mapstructure:"partition"`
}

type kafkaProducerConfig struct {
	Broker string `mapstructure:"broker"`
	Topic  string `mapstructure:"topic"`
}

type mongoConfig struct {
	Uri            string `mapstructure:"uri"`
	Db             string `mapstructure:"db"`
	ConnectTimeout uint64 `mapstructure:"connect_timeout"`
	MaxOpenConn    uint64 `mapstructure:"max_open_conn"`
	MaxPoolSize    uint64 `mapstructure:"max_pool_size"`
	MinPoolSize    uint64 `mapstructure:"min_pool_size"`
}

type signConfig struct {
	Secret string `mapstructure:"secret"`
	Salt   string `mapstructure:"salt"`
}

type tsConfig struct {
	Expire string `mapstructure:"expire"`
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
