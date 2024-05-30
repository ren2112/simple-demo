package registry

import (
	"context"
	"github.com/RaymondCode/simple-demo/etcd"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
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
	var grantLease bool
	var leaseId clientv3.LeaseID
	ctx := context.Background()

	//	查找etcd服务端是否存在服务
	getRes, err := cli.Get(ctx, s.Name, clientv3.WithCountOnly())
	if err != nil {
		log.Fatal(err)
	}
	if getRes.Count == 0 {
		grantLease = true
	}
	if grantLease {
		leaseRes, err := cli.Grant(ctx, 10)
		if err != nil {
			log.Fatal(err)
		}
		leaseId = leaseRes.ID
	}

	//	开启事务，进行注册
	kv := clientv3.NewKV(cli)
	txn := kv.Txn(ctx)
	_, err = txn.If(clientv3.Compare(clientv3.CreateRevision(s.Name), "=", 0)).
		Then(
			clientv3.OpPut(s.Name, s.Name, clientv3.WithLease(leaseId)),
			clientv3.OpPut(s.Name+".ip", s.IP, clientv3.WithLease(leaseId)),
			clientv3.OpPut(s.Name+".port", s.Port, clientv3.WithLease(leaseId)),
			clientv3.OpPut(s.Name+".protocol", s.Protocol, clientv3.WithLease(leaseId)),
		).
		Else(
			clientv3.OpPut(s.Name, s.Name, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".ip", s.IP, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".port", s.Port, clientv3.WithIgnoreLease()),
			clientv3.OpPut(s.Name+".protocol", s.Protocol, clientv3.WithIgnoreLease()),
		).
		Commit()

	if err != nil {
		log.Fatal(err)
	}
	if grantLease {
		leaseKeepalive, err := cli.KeepAlive(ctx, leaseId)
		if err != nil {
			log.Fatal(err)
		}
		//需要及时获取反馈消息，否则当做挂死
		for range leaseKeepalive {
		}
	}
}

type Services struct {
	services map[string]*Service
	sync.RWMutex
}

var douyinServices = &Services{
	services: map[string]*Service{},
}

// 服务发现
func ServiceDiscovery(serviceName string) *Service {
	var s *Service = nil
	douyinServices.RLock()
	s, _ = douyinServices.services[serviceName]
	douyinServices.RUnlock()
	return s
}

// 监视服务端的任何变化及时更新
func WatchServiceName(serviceName string) {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcd.GetEtcdEndpoints(),
		DialTimeout: etcd.DailTime,
	})
	if err != nil {
		log.Fatal(err)
		return
	}
	defer cli.Close()
	getRes, err := cli.Get(context.Background(), serviceName, clientv3.WithPrefix())
	if err != nil {
		log.Fatal(err)
		return
	}
	//如果存在服务，则将服务存储到本进程的douyinService
	if getRes.Count > 0 {
		mp := sliceToMap(getRes.Kvs)
		s := &Service{}
		if kv, ok := mp[serviceName]; ok {
			s.Name = string(kv.Value)
		}
		if kv, ok := mp[serviceName+".ip"]; ok {
			s.IP = string(kv.Value)
		}
		if kv, ok := mp[serviceName+".port"]; ok {
			s.Port = string(kv.Value)
		}
		if kv, ok := mp[serviceName+".protocol"]; ok {
			s.Protocol = string(kv.Value)
		}
		douyinServices.Lock()
		douyinServices.services[serviceName] = s
		douyinServices.Unlock()
	}

	//	开启监视
	rch := cli.Watch(context.Background(), serviceName, clientv3.WithPrefix())
	for wres := range rch {
		for _, ev := range wres.Events {
			if ev.Type == clientv3.EventTypeDelete {
				douyinServices.Lock()
				delete(douyinServices.services, serviceName)
				douyinServices.Unlock()
			}
			if ev.Type == clientv3.EventTypePut {
				douyinServices.Lock()
				if _, ok := douyinServices.services[serviceName]; !ok {
					douyinServices.services[serviceName] = &Service{}
				}
				switch string(ev.Kv.Key) {
				case serviceName:
					douyinServices.services[serviceName].Name = string(ev.Kv.Value)
				case serviceName + ".ip":
					douyinServices.services[serviceName].IP = string(ev.Kv.Value)
				case serviceName + ".port":
					douyinServices.services[serviceName].Port = string(ev.Kv.Value)
				case serviceName + ".protocol":
					douyinServices.services[serviceName].Protocol = string(ev.Kv.Value)
				}
				douyinServices.Unlock()
			}
		}
	}
}

func sliceToMap(list []*mvccpb.KeyValue) map[string]*mvccpb.KeyValue {
	mp := make(map[string]*mvccpb.KeyValue, 0)
	for _, item := range list {
		mp[string(item.Key)] = item
	}
	return mp
}
