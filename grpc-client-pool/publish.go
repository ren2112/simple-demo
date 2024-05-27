package grpc_client_pool

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"log"
	"sync"
)

type ClientPublishPool struct {
	sync.Pool
}

func GetPublishPool(target string, opts ...grpc.DialOption) (*ClientPublishPool, error) {
	return &ClientPublishPool{
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

func (c *ClientPublishPool) Get() *grpc.ClientConn {
	conn := c.Pool.Get().(*grpc.ClientConn)
	if conn == nil || conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		if conn != nil {
			conn.Close()
		}
		conn = c.Pool.New().(*grpc.ClientConn)
	}
	return conn
}

func (c *ClientPublishPool) Put(conn *grpc.ClientConn) {
	if conn == nil {
		return
	}
	if conn.GetState() == connectivity.Shutdown || conn.GetState() == connectivity.TransientFailure {
		conn.Close()
		return
	}
	c.Pool.Put(conn)
}
