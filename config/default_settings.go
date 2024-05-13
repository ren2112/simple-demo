package config

const (
	// 资源路径
	SERVER_RESOURCES     = "http://" + LOCAL_IP_ADDRESS + ":8080/static/"
	DEFAULT_AVATAR_URL   = SERVER_RESOURCES + "initdata/avatar/"
	DEFAULT_BG_IMAGE_URL = SERVER_RESOURCES + "initdata/background/"

	// 默认用户信息
	DEFAULT_USER_AVATAR_URL   = DEFAULT_AVATAR_URL + "default.png"   // 默认头像地址
	DEFAULT_USER_BG_IMAGE_URL = DEFAULT_BG_IMAGE_URL + "default.png" // 默认背景图地址
	DEFAULT_USER_BIO          = "这个人很懒，什么也没有留下......"                // 默认简介内容
)
