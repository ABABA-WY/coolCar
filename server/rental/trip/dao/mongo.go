package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/id"
	mgutil "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
	statusField    = tripField + ".tripstatus"
	idField        = "idfield._id"
)

var IDField primitive.ObjectID

type Mongo struct {
	col *mongo.Collection
	//newObjectID func() primitive.ObjectID
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("trip"),
		//newObjectID: primitive.NewObjectID, //生成newobjectid
	}
}

type TripRecord struct {
	mgutil.IDField
	mgutil.UpdateAtField
	Trip *rentalpb.Trip `bson:"trip"`
}

func (m *Mongo) DeleteAllTrip(c context.Context) error {
	_, err := m.col.DeleteMany(context.TODO(), bson.D{{}})
	if err != nil {
		return err
	}
	return nil
}

// CreateTrip creates a trip
func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {

	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgutil.NewObjectID()
	r.UpdateAt = mgutil.UpdateAt()
	_, err := m.col.InsertOne(c, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
func (m *Mongo) GetTrip(c context.Context, tid id.TripID, accountID id.AccountID) (*TripRecord, error) {
	objID, err := objid.FromID(tid)
	if err != nil {
		return nil, fmt.Errorf("id invalid:%v\n", err)
	}
	//fmt.Println(tid)
	res := m.col.FindOne(c, bson.M{
		idField:        objID,
		accountIDField: accountID,
	})
	if err := res.Err(); err != nil {
		return nil, fmt.Errorf("find error:%v", err)
	}
	var tr TripRecord
	err = res.Decode(&tr)
	if err != nil {
		return nil, fmt.Errorf("cannot decode:%v", err)
	}
	return &tr, err
}

// GetTrips get trips
//if status not specified ,return all trips
func (m *Mongo) GetTrips(c context.Context, accountID id.AccountID, status rentalpb.TripStatus) ([]*TripRecord, error) {
	filter := bson.M{
		accountIDField: accountID.String(),
	}
	if status != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter[statusField] = status
	}
	find, err := m.col.Find(c, filter)
	if err != nil {
		return nil, err
	}
	var trips []*TripRecord
	for find.Next(c) {
		var tr TripRecord
		err := find.Decode(&tr)
		if err != nil {
			return nil, err
		}
		trips = append(trips, &tr)
	}
	return trips, nil
}

// UpdateTrip update trip
func (m *Mongo) UpdateTrip(c context.Context, tid id.TripID, aid id.AccountID, updateAt int64, trip *rentalpb.Trip) error {
	objID, err := objid.FromID(tid)
	fmt.Println(tid)
	if err != nil {
		return fmt.Errorf("id invalid:%v\n", err)
	}
	newUpdateAt := mgutil.UpdateAt()
	res, err := m.col.UpdateOne(c, bson.M{
		idField:                  objID,
		accountIDField:           aid.String(),
		mgutil.UpdateAtFieldName: updateAt,
	}, mgutil.Set(bson.M{
		tripField:                trip,
		mgutil.UpdateAtFieldName: newUpdateAt,
	}))
	if err != nil {
		return err
	}
	//The number of documents matched by the filter.
	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}
