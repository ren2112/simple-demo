package etcd

import "time"

const DailTime = 5 * time.Second

func GetEtcdEndpoints() []string {
	return []string{"localhost:2380"} //注册中心地址
}
