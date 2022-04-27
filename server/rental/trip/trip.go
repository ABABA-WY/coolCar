package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	auth "coolcar/shared/auth"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Logger *zap.Logger
	rentalpb.UnimplementedTripServiceServer
}

// CreateTrip
//type TripServiceServer interface {
//	CreateTrip(context.Context, *CreateTripRequest) (*CreateTripResponse, error)
//	mustEmbedUnimplementedTripServiceServer()
//}
func (s *Service) CreateTrip(c context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.CreateTripResponse, error) {
	//get account id from context
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	s.Logger.Info("create trip", zap.String("start:", request.Code), zap.String("accountID", aid))
	return nil, status.Error(codes.Unimplemented, "")
}
