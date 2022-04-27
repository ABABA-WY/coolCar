package mongotesting

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const IDField = "_id"

type ObjectID struct {
	ID     primitive.ObjectID `bson:"_id"`
	OpenID string             `bson:"open_id"`
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
