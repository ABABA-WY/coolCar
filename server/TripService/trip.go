package TripService

import (
	"context"
	trippb "coolcar/proto/gen/go"
	rentalpb "coolcar/rental/api/gen/v1"
)

//type TripServiceServer interface {
//	GetTrip(context.Context, *GetTripRequest) (*GetTripResponse, error)
//	mustEmbedUnimplementedTripServiceServer()
//}

type Service struct {
}

func (s *Service) CreateTrip(ctx context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.CreateTripResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) mustEmbedUnimplementedTripServiceServer() {
	//TODO implement me

}

func (s *Service) GetTrip(c context.Context, req *trippb.GetTripRequest) (*trippb.GetTripResponse, error) {
	return &trippb.GetTripResponse{
		Id: req.Id,
		Trip: &trippb.Trip{
			Start:       "aaa",
			End:         "bbb",
			DurationSec: 12,
			FeeCent:     34,
			StartPos: &trippb.Location{
				Latitude:  12,
				Longitude: 34,
			},
			EndPos: &trippb.Location{
				Latitude:  56,
				Longitude: 78,
			},
			PathLocation: []*trippb.Location{
				{
					Latitude:  22,
					Longitude: 33,
				},
				{
					Latitude:  33,
					Longitude: 44,
				},
			},
			Status: trippb.TripStatus_paid,
		},
	}, nil
}
