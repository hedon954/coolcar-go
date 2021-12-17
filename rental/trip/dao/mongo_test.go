package dao

import (
	"context"
	"os"
	"testing"

	rentalpb "coolcar/rental/api/gen/v1"
	shared_id "coolcar/shared/id"
	shared_mongo "coolcar/shared/mongo"
	"coolcar/shared/mongo/objid"
	mongotesting "coolcar/shared/testing"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestCreateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connet mongodb: %v", err)
	}

	db := mc.Database("coolcar")
	err = mongotesting.SetupIndexes(c, db)
	if err != nil {
		t.Fatalf("cannot setup indexes: %v", err)
	}

	m := NewMongo(db)

	cases := []struct {
		name       string
		tripID     string
		accounID   string
		tripstatus rentalpb.TripStatus
		wantErr    bool
	}{
		{
			name:       "finished",
			tripID:     "61bb468dc648e913cded5652",
			accounID:   "account1",
			tripstatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "another_finished",
			tripID:     "61bb48d44e17e2f2040c1b59",
			accounID:   "account1",
			tripstatus: rentalpb.TripStatus_FINISHED,
		},
		{
			name:       "in_progress",
			tripID:     "61bb48e569239c93c16c96bd",
			accounID:   "account1",
			tripstatus: rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			name:       "anther_in_progress",
			tripID:     "61bb490936560e9b4c1dd2a3",
			accounID:   "account1",
			tripstatus: rentalpb.TripStatus_IN_PROGRESS,
			wantErr:    true,
		},
		{
			name:       "in_progress_by_another_account",
			tripID:     "61bb48d44e17e2f2040c1b29",
			accounID:   "account2",
			tripstatus: rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	for _, testCase := range cases {
		shared_mongo.NewObjID = func() primitive.ObjectID {
			return objid.MustFromID(shared_id.TripID(testCase.tripID))
		}
		tr, err := m.CreateTrip(c, &rentalpb.Trip{
			AccountId: testCase.accounID,
			Status:    testCase.tripstatus,
		})
		if testCase.wantErr {
			if err == nil {
				t.Errorf("%s: error expected; got none", testCase.name)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: error creating trip: %v", testCase.name, err)
			continue
		}
		if tr.ID.Hex() != testCase.tripID {
			t.Errorf("%s: incorrent trip id; want: %q, got: %q", testCase.name, testCase.tripID, tr.ID.Hex())
		}
	}
}

func TestGetTrip(t *testing.T) {
	c := context.Background()

	// 获取客户端对象
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}

	// 获取集合操作对象
	m := NewMongo(mc.Database("coolcar"))

	shared_mongo.NewObjID = primitive.NewObjectID

	// 创建行程
	tripRecord, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: "accountID2",
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
		Status: rentalpb.TripStatus_IN_PROGRESS,
	})

	if err != nil {
		t.Fatalf("cannot create trip: %v", err)
	}

	got, err := m.GetTrip(c, objid.ToTripID(tripRecord.ID), "accountID2")
	if err != nil {
		t.Fatalf("cannot get trip: %v", err)
	}

	if diff := cmp.Diff(tripRecord, got, protocmp.Transform()); diff != "" {
		t.Errorf("results differs; -want +got:%s", diff)
	}
}

