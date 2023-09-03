package redis

import (
	"context"
	"fmt"

	"github.com/lyujsally/MinTikTok-lyujsally/settings"

	"github.com/go-redis/redis"
)

// 声明全局的rdb变量
var Ctx = context.Background()
var RdbFollowee *redis.Client
var RdbFollower *redis.Client
var RdbIsFollow *redis.Client
var RDB *redis.Client

// 初始化连接
func InitRedisCli(cfg *settings.RedisConfig) (err error) {

	//用户完整关注列表信息存入DB1
	RdbFollowee = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s%d", cfg.Host, ":", cfg.Port),
		Password: cfg.Password,
		DB:       1,
		PoolSize: cfg.PoolSize,
	})

	//用户完整粉丝列表信息存入DB2
	RdbFollower = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s%d", cfg.Host, ":", cfg.Port),
		Password: cfg.Password,
		DB:       2,
		PoolSize: cfg.PoolSize,
	})

	//一些用户关注热数据存入DB3
	RdbIsFollow = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s%d", cfg.Host, ":", cfg.Port),
		Password: cfg.Password,
		DB:       3,
		PoolSize: cfg.PoolSize,
	})

	RDB = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s%s%d", cfg.Host, ":", cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
		PoolSize: cfg.PoolSize,
	})

	_, err = RDB.Ping().Result()
	return err
}

func RedisClose() {
	_ = RDB.Close()
}
