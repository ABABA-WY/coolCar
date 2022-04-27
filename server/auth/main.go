package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	auth "coolcar/auth/auth"
	"coolcar/auth/auth/token"
	"coolcar/auth/dao"
	"coolcar/auth/wechat"
	"coolcar/shared/server"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io/ioutil"
	"os"
	"time"
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

	//开文件
	file, err := os.Open("auth/private.key")
	if err != nil {
		logger.Fatal("cannot open private key file ", zap.Error(err))
	}
	//读文件
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		logger.Fatal("cannot read private key ", zap.Error(err))
	}
	//翻译
	key, err := jwt.ParseRSAPrivateKeyFromPEM(bytes)
	if err != nil {
		logger.Fatal("cannot parse private key ", zap.Error(err))
	}

	logger.Sugar().Fatal(server.RunGRPCServer(&server.GRPCConfig{
		Name:    "auth",
		TCPAddr: "127.0.0.1:8001",
		RegisterFunc: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				Logeer:         logger,
				Mongo:          dao.NewMongo(client.Database("coolcar"), "account"),
				TokenExpire:    10 * time.Second,
				TokenGenerator: token.NewJWTTokenGen("coolcar/auth", key),
				OpenIDResolver: &wechat.Service{
					AppID:     "wx2d6e57c8ad43a3ff",
					AppSecret: "0be37f759a6f299dec8f691cbcfdd873",
				},
			})
		},
		Logger: logger,
	}))
}
