package redis_utils

import (
	"fmt"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

type RedisUtilsConfig struct {
	RedisConfig redis.RedisConf
	BizName     string
}

type RedisUtils struct {
	rds     *redis.Redis
	bizName string
}

func (r *RedisUtils) fmtStrWithPrefix(s string) string {
	return fmt.Sprintf("%s/%s", r.bizName, s)
}

func (r *RedisUtils) Set(k string, v string) error {
	return r.rds.Set(r.fmtStrWithPrefix(k), v)
}

func (r *RedisUtils) SetEx(k string, v string, seconds int) error {
	return r.rds.Setex(r.fmtStrWithPrefix(k), v, seconds)
}

func (r *RedisUtils) Get(k string) (string, error) {
	return r.rds.Get(r.fmtStrWithPrefix(k))
}

func (r *RedisUtils) Delete(ks ...string) (int, error) {
	var keys []string = make([]string, len(ks))
	for _, k := range ks {
		keys = append(keys, r.fmtStrWithPrefix(k))
	}
	return r.rds.Del(keys...)
}

func (r *RedisUtils) Exists(k string) (bool, error) {
	return r.rds.Exists(r.fmtStrWithPrefix(k))
}

func NewRedisUtils(rds *redis.Redis, bizName string) *RedisUtils {
	return &RedisUtils{
		rds:     rds,
		bizName: bizName,
	}
}

func InitRedisUtils(rc RedisUtilsConfig) *RedisUtils {
	rds := redis.MustNewRedis(rc.RedisConfig)
	return NewRedisUtils(rds, rc.BizName)
}
