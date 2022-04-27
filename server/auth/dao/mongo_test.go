package dao

import (
	"context"
	mongotesting "coolcar/shared/testing"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

const containerPort = "27017/tcp"

var mongoURI = "mongodb://admin:123456@localhost:27017"

func TestResolveAccountID(t *testing.T) {
	fmt.Println(mongoURI)
	c := context.Background()
	client, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect mongodb:%v\n", err)
	}
	m := NewMongo(client.Database("coolcar"), "account")
	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mongotesting.IDField: mustObjID("625905883b4ae1a7857692aa"),
			openIDField:          "456",
		},
		bson.M{
			mongotesting.IDField: mustObjID("625905883b4ae1a7857bbbbb"),
			openIDField:          "789",
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
		if want != id {
			t.Errorf("not this openID")
		}

	}
	fmt.Println("ID:", id)
}
func mustObjID(hex string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return objID
}
func TestMain(m *testing.M) {
	fmt.Println("main")
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
