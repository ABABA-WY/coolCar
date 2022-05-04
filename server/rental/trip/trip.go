package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	auth "coolcar/shared/auth"
	"coolcar/shared/id"
	"fmt"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	Mongo  *dao.Mongo
	Logger *zap.Logger
	rentalpb.UnimplementedTripServiceServer
	ProfileManager ProfileManager
	CarManager     CarManager
	POIManager     POIManager
}

//ProfileManager访问控制列表(ACL)
//for profile verification logic
type ProfileManager interface {
	Verify(context.Context, id.AccountID) (id.IdentityID, error)
}

//define the acl for car manager
type CarManager interface {
	//认证
	Verify(context.Context, id.CarID, *rentalpb.Location) error
	//检查是否被锁
	Unlock(context.Context, id.CarID) error
}

//POI point of interest
//地标
type POIManager interface {
	Resolve(context.Context, *rentalpb.Location) (string, error)
}

func (s *Service) CreateTrip(c context.Context, request *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	//get account id from context

	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	s.Logger.Info("create trip", zap.String("accountID", aid.String()))
	//return nil,status.Errorf(codes.OK,"")
	//验证驾驶者身份
	iID, err := s.ProfileManager.Verify(c, aid)
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}
	s.Logger.Info(iID.String())
	//fmt.Println(iID)
	//要将iID保存，方便以后查找
	//检查车辆状态
	carID := id.CarID(request.GetCarId())
	err = s.CarManager.Verify(c, carID, request.Start) //要求知道当前位置
	if err != nil {
		return nil, status.Errorf(codes.FailedPrecondition, err.Error())
	}
	s.Logger.Info(carID.String())
	//创建行程：写入数据库，开始计费。
	//上次更新就是起始位置
	ls := s.calcCurrentStatus(c, &rentalpb.LocationStatus{
		Location:     request.Start,
		TimestampSec: nowFunc(),
	}, request.Start)
	trip, err := s.Mongo.CreateTrip(c, &rentalpb.Trip{
		AccountId:  aid.String(),
		CarId:      carID.String(),
		IdentityId: iID.String(),
		TripStatus: rentalpb.TripStatus_IN_PROGRESS,
		Start:      ls,
		Current:    ls,
	})
	if err != nil {
		s.Logger.Warn("cannot create trip", zap.Error(err))
		return nil, status.Errorf(codes.AlreadyExists, err.Error())
	}
	//fmt.Println(trip)
	//车辆开锁
	//在后台开锁
	go func() {
		err = s.CarManager.Unlock(c, carID)
		if err != nil {
			s.Logger.Error("cannot unlock car", zap.Error(err))
		}
	}()
	//fmt.Println(trip.ID.Hex())
	//返回一个trip实体
	return &rentalpb.TripEntity{
		Trip: trip.Trip,
		Id:   trip.ID.Hex(),
	}, nil
}

func (s *Service) GetTrip(c context.Context, request *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	tr, err := s.Mongo.GetTrip(c, id.TripID(request.Tid), id.AccountID(request.Aid))
	if err != nil {
		return nil, err
	}
	return tr.Trip, err
}
func (s *Service) GetTrips(c context.Context, request *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	return nil, nil
}

// UpdateTrip update trip
func (s *Service) UpdateTrip(c context.Context, request *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	if request.Id == "" {
		return nil, status.Errorf(codes.InvalidArgument, "")
	}
	tid := id.TripID(request.Id)
	//fmt.Println(aid)
	//fmt.Println(tid)
	trip, err := s.Mongo.GetTrip(c, tid, aid)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, err.Error())
	}
	if trip.Trip.Current == nil {
		s.Logger.Error("trip without current set", zap.String("id", aid.String()))
		return nil, status.Error(codes.Internal, "")
	}
	//更新当前信息

	cur := trip.Trip.Current.Location //数据库中的信息
	if request.Current != nil {
		cur = request.Current

	}
	trip.Trip.Current = s.calcCurrentStatus(c, trip.Trip.Current, cur)
	//是否结束
	if request.EndTrip {
		trip.Trip.End = trip.Trip.Current
		trip.Trip.TripStatus = rentalpb.TripStatus_FINISHED
	}
	//UpdateAt乐观锁
	err = s.Mongo.UpdateTrip(c, tid, aid, trip.UpdateAt, trip.Trip)
	if err != nil {
		return nil, err
	}
	return trip.Trip, nil
}

var nowFunc = func() int64 {
	return time.Now().Unix()
}

const (
	centPerSec = 0.7
	kmPerSec   = 0.02
)

func (s *Service) calcCurrentStatus(c context.Context, last *rentalpb.LocationStatus, current *rentalpb.Location) *rentalpb.LocationStatus {
	now := nowFunc()
	elapsedSec := float64(now - last.TimestampSec) //距离上次的时间
	poi, err := s.POIManager.Resolve(c, current)
	//s.Logger.Info( request.Start.String())
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.String("location", current.String()))
	}
	return &rentalpb.LocationStatus{
		Location:     current,
		FeeCent:      last.FeeCent + int32(centPerSec*elapsedSec),
		KmDriven:     last.KmDriven + kmPerSec*elapsedSec,
		TimestampSec: now,
		PoiName:      poi,
	}
}
