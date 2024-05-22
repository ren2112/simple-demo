package config

const (
	LOCAL_IP_ADDRESS        = "192.168.249.67"      // 填入本机 IP 地址
	VIDEO_STREAM_BATCH_SIZE = 30                    // 每次获取视频流的数量限制
	DATETIME_FORMAT         = "2006-01-02 15:04:05" // 固定的时间格式
	AUTH_KEY                = "a_secret_key"        // JWT 密钥
	TOKENBUCKET_CAPACITY    = 100                   //令牌桶容量
	TOKENBUCKET_RATE        = 10000                 //令牌桶令牌生成频率

)
