package trip

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	rentalpb "coolcar/rental/api/gen/v1"
)

// Service implements a trip service
type Service struct {
	Logger *zap.Logger
}

// CreateTrip creates a trip
func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.CreateTripResponse, error) {
	s.Logger.Info("create trip", zap.String("start", req.Start))
	return nil, status.Error(codes.Unimplemented, "")
}
