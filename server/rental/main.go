package main

import (
	"context"

	rentalpb "coolcar/rental/api/gen/v1"
	profile "coolcar/rental/profile"
	profileDao "coolcar/rental/profile/dao"
	trip "coolcar/rental/trip"
	"coolcar/rental/trip/client/car"
	"coolcar/rental/trip/client/poi"
	profileClient "coolcar/rental/trip/client/profile"
	tripDao "coolcar/rental/trip/dao"
	server "coolcar/shared/server"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	//连接数据库
	c := context.Background() //后台任务
	client, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://admin:123456@localhost:27017"))
	if err != nil {
		logger.Fatal("cannot connect mongodb", zap.Error(err))
	}
	//用语法糖来处理错误
	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:           "rental",
		TCPAddr:        "127.0.0.1:8002",
		AuthPublicFile: "shared/auth/public.key",
		RegisterFunc: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger:         logger,
				POIManager:     &poi.Manager{},
				ProfileManager: &profileClient.Manager{},
				CarManager:     &car.Manager{},
				Mongo:          tripDao.NewMongo(client.Database("coolcar")),
			})
			rentalpb.RegisterProfileServiceServer(s, &profile.Service{
				Logger: logger,
				Mongo:  profileDao.NewMongo(client.Database("coolcar")),
			})
		},
		Logger: logger,
	}))
}
