package objid

import (
	shared_id "coolcar/shared/id"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// FromID coverts an id to object id
func FromID(id fmt.Stringer) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(id.String())
}

// MustFromID coverts an id to object id without error
func MustFromID(id fmt.Stringer) primitive.ObjectID {
	objid, err := primitive.ObjectIDFromHex(id.String())
	if err != nil {
		panic(err)
	}
	return objid
}

// ToAccountID converts object id to account id
func ToAccountID(objid primitive.ObjectID) shared_id.AccountID {
	return shared_id.AccountID(objid.Hex())
}

// ToTripID converts object id to account id
func ToTripID(objid primitive.ObjectID) shared_id.TripID {
	return shared_id.TripID(objid.Hex())
}
