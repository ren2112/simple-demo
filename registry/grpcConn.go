package registry

import (
	"github.com/RaymondCode/simple-demo/config"
	grpc_client_pool "github.com/RaymondCode/simple-demo/grpc-client-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"math/rand"
)

var AllPools = make(map[string][]*grpc_client_pool.ClientPool)

// 初始化所有连接池
func InitAllConnPool() {
	for _, s := range config.Server_list {
		initializeConnectionPool(s)
	}
}

func Balance(pools []*grpc_client_pool.ClientPool) (*grpc_client_pool.ClientPool, bool) {
	if len(pools) != 0 {
		return pools[rand.Intn(len(pools))], true
	} else {
		return nil, false
	}
}

// 初始化所有服务的连接池
func initializeConnectionPool(serverStr string) {
	services := ServiceDiscovery(serverStr)
	AllPools[serverStr] = nil
	for _, s := range services {
		addr := s.IP + ":" + s.Port
		ConnPool, err := grpc_client_pool.GetPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		AllPools[serverStr] = append(AllPools[serverStr], ConnPool)
	}
}

func GetPool(serviceName string) (*grpc_client_pool.ClientPool, bool) {
	res, ok := Balance(AllPools[serviceName])
	return res, ok
}
