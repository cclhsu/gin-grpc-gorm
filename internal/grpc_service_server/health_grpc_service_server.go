package grpc_service_server

import (
	"context"
	"log"

	"github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/health"
	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/sirupsen/logrus"
)

type HealthServiceServer struct {
	ctx    context.Context
	logger *logrus.Logger
	health.UnimplementedHealthServiceServer
	healthService *service.HealthService
}

func NewHealthServiceServer(ctx context.Context, logger *logrus.Logger, hcs *service.HealthService) *HealthServiceServer {
	return &HealthServiceServer{
		ctx:           ctx,
		logger:        logger,
		healthService: hcs,
	}
}

// grpcurl -plaintext -d '{"service": "Health"}' 0.0.0.0:50051 health.HealthService/IsHealthy | jq
// grpcurl -plaintext -d '{"service": "Health"}' -proto api/grpc/proto/health.proto -import-path api/grpc/proto 0.0.0.0:50051 health.HealthService/IsHealthy | jq
// grpcurl -plaintext -d '{"service": "Health"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/health.proto -import-path api/grpc/proto 0.0.0.0:50051 health.HealthService/IsHealthy | jq
func (rhcss *HealthServiceServer) IsHealthy(ctx context.Context, request *health.HealthRequest) (*health.HealthResponse, error) {
	// Convert the gRPC request to the internal model
	healthRequest := model.HealthRequest{
		Service: request.Service,
	}

	// Call the internal health service
	response, err := rhcss.healthService.IsHealthy(healthRequest)
	if err != nil {
		log.Printf("Failed to get health: %v", err)
		return nil, err
	}

	// Convert the internal model response to the gRPC response
	healthResponse := &health.HealthResponse{
		Status: health.ServingStatus(response.Status),
	}

	return healthResponse, nil
}

// grpcurl -plaintext -d '{"service": "Health"}' 0.0.0.0:50051 health.HealthService/IsALive | jq
// grpcurl -plaintext -d '{"service": "Health"}' -proto api/grpc/proto/health.proto -import-path api/grpc/proto 0.0.0.0:50051 health.HealthService/IsALive | jq
// grpcurl -plaintext -d '{"service": "Health"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/health.proto -import-path api/grpc/proto 0.0.0.0:50051 health.HealthService/IsALive | jq
func (rhcss *HealthServiceServer) IsALive(ctx context.Context, request *health.HealthRequest) (*health.HealthResponse, error) {
	// Convert the gRPC request to the internal model
	healthRequest := model.HealthRequest{
		Service: request.Service,
	}

	// Call the internal health service
	response, err := rhcss.healthService.IsALive(healthRequest)
	if err != nil {
		log.Printf("Failed to get health: %v", err)
		return nil, err
	}

	// Convert the internal model response to the gRPC response
	healthResponse := &health.HealthResponse{
		Status: health.ServingStatus(response.Status),
	}

	return healthResponse, nil
}

// grpcurl -plaintext -d '{"service": "Health"}' 0.0.0.0:50051 health.HealthService/IsReady | jq
// grpcurl -plaintext -d '{"service": "Health"}' -proto api/grpc/proto/health.proto -import-path api/grpc/proto 0.0.0.0:50051 health.HealthService/IsReady | jq
// grpcurl -plaintext -d '{"service": "Health"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/health.proto -import-path api/grpc/proto 0.0.0.0:50051 health.HealthService/IsReady | jq
func (rhcss *HealthServiceServer) IsReady(ctx context.Context, request *health.HealthRequest) (*health.HealthResponse, error) {
	// Convert the gRPC request to the internal model
	healthRequest := model.HealthRequest{
		Service: request.Service,
	}

	// Call the internal health service
	response, err := rhcss.healthService.IsReady(healthRequest)
	if err != nil {
		log.Printf("Failed to get health: %v", err)
		return nil, err
	}

	// Convert the internal model response to the gRPC response
	healthResponse := &health.HealthResponse{
		Status: health.ServingStatus(response.Status),
	}

	return healthResponse, nil
}
