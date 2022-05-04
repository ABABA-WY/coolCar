package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"math/rand"
	"os"
	"testing"
	"time"
)

//ProfileManager访问控制列表(ACL)
//for profile verification logic
type profileManager struct {
	iID id.IdentityID
	err error
}

func (p *profileManager) Verify(c context.Context, aid id.AccountID) (id.IdentityID, error) {
	return p.iID, nil
}

//define the acl for car manager
type carManager struct {
	verr error
	uerr error
}

func (car *carManager) Verify(c context.Context, cid id.CarID, location *rentalpb.Location) error {
	return car.verr
}
func (car *carManager) Unlock(c context.Context, cid id.CarID) error {
	return car.uerr
}

func TestCreateTrip(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), "account1")
	pm := &profileManager{}
	cm := &carManager{}
	//建服务
	s := newService(c, t, pm, cm)

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  30,
			Longitude: 114,
		},
	}
	cases := []struct {
		name       string
		accountID  string
		tripID     string
		profileErr error
		carErr     error
		want       string
	}{
		{
			name:      "normal",
			tripID:    "625905883b4ae1a7857692aa",
			accountID: "a1",
		},
		{
			name:      "profile_err",
			accountID: "a2",
			tripID:    "625905883b4ae1a7857692ab",
		},
	}
	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			//生成object id
			mgutil.NewObjectIDWithValue(id.TripID(cc.tripID))
			c := auth.ContextWithAccountID(context.Background(), id.AccountID(cc.accountID))
			res, err := s.CreateTrip(c, req)
			if err != nil {
				t.Errorf("create trip error\n")
				return
			}
			if res.Id != cc.tripID {
				t.Fatal("invalid id\n")
			}
			fmt.Printf("%+v\n", res.Trip.Current)
		})
	}
}
func TestLifecycle(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), id.AccountID("account8"))
	tid := id.TripID("625905883b4ae1a7857692ab")
	mgutil.NewObjectIDWithValue(tid)
	pm := &profileManager{}
	cm := &carManager{}
	s := newService(c, t, pm, cm)
	//err := s.Mongo.DeleteAllTrip(c)
	//if err != nil {
	//	t.Fatalf("cannot delete:%v\n", err)
	//}
	cases := []struct {
		name string
		now  int64
		op   func() (*rentalpb.Trip, error)
	}{
		{
			name: "create_trip",
			now:  10000,
			op: func() (*rentalpb.Trip, error) {
				e, err := s.CreateTrip(c, &rentalpb.CreateTripRequest{
					CarId: "car1",
					Start: &rentalpb.Location{
						Latitude:  30,
						Longitude: 114,
					},
				})
				if err != nil {
					return nil, err
				}
				return e.Trip, nil
			},
		},
		{

			name: "update_trip",
			now:  20000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id: tid.String(),
					Current: &rentalpb.Location{
						Latitude:  50,
						Longitude: 60,
					},
				})
			},
		},
		{
			name: "finished_trip",
			now:  30000,
			op: func() (*rentalpb.Trip, error) {
				return s.UpdateTrip(c, &rentalpb.UpdateTripRequest{
					Id:      tid.String(),
					EndTrip: true,
				})
			},
		},
		{
			name: "query_trip",
			now:  40000,
			op: func() (*rentalpb.Trip, error) {
				return s.GetTrip(c, &rentalpb.GetTripRequest{
					Tid: tid.String(),
					Aid: "account8",
				})
			},
		},
	}
	for _, cc := range cases {
		nowFunc = func() int64 {
			return cc.now
		}
		rand.Seed(time.Now().UnixNano())
		trip, err := cc.op()
		if err != nil {
			t.Errorf("%s:operation failed:%v\n", cc.name, err)
			continue
		}
		b, err := json.Marshal(trip)
		if err != nil {
			t.Errorf("%s:json:%v\n", cc.name, err)
		}
		fmt.Println(string(b))
	}
}
func newService(c context.Context, t *testing.T, pm ProfileManager, cm CarManager) *Service {
	client, err := mongotesting.NewMongoClient(c)
	//client,err := mongotesting.NewDefaultMongoClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v\n", err)
	}
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	db := client.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("SetupIndexes error")
	}
	//建服务
	s := &Service{
		ProfileManager: pm,
		CarManager:     cm,
		POIManager:     &poi.Manager{},
		Mongo:          dao.NewMongo(client.Database("coolcar")),
		Logger:         logger,
	}
	return s
}
func TestMain(m *testing.M) {
	fmt.Println("main")
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
