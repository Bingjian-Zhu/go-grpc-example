package recovery

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

var customFunc grpc_recovery.RecoveryHandlerFunc

func RecoveryInterceptor(p interface{}) (err error) {
	return grpc.Errorf(codes.Unknown, "panic triggered: %v", p)
}
