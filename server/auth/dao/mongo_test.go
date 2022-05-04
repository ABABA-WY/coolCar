package dao

import (
	"context"
	id2 "coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/mongo/testing"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"os"
	"testing"
)

const containerPort = "27017/tcp"

var mongoURI = "mongodb://admin:123456@localhost:27017"

func TestResolveAccountID(t *testing.T) {
	fmt.Println(mongoURI)
	c := context.Background()
	client, err := mongotesting.NewMongoClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v\n", err)
	}
	m := NewMongo(client.Database("coolcar"))
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id2.AccountID("625905883b4ae1a7857692aa")),
			openIDField:        "456",
		},
		bson.M{
			mgutil.IDFieldName: objid.MustFromID(id2.AccountID("625905883b4ae1a7857bbbbb")),
			openIDField:        "789",
		},
	})
	if err != nil {
		t.Fatalf("cannot InsertMany:%v\n", err)
	}
	//m.newObjectID = func() primitive.ObjectID {
	//	objID := mustObjID("625905883b4ae1a7857692aa")
	//	return objID
	//}
	id, err := m.ResolveAccountID(c, "789")
	if err != nil {
		t.Error(err)
	} else {
		want := "625905883b4ae1a7857bbbbb"
		if want != id.String() {
			t.Errorf("not this openID")
		}

	}
	fmt.Println("ID:", id)
}

func TestMain(m *testing.M) {
	fmt.Println("main")
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
