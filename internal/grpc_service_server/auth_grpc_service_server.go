package grpc_service_server

import (
	// "context"

	"context"

	"github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/auth"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/sirupsen/logrus"
	// "github.com/cclhsu/gin-grpc-gorm/internal/utils"
)

type AuthServiceServer struct {
	ctx    context.Context
	logger *logrus.Logger
	auth.UnimplementedAuthServiceServer
	authService *service.AuthService
}

func NewAuthServiceServer(ctx context.Context, logger *logrus.Logger, hs *service.AuthService) *AuthServiceServer {
	return &AuthServiceServer{
		ctx:         ctx,
		logger:      logger,
		authService: hs,
	}
}

// // grpcurl -plaintext -d '{"username": "john.doe", "password": "changeme"}' 0.0.0.0:50051 auth.AuthService/Login
// func (ass *AuthServiceServer) Login(ctx context.Context, in *auth.LoginRequest) (*auth.LoginResponse, error) {
//	token, err := ass.authService.Login(in.Username, in.Password)
//	if err != nil {
//		return nil, err
//	}
//	loginResponse := &auth.LoginResponse{
//		Token: token,
//	}
//	return loginResponse, nil
// }

// // grpcurl -plaintext -H "authorization: Bearer <JWT_TOKEN>" -plaintext 0.0.0.0:50051 auth.AuthService/GetProtectedData
// func (ass *AuthServiceServer) GetProtectedData(ctx context.Context, in *auth.SecuredRequest) (*auth.SecuredResponse, error) {
//	message := "This is a protected data"
//	securedResponse := &auth.SecuredResponse{
//		Success:   true,
//		Message:   message,
//		ErrorCode: 200,
//	}
//	return securedResponse, nil
// }

// // grpcurl -plaintext -d '{"token": "<YOUR_TOKEN>"}' 0.0.0.0:50051 auth.AuthService/GetProfile
// func (ass *AuthServiceServer) GetProfile(ctx context.Context, in *auth.SecuredRequest) (*auth.GetUserProfileResponse, error) {
//	UUID, err := utils.GetUUIDFromToken(in.Token)
//	if err != nil {
//		return nil, err
//	}
//	userProfile, err := ass.authService.GetUserProfile(UUID)
//	if err != nil {
//		return nil, err
//	}

//	dob := userProfile.DateOfBirth.String()
//	phoneNumber := userProfile.PhoneNumber
//	userProfileResponse := &auth.GetUserProfileResponse{
//		UUID:		 userProfile.UUID,
//		LastName:	 userProfile.LastName,
//		FirstName:	 userProfile.FirstName,
//		Role:		 userProfile.Role,
//		DateOfBirth: &dob,
//		Email:		 userProfile.Email,
//		PhoneNumber: &phoneNumber,
//		Username:	 userProfile.Username,
//	}
//	return userProfileResponse, nil

// }

// // grpcurl -plaintext -d '{"token":	 "<YOUR_TOKEN>"}' 0.0.0.0:50051 auth.AuthService/Logout
// func (ass *AuthServiceServer) Logout(ctx context.Context, in *auth.SecuredRequest) (*auth.LogoutResponse, error) {

//	// err := ass.authService.Logout(UUID)
//	// if err == nil {
//	//	return nil, err
//	// }
//	logoutResponse := &auth.LogoutResponse{
//		Message: "Successfully logged out",
//	}
//	return logoutResponse, nil
// }
