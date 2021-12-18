package trip

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/client/poi"
	"coolcar/rental/trip/dao"
	shared_auth "coolcar/shared/auth"
	shared_id "coolcar/shared/id"
	shared_mongo "coolcar/shared/mongo"
	mongotesting "coolcar/shared/testing"

	"go.uber.org/zap"
)

func TestCreateTrip(t *testing.T) {
	c := shared_auth.ContextWithAccountID(
		context.Background(),
		shared_id.AccountID("account1"),
	)
	mc, err := mongotesting.NewClient(c)
	if err != nil {
		t.Fatalf("cannot create mongo client: %v", err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		t.Fatalf("cannot create logger: %v", err)
	}

	pm := &profileManager{}
	cm := &carManager{}

	s := Service{
		ProfileManager:  pm,
		CarManager:      cm,
		PositionManager: &poi.Manager{},
		Mongo:           dao.NewMongo(mc.Database("coolcar")),
		Logger:          logger,
	}

	req := &rentalpb.CreateTripRequest{
		CarId: "car1",
		Start: &rentalpb.Location{
			Latitude:  32.123,
			Longitude: 114.2525,
		},
	}

	pm.iID = "identity1"

	golden := `{"account_id":"account1","car_id":"car1","start":{"Location":{"latitude":32.123,"longitude":114.2525},"position_name":"天安门"},"current":{"Location":{"latitude":32.123,"longitude":114.2525},"position_name":"天安门"},"status":1,"identity_id":"identity1"}`

	cases := []struct {
		name         string
		tripID       string
		profileErr   error
		carVerifyErr error
		carUnlockErr error
		want         string
		wantErr      bool
	}{
		{
			name:    "normal_create",
			tripID:  "61bb48d44e17e2f2040c1b59",
			want:    golden,
			wantErr: false,
		},
		{
			name:       "profile_err",
			tripID:     "61bb468dc648e913cded5652",
			profileErr: fmt.Errorf("profile"),
			wantErr:    true,
		},
		{
			name:         "car_verify_err",
			tripID:       "61bb48e569239c93c16c96bd",
			carVerifyErr: fmt.Errorf("verify"),
			wantErr:      true,
		},
		{
			name:         "car_unlock_err",
			tripID:       "61bb490936560e9b4c1dd2a3",
			carUnlockErr: fmt.Errorf("unlock"),
			wantErr:      false,
			want:         golden,
		},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			shared_mongo.NewObjIDWithValue(shared_id.TripID(testCase.tripID))
			pm.err = testCase.profileErr
			cm.unlockErr = testCase.carUnlockErr
			cm.verifyErr = testCase.carVerifyErr

			res, err := s.CreateTrip(c, req)
			if testCase.wantErr {
				if err == nil {
					t.Errorf("want error; got none")
				} else {
					return
				}
			}
			if err != nil {
				t.Errorf("error creating trip: %v", err)
				return
			}
			if res.Id != testCase.tripID {
				t.Errorf("incorrect id; want %q, got %q", testCase.tripID, res.Id)
			}
			b, err := json.Marshal(res.Trip)
			if err != nil {
				t.Errorf("cannot marshal response: %v", err)
			}
			tripStr := string(b)
			if testCase.want != tripStr {
				t.Errorf("incorrect response; want %s, got %s", testCase.want, tripStr)
			}
		})
	}
}

type profileManager struct {
	iID shared_id.IdentityID
	err error
}

func (p *profileManager) Verify(c context.Context, aid shared_id.AccountID) (shared_id.IdentityID, error) {
	return p.iID, p.err
}

type carManager struct {
	verifyErr error
	unlockErr error
}

func (c *carManager) Verify(ctx context.Context, cid shared_id.CarID, l *rentalpb.Location) error {
	return c.verifyErr
}

func (c *carManager) Unlock(ctx context.Context, cid shared_id.CarID) error {
	return c.unlockErr
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m))
}
