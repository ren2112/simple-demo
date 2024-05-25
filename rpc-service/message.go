package main

import (
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	rpcService "github.com/RaymondCode/simple-demo/rpc-service/service"
	"github.com/RaymondCode/simple-demo/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
)

func main() {
	utils.InitConfig()
	common.InitDB()
	listen, err := net.Listen("tcp", ":9097")
	if err != nil {
		fmt.Printf("无法启动监听：%v\n", err)
		return
	}

	// 创建 gRPC 服务器对象
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	// 在 gRPC 服务器上注册服务（本地注册）
	pb.RegisterChatServiceServer(grpcServer, &rpcService.MessageService{})

	// 启动 gRPC 服务
	fmt.Println("启动chat gRPC 服务...")
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("启动message gRPC 服务失败：%v\n", err)
		return
	}
}
