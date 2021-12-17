package dao

import (
	"context"
	"testing"

	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/mongo/objid"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var mongoURI string

func TestCreateTrip(t *testing.T) {
	c := context.Background()

	mongoURI = "mongodb://localhost:27017"

	// 获取客户端对象
	mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	// 获取集合操作对象
	m := NewMongo(mc.Database("coolcar"))

	// 创建行程
	tripRecord, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: "accountID",
		CarId:     "car1",
		Start: &rentalpb.LocationStatus{
			PositionName: "startpoint",
			Location: &rentalpb.Location{
				Latitude:  30,
				Longitude: 120,
			},
		},
		End: &rentalpb.LocationStatus{
			PositionName: "endpoint",
			Location: &rentalpb.Location{
				Latitude:  35,
				Longitude: 115,
			},
		},
		Status: rentalpb.TripStatus_FINISHED,
	})

	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}

	t.Errorf("%+v", tripRecord)

	got, err := m.GetTrip(c, objid.ToTripID(tripRecord.ID), "accountID")
	if err != nil {
		t.Fatalf("cannot get trip: %v", err)
	}
	t.Errorf("%+v", got)
}

// // 必须这么命名
// func TestMain(m *testing.M) {
// 	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
// }

// // 生成固定的 ObjectID
// func MustObjID(hex string) primitive.ObjectID {
// 	objID, _ := primitive.ObjectIDFromHex(hex)
// 	return objID
// }
