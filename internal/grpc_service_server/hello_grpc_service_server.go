package grpc_service_server

import (
	"context"
	"log"

	"github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/hello"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type HelloServiceServer struct {
	ctx    context.Context
	logger *logrus.Logger
	hello.UnimplementedHelloServiceServer
	helloService *service.HelloService
}

func NewHelloServiceServer(ctx context.Context, logger *logrus.Logger, hs *service.HelloService) *HelloServiceServer {
	return &HelloServiceServer{
		ctx:          ctx,
		logger:       logger,
		helloService: hs,
	}
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 hello.HelloService/GetHelloString
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/hello.proto -import-path api/grpc/proto 0.0.0.0:50051 hello.HelloService/GetHelloString
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/hello.proto -import-path api/grpc/proto localhost:50051 hello.HelloService/GetHelloString
func (hss *HelloServiceServer) GetHelloString(ctx context.Context, in *emptypb.Empty) (*hello.HelloStringResponse, error) {
	// response, err := hss.helloService.GetHelloString()
	// if err != nil {
	//	log.Printf("Failed to get hello string: %v", err)
	//	return nil, err
	// }

	return &hello.HelloStringResponse{}, nil
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 hello.HelloService/GetHelloJson
// grpcurl -plaintext -d '{}'  -proto api/grpc/proto/hello.proto -import-path api/grpc/proto 0.0.0.0:50051 hello.HelloService/GetHelloJson
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/hello.proto -import-path api/grpc/proto localhost:50051 hello.HelloService/GetHelloJson
func (hss *HelloServiceServer) GetHelloJson(ctx context.Context, in *emptypb.Empty) (*hello.HelloJsonResponse, error) {
	response, err := hss.helloService.GetHelloJson()
	if err != nil {
		log.Printf("Failed to get hello JSON: %v", err)
		return nil, err
	}

	helloJsonResponse := &hello.HelloJsonResponse{
		Data: &hello.Data{
			Message: response.Data.Message,
		},
	}
	return helloJsonResponse, nil
}
