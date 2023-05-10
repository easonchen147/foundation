package cache

import (
	"errors"
	"fmt"
	"time"

	"github.com/easonchen147/foundation/cfg"
	"github.com/redis/go-redis/v9"
)

var (
	client        *redis.Client
	clusterClient *redis.ClusterClient
)

func init() {
	err := InitRedis(cfg.AppConf)
	if err != nil {
		panic(fmt.Sprintf("init redis failed: %s", err))
	}

	err = InitRedisCluster(cfg.AppConf)
	if err != nil {
		panic(fmt.Sprintf("init redis cluster failed: %s", err))
	}
}

// InitRedis 初始化redis
func InitRedis(cfg *cfg.AppConfig) error {
	if cfg.RedisConfig == nil {
		return nil
	}
	client = redis.NewClient(&redis.Options{
		Addr:         cfg.RedisConfig.Addr,
		Username:     cfg.RedisConfig.User,
		Password:     cfg.RedisConfig.Pass,
		DB:           cfg.RedisConfig.Db,
		MinIdleConns: cfg.RedisConfig.MinIdle,
		PoolSize:     cfg.RedisConfig.PoolSize,
		DialTimeout:  time.Second * time.Duration(cfg.RedisConfig.ConnectTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.RedisConfig.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.RedisConfig.WriteTimeout),
	})
	return nil
}

func Redis() *redis.Client {
	if client == nil {
		panic(errors.New("cache is not ready"))
	}
	return client
}

// InitRedisCluster 初始化redis cluster
func InitRedisCluster(cfg *cfg.AppConfig) error {
	if cfg.RedisClusterConfig == nil {
		return nil
	}
	clusterClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        cfg.RedisClusterConfig.Addrs,
		Password:     cfg.RedisClusterConfig.Pass,
		MinIdleConns: cfg.RedisClusterConfig.MinIdle,
		PoolSize:     cfg.RedisClusterConfig.PoolSize,
		DialTimeout:  time.Second * time.Duration(cfg.RedisConfig.ConnectTimeout),
		ReadTimeout:  time.Second * time.Duration(cfg.RedisConfig.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.RedisConfig.WriteTimeout),
	})
	return nil
}

func RedisCluster() *redis.ClusterClient {
	if clusterClient == nil {
		panic(errors.New("foundation cluster is not ready"))
	}
	return clusterClient
}

func Close() {
	if client != nil {
		_ = client.Close()
	}
	if clusterClient != nil {
		_ = clusterClient.Close()
	}
}
