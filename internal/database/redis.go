package database

import (
	"fmt"
	"github.com/kilicmu/user-service/redis_utils"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"os"
)

var redisUtils *redis_utils.RedisUtils

func InitRedisUtils() error {
	bizName := os.Getenv("BIZ_NAME")
	if bizName == "" {
		panic("init redis util fail, env [BIZ_NAME] is required.")
	}
	redisUtils = redis_utils.InitRedisUtils(redis_utils.RedisUtilsConfig{
		RedisConfig: redis.RedisConf{
			Host: fmt.Sprintf("%s:%s", os.Getenv("RDS_ADDRESS"), os.Getenv("RDS_PORT")),
			Type: "node",
			Pass: "",
			Tls:  false,
		},
		BizName: bizName,
	})
	return nil
}

func GetRedisUtils() *redis_utils.RedisUtils {
	return redisUtils
}
