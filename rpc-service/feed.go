package main

import (
	"flag"
	"fmt"
	"github.com/RaymondCode/simple-demo/common"
	"github.com/RaymondCode/simple-demo/config"
	"github.com/RaymondCode/simple-demo/registry"
	pb "github.com/RaymondCode/simple-demo/rpc-service/proto"
	rpcService "github.com/RaymondCode/simple-demo/rpc-service/service"
	"github.com/RaymondCode/simple-demo/utils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strings"
)

func main() {
	var ip string
	var port string

	// 使用flag包定义命令行参数
	flag.StringVar(&ip, "i", strings.Split(config.FEED_SERVER_ADDR, ":")[0], "The IP address to listen on")
	flag.StringVar(&port, "p", strings.Split(config.FEED_SERVER_ADDR, ":")[1], "The port number to listen on")

	// 解析命令行参数
	flag.Parse()

	utils.InitConfig()
	common.InitDB()
	listen, err := net.Listen("tcp", ip+":"+port)
	if err != nil {
		fmt.Printf("无法启动监听：%v\n", err)
		return
	}

	// 创建 gRPC 服务器对象
	grpcServer := grpc.NewServer(grpc.Creds(insecure.NewCredentials()))

	// 在 gRPC 服务器上注册服务（本地注册）
	pb.RegisterVideoFeedServiceServer(grpcServer, &rpcService.VideoFeedService{})
	//注册中心注册
	s := &registry.Service{
		Name:     "feed",
		IP:       ip,
		Port:     port,
		Protocol: "grpc",
	}
	go registry.ServiceRegister(s)

	// 启动 gRPC 服务
	fmt.Println("启动视频流 gRPC 服务...")
	if err := grpcServer.Serve(listen); err != nil {
		fmt.Printf("启动视频流 gRPC 服务失败：%v\n", err)
		return
	}
}
