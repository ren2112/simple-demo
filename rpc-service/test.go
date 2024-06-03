package main

import (
	"context"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	redislock "github.com/jefferyjob/go-redislock"
	"time"
)

func main() {
	// Create a Redis client

	// Create a context for canceling lock operations
	ctx := context.Background()

	// Create a RedisLock object
	lock := redislock.New(ctx, common.RedisClient, "test_key", redislock.WithAutoRenew())

	// acquire lock
	err := lock.Lock()
	if err != nil {
		fmt.Println("lock acquisition failed：", err)
		return
	}
	defer lock.UnLock() // unlock

	// Perform tasks during lockdown
	// ...
	fmt.Println("正在执行")
	time.Sleep(10 * time.Second)
	fmt.Println("task execution completed")
}
