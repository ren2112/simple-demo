package config

import "time"

var Server_list = []string{"feed", "user", "publish", "relation", "message", "comment", "friend", "favorite"}

const (
	TIMEOUT                 = 5 * time.Second
	LOCAL_IP_ADDRESS        = "192.168.241.67"      // 填入本机 IP 地址
	VIDEO_STREAM_BATCH_SIZE = 30                    // 每次获取视频流的数量限制
	DATETIME_FORMAT         = "2006-01-02 15:04:05" // 固定的时间格式
	AUTH_KEY                = "a_secret_key"        // JWT 密钥
	TOKENBUCKET_CAPACITY    = 100                   //令牌桶容量
	TOKENBUCKET_RATE        = 10000                 //令牌桶生成令牌频率
	FEED_SERVER_ADDR        = "127.0.0.1:9091"
	USER_SERVER_ADDR        = "127.0.0.1:9093"
	PUBLISH_SERVER_ADDR     = "127.0.0.1:9093"
	RELATION_SERVER_ADDR    = "127.0.0.1:9094"
	FAVORITE_SERVER_ADDR    = "127.0.0.1:9095"
	COMMENT_SERVER_ADDR     = "127.0.0.1:9096"
	MESSAGE_SERVER_ADDR     = "127.0.0.1:9097"
	FRIEND_SERVER_ADDR      = "127.0.0.1:9098"
)
