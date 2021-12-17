package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	shared_id "coolcar/shared/id"
	shared_mongo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
)

type Mongo struct {
	col *mongo.Collection
}

func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("trip"),
	}
}

type TripRecord struct {
	shared_mongo.IDField       `bson:"inline"`
	shared_mongo.UpdateAtField `bson:"inline"`
	Trip                       *rentalpb.Trip `bson:"trip"`
}

// CreateTrip creates a trip record in mongodb
func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = shared_mongo.NewObjID()
	r.UpdateAt = shared_mongo.UpdateAt()
	_, err := m.col.InsertOne(c, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// GetTrip returns a trip suitable for tripID and accountID
func (m *Mongo) GetTrip(c context.Context, id shared_id.TripID, accountID shared_id.AccountID) (*TripRecord, error) {
	objID, err := objid.FromID(id)
	if err != nil {
		return nil, fmt.Errorf("invalid id: %v", err)
	}
	res := m.col.FindOne(c, bson.M{
		shared_mongo.IDFieldName: objID,
		accountIDField:           accountID,
	})
	if err := res.Err(); err != nil {
		return nil, err
	}

	var t TripRecord
	err = res.Decode(&t)
	if err != nil {
		return nil, fmt.Errorf("cannot decode insertOne result: %v", err)
	}
	return &t, nil
}
