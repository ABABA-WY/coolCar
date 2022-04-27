package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	trip "coolcar/rental/trip"
	server "coolcar/shared/server"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	//用语法糖来处理错误
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:           "rental",
		TCPAddr:        "127.0.0.1:8002",
		AuthPublicFile: "shared/auth/public.key",
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
			})
		},
		Logger: logger,
	}))
}
