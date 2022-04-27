package server

import (
	auth "coolcar/shared/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

type GRPCConfig struct {
	Name           string
	TCPAddr        string
	AuthPublicFile string
	RegisterFunc   func(*grpc.Server)
	Logger         *zap.Logger
}

// RunGRPCServer run grpc server
func RunGRPCServer(g *GRPCConfig) error {
	nameFile := zap.String("name", g.Name)
	listen, err := net.Listen("tcp", g.TCPAddr)
	if err != nil {
		g.Logger.Fatal("can't listen ", nameFile, zap.Error(err))
	}
	//NewServer
	var opts []grpc.ServerOption
	if g.AuthPublicFile != "" {
		in, err := auth.Interceptor(g.AuthPublicFile)
		if err != nil {
			g.Logger.Fatal("cannot create auth interceptor ", nameFile, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}
	//s := grpc.NewServer()
	s := grpc.NewServer(opts...) //添加拦截池,仅限opts不为空
	//业务逻辑
	g.RegisterFunc(s)
	g.Logger.Sugar().Infof("server started! name:%s, addr:%s", g.Name, g.TCPAddr)
	return s.Serve(listen)
}
