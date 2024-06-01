package registry

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/etcd"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"strings"
	"sync"
)

type Service struct {
	Name     string
	IP       string
	Port     string
	Protocol string
}

// 服务注册
func ServiceRegister(s *Service) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcd.GetEtcdEndpoints(),
		DialTimeout: etcd.DailTime,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	var leaseId clientv3.LeaseID
	ctx := context.Background()

	servicekey := fmt.Sprintf("%s/%s:%s", s.Name, s.IP, s.Port)
	leaseRes, err := cli.Grant(ctx, 10)
	if err != nil {
		log.Fatal(err)
	}
	leaseId = leaseRes.ID
	kv := clientv3.NewKV(cli)
	_, err = kv.Put(ctx, servicekey, fmt.Sprintf("%s:%s:%s", s.IP, s.Port, s.Protocol), clientv3.WithLease(leaseId))
	if err != nil {
		log.Fatal(err)
	}

	leaseKeepalive, err := cli.KeepAlive(ctx, leaseId)
	if err != nil {
		log.Fatal(err)
	}
	//需要及时获取反馈消息，否则当做挂死
	for range leaseKeepalive {
	}
}

type Services struct {
	services map[string][]*Service
	sync.RWMutex
}

var douyinServices = &Services{
	services: map[string][]*Service{},
}

// 服务发现
func ServiceDiscovery(serviceName string) []*Service {
	var services []*Service = nil
	douyinServices.RLock()
	services, _ = douyinServices.services[serviceName]
	douyinServices.RUnlock()
	return services
}

// 监视服务端的任何变化及时更新
func WatchServiceName(serviceName string, wait *sync.WaitGroup) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcd.GetEtcdEndpoints(),
		DialTimeout: etcd.DailTime,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cli.Close()
	getRes, err := cli.Get(context.Background(), serviceName+"/", clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
		return
	}
	//如果存在服务，则将服务存储到本进程的douyinService(相当于需要初始化我们的服务注册表格，使得客户端能使用服务发现的注册表来获得地址)
	if getRes.Count > 0 {
		var services []*Service
		for _, kv := range getRes.Kvs {
			parts := strings.Split(string(kv.Value), ":")
			if len(parts) != 3 {
				continue
			}
			s := &Service{
				Name:     serviceName,
				IP:       parts[0],
				Port:     parts[1],
				Protocol: parts[2],
			}
			services = append(services, s)
		}
		douyinServices.Lock()
		douyinServices.services[serviceName] = services
		douyinServices.Unlock()
	}

	wait.Done()
	//	开启监视
	rch := cli.Watch(context.Background(), serviceName+"/", clientv3.WithPrefix())
	for wres := range rch {
		for _, ev := range wres.Events {
			parts := strings.Split(string(ev.Kv.Key), "/")
			if len(parts) != 2 {
				continue
			}
			addr := parts[1]
			switch ev.Type {
			case clientv3.EventTypeDelete:
				douyinServices.Lock()
				var updatedServices []*Service
				for _, s := range douyinServices.services[serviceName] {
					if s.IP+":"+s.Port != addr {
						updatedServices = append(updatedServices, s)
					}
				}
				douyinServices.services[serviceName] = updatedServices
				douyinServices.Unlock()

			case clientv3.EventTypePut:
				partsPut := strings.Split(string(ev.Kv.Value), ":")
				if len(partsPut) != 3 {
					continue
				}
				newService := &Service{
					Name:     serviceName,
					IP:       partsPut[0],
					Port:     partsPut[1],
					Protocol: partsPut[2],
				}
				douyinServices.Lock()
				//检查服务实例是否存在，若存在就更新（就只有protocol能更新）
				found := false
				for i, s := range douyinServices.services[serviceName] {
					if s.IP == newService.IP && s.Port == newService.Port {
						douyinServices.services[serviceName][i] = newService
						found = true
						break
					}
				}
				if !found {
					douyinServices.services[serviceName] = append(douyinServices.services[serviceName], newService)
				}
				douyinServices.Unlock()
			}
		}
	}
}

func addPool() {

}