func TestGetTrips(t *testing.T) {
	// 准备测试数据
	rows := []struct {
		id        shared_id.TripID
		accountID shared_id.AccountID
		status    rentalpb.TripStatus
	}{
		{
			id:        "61bb468dc648e913cded5652",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "61bb48d44e17e2f2040c1b59",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "61bb48d44e17e2f2040c1b58",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_FINISHED,
		},
		{
			id:        "61bb48d44e17e2f2040c1b56",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
		{
			id:        "61bb48d44e17e2f2040c1b52",
			accountID: "account_id_for_get_trips_1",
			status:    rentalpb.TripStatus_IN_PROGRESS,
		},
	}

	// 连接 mongodb
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))

	// 先插入数据
	for _, r := range rows {
		shared_mongo.NewObjIDWithValue(shared_id.TripID(r.id))

		_, err = m.CreateTrip(c, &rentalpb.Trip{
			AccountId: r.accountID.String(),
			Status:    r.status,
		})

		if err != nil {
			t.Fatalf("cannot create rows: %v", err)
		}
	}

	// 测试样例
	cases := []struct {
		name       string
		accountID  string
		status     rentalpb.TripStatus
		wantCount  int
		wantOnlyID string
	}{
		{
			name:      "get_all",
			accountID: "account_id_for_get_trips",
			status:    rentalpb.TripStatus_TS_NOT_SPECIFIED,
			wantCount: 4,
		},
		{
			name:       "get_in_progress",
			accountID:  "account_id_for_get_trips",
			status:     rentalpb.TripStatus_IN_PROGRESS,
			wantCount:  1,
			wantOnlyID: "61bb48d44e17e2f2040c1b56",
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			res, err := m.GetTrips(context.Background(),
				shared_id.AccountID(testCase.accountID),
				testCase.status)

			if err != nil {
				t.Errorf("cannot get trips: %v", err)
			}

			if testCase.wantCount != len(res) {
				t.Errorf("incorrect result count; want: %d, got: %d", testCase.wantCount, len(res))
			}

			if testCase.wantOnlyID != "" && len(res) > 0 {
				if testCase.wantOnlyID != res[0].ID.Hex() {
					t.Errorf("incorrect only id; want: %q, got: %q", testCase.wantOnlyID, res[0].ID.Hex())
				}
			}
		})
	}
}

func TestUpdateTrip(t *testing.T) {
	c := context.Background()
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot connect mongodb: %v", err)
	}
	m := NewMongo(mc.Database("coolcar"))

	// insert a trip before
	tripID := shared_id.TripID("61bb48e569239c93c16d96bd")
	accountID := shared_id.AccountID("account_for_update")

	// set time
	var now int64 = 10000
	shared_mongo.NewObjIDWithValue(tripID)
	shared_mongo.UpdateAt = func() int64 {
		return now
	}

	tr, err := m.CreateTrip(c, &rentalpb.Trip{
		AccountId: accountID.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PositionName: "start_poi",
		},
	})
	if err != nil {
		t.Fatalf("cannot create a trip: %v", err)
	}

	if tr.UpdateAt != 10000 {
		t.Fatalf("wrong updateAt; want: %d, got: %d", now, tr.UpdateAt)
	}

	// test cases
	update := &rentalpb.Trip{
		AccountId: accountID.String(),
		Status:    rentalpb.TripStatus_IN_PROGRESS,
		Start: &rentalpb.LocationStatus{
			PositionName: "start_poi_updated",
		},
	}

	cases := []struct {
		name         string
		now          int64
		withUpdateAt int64
		wantErr      bool
	}{
		{
			name:         "normal_update",
			now:          20000,
			withUpdateAt: 10000,
			wantErr:      false,
		},
		{
			name:         "update_with_stale_timestamp",
			now:          3000,
			withUpdateAt: 10000,
			wantErr:      true,
		},
		{
			name:         "update_with_refetch",
			now:          40000,
			withUpdateAt: 20000,
			wantErr:      false,
		},
	}

	// update
	for _, testCase := range cases {
		now = testCase.now
		err = m.UpdateTrip(c, tripID, accountID, testCase.withUpdateAt, update)
		if testCase.wantErr {
			if err == nil {
				t.Errorf("%s: want error, got none", testCase.name)
			} else {
				continue
			}
		} else {
			if err != nil {
				t.Errorf("%s: cannot update: %v", testCase.name, err)
			}
		}

		updatedTrip, err := m.GetTrip(c, tripID, accountID)
		if err != nil {
			if err != nil {
				t.Errorf("%s: cannot get trip after update: %v", testCase.name, err)
			}
		}

		if now != updatedTrip.UpdateAt {
			t.Errorf("%s: incorrect updateat: want: %d, got: %d", testCase.name, testCase.now, updatedTrip.UpdateAt)
		}
	}
}

// 必须这么命名
func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
