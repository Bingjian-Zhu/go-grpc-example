package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "go-grpc-example/6-grpc_deadlines/proto"
)

// SimpleService 定义我们的服务
type SimpleService struct{}

const (
	// Address 监听地址
	Address string = ":8000"
	// Network 网络通信协议
	Network string = "tcp"
)

func main() {
	// 监听本地端口
	listener, err := net.Listen(Network, Address)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	log.Println(Address + " net.Listing...")
	// 新建gRPC服务器实例
	grpcServer := grpc.NewServer()
	// 在gRPC服务器注册我们的服务
	pb.RegisterSimpleServer(grpcServer, &SimpleService{})

	//用服务器 Serve() 方法以及我们的端口信息区实现阻塞等待，直到进程被杀死或者 Stop() 被调用
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("grpcServer.Serve err: %v", err)
	}
}

// Route 实现Route方法
func (s *SimpleService) Route(ctx context.Context, req *pb.SimpleRequest) (*pb.SimpleResponse, error) {
	timeout := make(chan struct{}, 1)
	data := make(chan *pb.SimpleResponse, 1)
	go func() {
		time.Sleep(4 * time.Second)
		res := pb.SimpleResponse{
			Code:  200,
			Value: "hello " + req.Data,
		}
		log.Println("goroutine still running")
		data <- &res
	}()
	go func() {
		for {
			if ctx.Err() == context.Canceled {
				timeout <- struct{}{}
			}
		}
	}()
	select {
	case res := <-data:
		return res, nil
	case <-timeout:
		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	}
	// for n := 0; n <= 5; n++ {
	// 	if ctx.Err() == context.Canceled {
	// 		return nil, status.Errorf(codes.Canceled, "Client cancelled, abandoning.")
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }
	// res := pb.SimpleResponse{
	// 	Code:  200,
	// 	Value: "hello " + req.Data,
	// }
	// return &res, nil
}
