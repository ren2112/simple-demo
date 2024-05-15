package common

import (
	"context"
	"encoding/json"
	"github.com/RaymondCode/simple-demo/model"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"time"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:         viper.GetString("redis.addr"),
		Password:     viper.GetString("redis.password"),
		DB:           viper.GetInt("redis.DB"),
		PoolSize:     viper.GetInt("redis.poolSize"),
		MinIdleConns: viper.GetInt("redis.midIdleConn"),
	})
}

// 关于用户个人信息的缓存
// 检查token是否在黑名单
func CheckTokenInBlacklist(c context.Context, tokenString string) bool {
	_, err := RedisClient.Get(c, tokenString).Result()
	return err == nil
}

// 缓存用户
func CacheUser(c context.Context, userName string, user model.User) error {
	userData, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return RedisClient.Set(c, userName, userData, time.Hour*24).Err()
}

// 获取缓存用户信息
func GetCachedUser(c context.Context, userName string) (*model.User, error) {
	userData, err := RedisClient.Get(c, userName).Result()
	if err != nil {
		return nil, err
	}
	var user model.User
	err = json.Unmarshal([]byte(userData), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

//关于视频redis点赞与数据库的同步
