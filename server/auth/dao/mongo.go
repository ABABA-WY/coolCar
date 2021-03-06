package dao

import (
	"context"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const openIDField = "open_id"

var IDField primitive.ObjectID

type Mongo struct {
	col         *mongo.Collection
	newObjectID func() primitive.ObjectID
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:         db.Collection("account"),
		newObjectID: primitive.NewObjectID, //生成newobjectid
	}
}

func (m *Mongo) ResolveAccountID(c context.Context, openID string) (id.AccountID, error) {
	insertID := mgutil.NewObjectID()
	//查到了就返回，查不到就插入
	res := m.col.FindOneAndUpdate(c, bson.M{
		openIDField: openID,
	}, mgutil.SetOnInsert(bson.M{
		mgutil.IDFieldName: insertID,
		openIDField:        openID,
	}), options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After), //返回更新后
	)
	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot FindOneAndUpdate:%v\n", err)
	}
	var row mgutil.IDField
	err := res.Decode(&row)
	if err != nil {
		return "", fmt.Errorf("cannot Decode:%v\n", err)
	}
	return objid.ToAccountID(row.ID), nil
}
