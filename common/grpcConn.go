package common

import (
	"github.com/RaymondCode/simple-demo/config"
	grpc_client_pool "github.com/RaymondCode/simple-demo/grpc-client-pool"
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
	ConnFeedPool, err = grpc_client_pool.GetFeedPool(config.FEED_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化User连接池
func initializeUserConnectionPool() {
	var err error
	ConnUserPool, err = grpc_client_pool.GetUserPool(config.USER_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化Publish连接池
func initializePublishConnectionPool() {
	var err error
	ConnPublishPool, err = grpc_client_pool.GetPublishPool(config.PUBLISH_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化Relation连接池
func initializeRelationConnectionPool() {
	var err error
	ConnRelationPool, err = grpc_client_pool.GetRelationPool(config.RELATION_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化favorite连接池
func initializeFavoriteConnectionPool() {
	var err error
	ConnFavoritePool, err = grpc_client_pool.GetFavoritePool(config.FAVORITE_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化连接池
func initializeCommentConnectionPool() {
	var err error
	ConnCommentPool, err = grpc_client_pool.GetCommentPool(config.COMMENT_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化连接池
func initializeMessageConnectionPool() {
	var err error
	ConnMessagePool, err = grpc_client_pool.GetMessagePool(config.MESSAGE_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 初始化连接池
func initializeFriendConnectionPool() {
	var err error
	ConnFriendPool, err = grpc_client_pool.GetFriendPool(config.FRIEND_SERVER_ADDR, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
}

// 从连接池中获取连接
//func GetFeedConnection() *grpc.ClientConn {
//	conn := ConnFeedPool.Get()
//
//	//onceFeed.Do(func() {
//	//	initializeFeedConnectionPool()
//	//})
//	//return connFeedPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}

//func GetUserConnection() *grpc.ClientConn {
//	onceUser.Do(func() {
//		initializeUserConnectionPool()
//	})
//	return connUserPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}

//func GetPublishConnection() *grpc.ClientConn {
//	oncePublish.Do(func() {
//		initializePublishConnectionPool()
//	})
//	return connPublishPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}

//func GetRelationConnection() *grpc.ClientConn {
//	onceRelation.Do(func() {
//		initializeRelationConnectionPool()
//	})
//	return connRelationPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}
//
//func GetFavoriteConnection() *grpc.ClientConn {
//	onceFavorite.Do(func() {
//		initializeFavoriteConnectionPool()
//	})
//	return connFavoritePool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}
//
//func GetCommentConnection() *grpc.ClientConn {
//	onceComment.Do(func() {
//		initializeCommentConnectionPool()
//	})
//	return connCommentPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}
//
//func GetMessageConnection() *grpc.ClientConn {
//	onceMessage.Do(func() {
//		initializeMessageConnectionPool()
//	})
//	return connMessagePool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}
//
//func GetFriendConnection() *grpc.ClientConn {
//	onceFriend.Do(func() {
//		initializeFriendConnectionPool()
//	})
//	return connFriendPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
//}
