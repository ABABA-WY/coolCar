package objid

import (
	"coolcar/shared/id"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FromID 将id转换成objected id
func FromID(id fmt.Stringer) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id.String())
}
func MustFromID(id fmt.Stringer) primitive.ObjectID {
	aid, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		panic(err)
	}
	return aid
}

// ToAccountID 将objected id转换成account id
func ToAccountID(oid primitive.ObjectID) id.AccountID {
	return id.AccountID(oid.Hex())
}

// ToTripID 将objected id转换成trip id
func ToTripID(oid primitive.ObjectID) id.TripID {
	return id.TripID(oid.Hex())
}
