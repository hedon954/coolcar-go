package shared_mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	IDFieldName       = "_id"
	UpdateAtFieldName = "updateat"
)

type IDField struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdateAtField struct {
	UpdateAt int64 `bson:"updateat"`
}

// NewObjID generates a new object id
var NewObjID = primitive.NewObjectID

// UpdateAt returns a  value suitable for UpdateAt field
var UpdateAt = func() int64 {
	return time.Now().UnixNano()
}

// Set returns a $set update document
func Set(v interface{}) bson.M {
	return bson.M{
		"$set": v,
	}
}

// SetOnInsert returns a $setOnInsert update document
func SetOnInsert(v interface{}) bson.M {
	return bson.M{
		"$setOnInsert": v,
	}
}
