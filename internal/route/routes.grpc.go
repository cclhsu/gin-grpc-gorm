package route

import (
	"context"

	"github.com/cclhsu/gin-grpc-gorm/internal/grpc_service_server"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	authGrpcService "github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/auth"
	healthGrpcService "github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/health"
	helloGrpcService "github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/hello"
	teamGrpcService "github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/team"
	userGrpcService "github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/user"
)

func SetupGrpcRoutes(ctx context.Context, logger *logrus.Logger, grpcServer *grpc.Server, authService *service.AuthService, userService *service.UserService, teamService *service.TeamService, helloService *service.HelloService, healthService *service.HealthService) {

	// Register Auth service server
	authGrpcServiceServer := grpc_service_server.NewAuthServiceServer(ctx, logger, authService)
	authGrpcService.RegisterAuthServiceServer(grpcServer, authGrpcServiceServer)

	// Register User service server
	userGrpcServiceServer := grpc_service_server.NewUserServiceServer(ctx, logger, userService)
	userGrpcService.RegisterUserServiceServer(grpcServer, userGrpcServiceServer)

	// Register team service server
	teamGrpcServiceServer := grpc_service_server.NewTeamServiceServer(ctx, logger, teamService)
	teamGrpcService.RegisterTeamServiceServer(grpcServer, teamGrpcServiceServer)

	// Register Hello service server
	helloGrpcServiceServer := grpc_service_server.NewHelloServiceServer(ctx, logger, helloService)
	helloGrpcService.RegisterHelloServiceServer(grpcServer, helloGrpcServiceServer)

	// Register Health service server
	healthGrpcServiceServer := grpc_service_server.NewHealthServiceServer(ctx, logger, healthService)
	healthGrpcService.RegisterHealthServiceServer(grpcServer, healthGrpcServiceServer)
}
