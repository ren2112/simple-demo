package grpc_client_pool

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"sync"
)

type ClientFeedPool struct {
	sync.Pool
}

func GetFeedPool(target string, opts ...grpc.DialOption) (*ClientFeedPool, error) {
	return &ClientFeedPool{
		Pool: sync.Pool{
			New: func() interface{} {
				conn, err := grpc.Dial(target, opts...)
				if err != nil {
					log.Fatal(err)
				}
				return conn
			},
		},
	}, nil
}

func (c *ClientFeedPool) Get() *grpc.ClientConn {
	conn := c.Pool.Get().(*grpc.ClientConn)
	if conn == nil || conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		if conn != nil {
			conn.Close()
		}
		conn = c.Pool.New().(*grpc.ClientConn)
	}
	return conn
}

func (c *ClientFeedPool) Put(conn *grpc.ClientConn) {
	if conn == nil {
		return
	}
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		return
	}
	c.Pool.Put(conn)
}
