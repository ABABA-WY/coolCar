package profile

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/profile/dao"
	"coolcar/shared/auth"
	"coolcar/shared/id"
	mongotesting "coolcar/shared/mongo/testing"
	"fmt"
	"go.uber.org/zap"
	"os"
	"testing"
)

func TestProfileLifecycle(t *testing.T) {
	c := auth.ContextWithAccountID(context.Background(), id.AccountID("account88"))
	s := newService(c, t)
	cases := []struct {
		name       string
		op         func() (*rentalpb.Profile, error)
		wantName   string
		wantStatus rentalpb.IdentityStatus
	}{
		{
			name: "get_empty",
			op: func() (*rentalpb.Profile, error) {
				return s.GetProfile(c, &rentalpb.GetProfileRequest{})
			},
		},
		{
			name: "submit",
			op: func() (*rentalpb.Profile, error) {
				return s.SubmitProfile(c, &rentalpb.Identity{
					Name: "sb",
				})
			},
			wantName:   "sb",
			wantStatus: rentalpb.IdentityStatus_PENDING,
		},
		{
			name: "verify",
			op: func() (*rentalpb.Profile, error) {
				p := &rentalpb.Profile{
					Identity: &rentalpb.Identity{
						Name: "sb",
					},
					IdentityStatus: rentalpb.IdentityStatus_VERIFIED,
				}
				err := s.Mongo.UpdateProfile(c, id.AccountID("account88"), rentalpb.IdentityStatus_PENDING, p)
				if err != nil {
					return nil, err
				}
				return p, nil
			},
			wantName:   "sb",
			wantStatus: rentalpb.IdentityStatus_VERIFIED,
		},
		{
			name: "clear",
			op: func() (*rentalpb.Profile, error) {
				return s.ClearProfile(c, &rentalpb.ClearProfileRequest{})
			},
			wantStatus: rentalpb.IdentityStatus_UNSUBMITTED,
		},
	}
	for _, cc := range cases {
		p, err := cc.op()
		if err != nil {
			t.Errorf("cc op error:%v\n", err)
		}
		if cc.wantName != "" && cc.wantName != p.Identity.Name {
			t.Errorf("wrong name,want:%s,got:%s", cc.wantName, p.Identity.Name)
		}
		if cc.wantStatus != p.IdentityStatus {
			t.Errorf("wrong status,want:%+v,got:%+v", cc.wantStatus, p.IdentityStatus)
		}
		fmt.Printf("\n%v\n", p)
	}
}

func newService(c context.Context, t *testing.T) *Service {
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

		Mongo:  dao.NewMongo(client.Database("coolcar")),
		Logger: logger,
	}
	return s
}
func TestMain(m *testing.M) {
	fmt.Println("main")
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
