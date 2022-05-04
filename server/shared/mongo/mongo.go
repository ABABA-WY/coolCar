package mgutil

import (
	"coolcar/shared/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const IDFieldName = "_id"
const UpdateAtFieldName = "updateatfield.update_at"

type IDField struct {
	ID     primitive.ObjectID `bson:"_id"`
	OpenID string             `bson:"open_id"`
}
type UpdateAtField struct {
	UpdateAt int64 `bson:"update_at"`
}

// NewObjectID generate a new object id
var NewObjectID = primitive.NewObjectID

func NewObjectIDWithValue(id fmt.Stringer) {
	NewObjectID = func() primitive.ObjectID {
		return objid.MustFromID(id)
	}
}

var UpdateAt = func() int64 {
	return time.Now().UnixNano()
}

func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}
func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}
