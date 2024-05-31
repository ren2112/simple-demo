package common

import (
	grpc_client_pool "github.com/RaymondCode/simple-demo/grpc-client-pool"
	"github.com/RaymondCode/simple-demo/registry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	ConnFeedPool     *grpc_client_pool.ClientFeedPool
	ConnUserPool     *grpc_client_pool.ClientUserPool
	ConnPublishPool  *grpc_client_pool.ClientPublishPool
	ConnRelationPool *grpc_client_pool.ClientRelationPool
	ConnFavoritePool *grpc_client_pool.ClientFavoritePool
	ConnCommentPool  *grpc_client_pool.ClientCommentPool
	ConnMessagePool  *grpc_client_pool.ClientMessagePool
	ConnFriendPool   *grpc_client_pool.ClientFriendPool
)

// 初始化所有连接池
func InitAllConnPool() {
	initializeFeedConnectionPool()
	initializeUserConnectionPool()
	initializePublishConnectionPool()
	initializeRelationConnectionPool()
	initializeFavoriteConnectionPool()
	initializeCommentConnectionPool()
	initializeMessageConnectionPool()
	initializeFriendConnectionPool()
}

// 初始化Feed连接池
func initializeFeedConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("feed")
	addr := s.IP + ":" + s.Port
	ConnFeedPool, err = grpc_client_pool.GetFeedPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化User连接池
func initializeUserConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("user")
	addr := s.IP + ":" + s.Port
	ConnUserPool, err = grpc_client_pool.GetUserPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化Publish连接池
func initializePublishConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("publish")
	addr := s.IP + ":" + s.Port
	ConnPublishPool, err = grpc_client_pool.GetPublishPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化Relation连接池
func initializeRelationConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("relation")
	addr := s.IP + ":" + s.Port
	ConnRelationPool, err = grpc_client_pool.GetRelationPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化favorite连接池
func initializeFavoriteConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("favorite")
	addr := s.IP + ":" + s.Port
	ConnFavoritePool, err = grpc_client_pool.GetFavoritePool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化连接池
func initializeCommentConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("comment")
	addr := s.IP + ":" + s.Port
	ConnCommentPool, err = grpc_client_pool.GetCommentPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化连接池
func initializeMessageConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("message")
	addr := s.IP + ":" + s.Port
	ConnMessagePool, err = grpc_client_pool.GetMessagePool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化连接池
func initializeFriendConnectionPool() {
	var err error
	s := registry.ServiceDiscovery("friend")
	addr := s.IP + ":" + s.Port
	ConnFriendPool, err = grpc_client_pool.GetFriendPool(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}
