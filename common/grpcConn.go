package common

import (
	"github.com/RaymondCode/simple-demo/config"
	grpc_client_pool "github.com/RaymondCode/simple-demo/grpc-client-pool"
	"github.com/RaymondCode/simple-demo/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var AllPools = make(map[string][]*grpc_client_pool.ClientPool)

// 初始化所有连接池
func InitAllConnPool() {
	for _, s := range config.Server_list {
		initializeConnectionPool(s)
	}
}

func balance(services []*registry.Service) *registry.Service {
	return services[0]
}

// 初始化Feed连接池
func initializeConnectionPool(serverStr string) {
	services := registry.ServiceDiscovery(serverStr)
	for _, s := range services {
		addr := s.IP + ":" + s.Port
		ConnPool, err := grpc_client_pool.GetPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatal(err)
		}
		AllPools[serverStr] = append(AllPools[serverStr], ConnPool)
	}
}
