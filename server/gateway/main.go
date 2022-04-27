package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net/http"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()
	mux := runtime.NewServeMux(runtime.WithMarshalerOption(
		runtime.MIMEWildcard, &runtime.JSONPb{},
	))
	serverConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         "127.0.0.1:8001",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         "127.0.0.1:8002",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}
	for _, s := range serverConfig {
		//88fmt.Println(s.name)
		err := s.registerFunc(c, mux, s.addr, []grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			logger.Sugar().Fatalf("cannot register %s service: %v\n", s.name, err)
		}
	}

	//总监听端口
	addr := "127.0.0.1:8000"
	logger.Sugar().Info("grpc gateway started at:", addr)
	err = http.ListenAndServe(addr, mux)
	if err != nil {
		logger.Sugar().Fatalf("cannot ListenAndServe: %s, error:%v\n", addr, err)
	}

}
