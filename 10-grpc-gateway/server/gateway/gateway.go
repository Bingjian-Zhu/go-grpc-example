package gateway

import (
	"context"
	"crypto/tls"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	pb "go-grpc-example/10-grpc-gateway/proto"
	"go-grpc-example/10-grpc-gateway/server/swagger"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// ProvideHTTP 把gRPC服务转成HTTP服务，让gRPC同时支持HTTP
func ProvideHTTP(endpoint string, grpcServer *grpc.Server) *http.Server {
	ctx := context.Background()
	//获取证书
	creds, err := credentials.NewClientTLSFromFile("../tls/server.pem", "go-grpc-example")
	if err != nil {
		log.Fatalf("Failed to create TLS credentials %v", err)
	}
	//添加证书
	dopts := []grpc.DialOption{grpc.WithTransportCredentials(creds)}
	//新建gwmux，它是grpc-gateway的请求复用器。它将http请求与模式匹配，并调用相应的处理程序。
	gwmux := runtime.NewServeMux()
	//将服务的http处理程序注册到gwmux。处理程序通过endpoint转发请求到grpc端点
	err = pb.RegisterSimpleHandlerFromEndpoint(ctx, gwmux, endpoint, dopts)
	if err != nil {
		log.Fatalf("Register Endpoint err: %v", err)
	}
	//新建mux，它是http的请求复用器
	mux := http.NewServeMux()
	//注册gwmux
	mux.Handle("/", gwmux)
	//注册swagger
	mux.HandleFunc("/swagger/", swagger.ServeSwaggerFile)
	swagger.ServeSwaggerUI(mux)
	log.Println(endpoint + " HTTP.Listing whth TLS and token...")
	return &http.Server{
		Addr:      endpoint,
		Handler:   grpcHandlerFunc(grpcServer, mux),
		TLSConfig: getTLSConfig(),
	}
}

// grpcHandlerFunc 根据不同的请求重定向到指定的Handler处理
func grpcHandlerFunc(grpcServer *grpc.Server, otherHandler http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}

// getTLSConfig获取TLS配置
func getTLSConfig() *tls.Config {
	cert, _ := ioutil.ReadFile("../tls/server.pem")
	key, _ := ioutil.ReadFile("../tls/server.key")
	var demoKeyPair *tls.Certificate
	pair, err := tls.X509KeyPair(cert, key)
	if err != nil {
		grpclog.Fatalf("TLS KeyPair err: %v\n", err)
	}
	demoKeyPair = &pair
	return &tls.Config{
		Certificates: []tls.Certificate{*demoKeyPair},
		NextProtos:   []string{http2.NextProtoTLS}, // HTTP2 TLS支持
	}
}
