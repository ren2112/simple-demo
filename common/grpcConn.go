package common

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
)

var (
	onceFeed     sync.Once
	onceUser     sync.Once
	oncePublish  sync.Once
	onceRelation sync.Once
	onceFavorite sync.Once
	onceComment  sync.Once

	connFeedPool     []*grpc.ClientConn
	connUserPool     []*grpc.ClientConn
	connPublishPool  []*grpc.ClientConn
	connRelationPool []*grpc.ClientConn
	connFavoritePool []*grpc.ClientConn
	connCommentPool  []*grpc.ClientConn
)

// 初始化连接池
func initializeFeedConnectionPool() {
	// 初始化连接池中的连接
	for i := 0; i < 10; i++ { // 这里可以根据需要设置连接池中连接的数量
		conn, err := grpc.Dial("127.0.0.1:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接：%v", err)
		}
		connFeedPool = append(connFeedPool, conn)
	}
}

// 初始化连接池
func initializeUserConnectionPool() {
	// 初始化连接池中的连接
	for i := 0; i < 10; i++ { // 这里可以根据需要设置连接池中连接的数量
		conn, err := grpc.Dial("127.0.0.1:9092", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接：%v", err)
		}
		connUserPool = append(connUserPool, conn)
	}
}

// 初始化连接池
func initializePublishConnectionPool() {
	// 初始化连接池中的连接
	for i := 0; i < 10; i++ { // 这里可以根据需要设置连接池中连接的数量
		conn, err := grpc.Dial("127.0.0.1:9093", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接：%v", err)
		}
		connPublishPool = append(connPublishPool, conn)
	}
}

func initializeRelationConnectionPool() {
	// 初始化连接池中的连接
	for i := 0; i < 10; i++ { // 这里可以根据需要设置连接池中连接的数量
		conn, err := grpc.Dial("127.0.0.1:9094", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接：%v", err)
		}
		connRelationPool = append(connRelationPool, conn)
	}
}

// 初始化连接池
func initializeFavoriteConnectionPool() {
	// 初始化连接池中的连接
	for i := 0; i < 10; i++ { // 这里可以根据需要设置连接池中连接的数量
		conn, err := grpc.Dial("127.0.0.1:9095", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接：%v", err)
		}
		connFavoritePool = append(connFavoritePool, conn)
	}
}

// 初始化连接池
func initializeCommentConnectionPool() {
	// 初始化连接池中的连接
	for i := 0; i < 10; i++ { // 这里可以根据需要设置连接池中连接的数量
		conn, err := grpc.Dial("127.0.0.1:9096", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接：%v", err)
		}
		connCommentPool = append(connCommentPool, conn)
	}
}

// 从连接池中获取连接
func GetFeedConnection() *grpc.ClientConn {
	onceFeed.Do(func() {
		initializeFeedConnectionPool()
	})
	return connFeedPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
}

func GetUserConnection() *grpc.ClientConn {
	onceUser.Do(func() {
		initializeUserConnectionPool()
	})
	return connUserPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
}

func GetPublishConnection() *grpc.ClientConn {
	oncePublish.Do(func() {
		initializePublishConnectionPool()
	})
	return connPublishPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
}

func GetRelationConnection() *grpc.ClientConn {
	onceRelation.Do(func() {
		initializeRelationConnectionPool()
	})
	return connRelationPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
}

func GetFavoriteConnection() *grpc.ClientConn {
	onceFavorite.Do(func() {
		initializeFavoriteConnectionPool()
	})
	return connFavoritePool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
}

func GetCommentConnection() *grpc.ClientConn {
	onceComment.Do(func() {
		initializeCommentConnectionPool()
	})
	return connCommentPool[0] // 返回连接池中的第一个连接，这里可以实现负载均衡或其他连接选择策略
}
