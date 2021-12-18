package trip

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip/dao"
	shared_auth "coolcar/shared/auth"
	shared_id "coolcar/shared/id"
)

// Service implements a trip service
type Service struct {
	Mongo           *dao.Mongo
	Logger          *zap.Logger
	ProfileManager  ProfileManager
	CarManager      CarManager
	PositionManager PositionManager
}

// ProfileManager defines the ACL (Anti Corruption Layer)
// for profile verification logic.
type ProfileManager interface {
	Verify(context.Context, shared_id.AccountID) (shared_id.IdentityID, error)
}

// CarManager defines the ACL for car management
type CarManager interface {
	Verify(context.Context, shared_id.CarID, *rentalpb.Location) error
	Unlock(context.Context, shared_id.CarID) error
}

// PositionManager resolves POI (point of position)
type PositionManager interface {
	Resolve(context.Context, *rentalpb.Location) (string, error)
}

// CreateTrip creates a trip
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	// verify driver
	accountID, err := shared_auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	identityID, err := s.ProfileManager.Verify(c, accountID)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	// check car status
	poi, err := s.PositionManager.Resolve(c, req.Start)
	if err != nil {
		s.Logger.Info("cannot resolve poi", zap.Stringer("location", req.Start), zap.Error(err))
	}
	carID := shared_id.CarID(req.CarId)
	err = s.CarManager.Verify(c, carID, req.Start)
	if err != nil {
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	// create trip: write db and bill
	ls := &rentalpb.LocationStatus{
		Location:     req.Start,
		PositionName: poi,
	}
	tripRecord, err := s.Mongo.CreateTrip(c, &rentalpb.Trip{
		AccountId:  accountID.String(),
		CarId:      carID.String(),
		IdentityId: identityID.String(),
		Status:     rentalpb.TripStatus_IN_PROGRESS,
		Start:      ls,
		Current:    ls,
	})
	if err != nil {
		s.Logger.Warn("cannot create trip", zap.Error(err))
		return nil, status.Error(codes.AlreadyExists, "")
	}
	// unlock car
	go func() {
		err = s.CarManager.Unlock(context.Background(), carID)
		if err != nil {
			s.Logger.Error("unlock car failed", zap.Error(err))
		}
	}()
	return &rentalpb.TripEntity{
		Id:   tripRecord.ID.Hex(),
		Trip: tripRecord.Trip,
	}, nil
}

// GetTrip gets a trip
func (s *Service) GetTrip(c context.Context, req *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// GetTrips gets trips
func (s *Service) GetTrips(c context.Context, req *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "")
}

// UpdateTrip updates a trip
func (s *Service) UpdateTrip(c context.Context, req *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	// get accountID from context
	accountID, err := shared_auth.AccountIDFromContext(c)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "")
	}
	// get trip id from request
	tripID := shared_id.TripID(req.Id)
	// get trip by trip id and account id
	tripRecord, err := s.Mongo.GetTrip(c, tripID, accountID)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "get trip error: %v", err)
	}
	// update current status
	if req.Current != nil {
		tripRecord.Trip.Current = calculateCurrentStatus(tripRecord.Trip, req.Current)
	}
	// if trip finished
	if req.EndTrip {
		tripRecord.Trip.End = tripRecord.Trip.Current
		tripRecord.Trip.Status = rentalpb.TripStatus_FINISHED
	}
	// update trip in mongodb
	err = s.Mongo.UpdateTrip(c, tripID, accountID, tripRecord.UpdateAt, tripRecord.Trip)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "update trip error: %v", err)
	}
	return tripRecord.Trip, nil
}

// calculateCurrentStatus
func calculateCurrentStatus(trip *rentalpb.Trip, cur *rentalpb.Location) *rentalpb.LocationStatus {
	return nil
}
