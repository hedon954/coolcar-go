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
	Mongo  *dao.Mongo
	Logger *zap.Logger
}

// CreateTrip creates a trip
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	return nil, status.Error(codes.Unimplemented, "")
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
