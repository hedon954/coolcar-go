package dao

import (
	"context"
	"fmt"

	rentalpb "coolcar/rental/api/gen/v1"
	shared_id "coolcar/shared/id"
	shared_mongo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	tripField      = "trip"
	accountIDField = tripField + ".accountid"
	statusField    = tripField + ".status"
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

// GetTrips gets trips for the account by status
// If status is not specified, gets all trip for the account
func (m *Mongo) GetTrips(c context.Context, accountID shared_id.AccountID, status rentalpb.TripStatus) ([]*TripRecord, error) {
	filter := bson.M{
		accountIDField: accountID.String(),
	}

	if status != rentalpb.TripStatus_TS_NOT_SPECIFIED {
		filter[statusField] = status
	}

	res, err := m.col.Find(c, filter)
	if err != nil {
		return nil, err
	}

	var trips []*TripRecord
	for res.Next(c) {
		var trip TripRecord
		err = res.Decode(&trip)
		if err != nil {
			return nil, err
		}
		trips = append(trips, &trip)
	}
	return trips, nil
}

// UpdateTrip updates a trip
// Need Optimistic Lock by using updateAt field
func (m *Mongo) UpdateTrip(c context.Context, tripID shared_id.TripID, accountID shared_id.AccountID, updateAt int64, trip *rentalpb.Trip) error {
	oid, err := objid.FromID(tripID)
	if err != nil {
		return err
	}
	newUpdateAt := shared_mongo.UpdateAt()

	res, err := m.col.UpdateOne(c, bson.M{
		shared_mongo.IDFieldName:       oid,
		accountIDField:                 accountID.String(),
		shared_mongo.UpdateAtFieldName: updateAt,
	}, shared_mongo.Set(bson.M{
		tripField:                      trip,
		shared_mongo.UpdateAtFieldName: newUpdateAt,
	}))

	if err != nil {
		return err
	}

	if res.MatchedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}
