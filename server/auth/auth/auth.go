package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	authpb.UnimplementedAuthServiceServer
	Logeer         *zap.Logger
	OpenIDResolver OpenIDResolver
	Mongo          *dao.Mongo
	TokenGenerator TokenGenerator
	TokenExpire    time.Duration //有效期
}

// OpenIDResolver 获取ID
type OpenIDResolver interface {
	Resolver(code string) (string, error)
}

// TokenGenerator 获取签名
type TokenGenerator interface {
	GenerateToken(accountID string, expire time.Duration) (string, error)
}

func (s *Service) Login(c context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	//获取openid
	openID, err := s.OpenIDResolver.Resolver(req.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openid: %v", err)
	}
	s.Logeer.Info("received code:", zap.String("code", req.Code))
	//通过openID获取accountID
	accountID, err := s.Mongo.ResolveAccountID(c, openID)
	if err != nil {
		s.Logeer.Error("cannot resolve accountID", zap.Error(err)) //记日志
		return nil, status.Errorf(codes.Internal, "")
	}

	//获取签名
	token, err := s.TokenGenerator.GenerateToken(accountID, s.TokenExpire)
	if err != nil {
		s.Logeer.Error("cannot token generate:%v\n", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "cannot token generate")
	}
	return &authpb.LoginResponse{
		AccessToken: token,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}
