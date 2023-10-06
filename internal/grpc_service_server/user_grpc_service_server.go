package grpc_service_server

import (
	"context"

	"github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/common"
	"github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/user"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	// "github.com/cclhsu/gin-grpc-gorm/internal/model"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
)

// type UserServiceClient interface {
// 	ListUserIdsAndUUIDs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUserIdUuid, error)
// 	ListUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersResponse, error)
// 	ListUsersMetadata(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersMetadataResponse, error)
// 	ListUsersContent(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListUsersContentResponse, error)
// 	GetUser(ctx context.Context, in *GetUserByUuidRequest, opts ...grpc.CallOption) (*User, error)
// 	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error)
// 	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error)
// 	DeleteUser(ctx context.Context, in *GetUserByUuidRequest, opts ...grpc.CallOption) (*User, error)
// 	GetUserById(ctx context.Context, in *GetUserByIdRequest, opts ...grpc.CallOption) (*User, error)
// 	GetUserByName(ctx context.Context, in *GetUserByUsernameRequest, opts ...grpc.CallOption) (*User, error)
// 	GetUserByEmail(ctx context.Context, in *GetUserByEmailRequest, opts ...grpc.CallOption) (*User, error)
// 	UpdateUserMetadata(ctx context.Context, in *UpdateUserMetadataRequest, opts ...grpc.CallOption) (*UserMetadataResponse, error)
// 	UpdateUserContent(ctx context.Context, in *UpdateUserContentRequest, opts ...grpc.CallOption) (*UserContentResponse, error)
// 	GetUserMetadata(ctx context.Context, in *GetUserByUuidRequest, opts ...grpc.CallOption) (*UserMetadataResponse, error)
// 	GetUserContent(ctx context.Context, in *GetUserByUuidRequest, opts ...grpc.CallOption) (*UserContentResponse, error)
// }

type UserServiceServer struct {
	ctx    context.Context
	logger *logrus.Logger
	user.UnimplementedUserServiceServer
	userService *service.UserService
}

func NewUserServiceServer(ctx context.Context, logger *logrus.Logger, hs *service.UserService) *UserServiceServer {
	return &UserServiceServer{
		ctx:         ctx,
		logger:      logger,
		userService: hs,
	}
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 user.UserService/ListUserIdsAndUUIDs
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/ListUserIdsAndUUIDs
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/ListUserIdsAndUUIDs
func (uss *UserServiceServer) ListUserIdsAndUUIDs(ctx context.Context, in *emptypb.Empty) (*user.ListUserIdUuid, error) {
	idUuids, err := uss.userService.ListUserIdsAndUUIDs()
	if err != nil {
		return nil, err
	}

	// convert []*model.IdUuid to user.ListUserIdUuid
	listUserIdUuid := &user.ListUserIdUuid{}
	for _, idUuid := range idUuids {
		listUserIdUuid.UserIdUuids = append(listUserIdUuid.UserIdUuids, &common.IdUuid{
			ID:   idUuid.ID,
			UUID: idUuid.UUID,
		})
	}

	uss.logger.Info(listUserIdUuid)
	return listUserIdUuid, nil
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 user.UserService/ListUsers
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/ListUsers
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/ListUsers
func (uss *UserServiceServer) ListUsers(ctx context.Context, in *emptypb.Empty) (*user.ListUsersResponse, error) {
	users, err := uss.userService.ListUsers()
	if err != nil {
		return nil, err
	}

	// convert []*model.UserResponse to user.ListUsersResponse
	listUsersResponse := &user.ListUsersResponse{}
	for _, userModel := range users {

		// Assuming you have an empty slice of the target type
		var targetProjectRoles []common.PROJECT_ROLE_TYPES

		// Iterate through userModel.Content.ProjectRoles and convert each element
		for _, role := range userModel.Content.ProjectRoles {
			// Perform the conversion and append to the target slice
			targetRole := common.PROJECT_ROLE_TYPES(role) // Assuming there's a valid conversion
			targetProjectRoles = append(targetProjectRoles, targetRole)
		}

		// Assuming you have an empty slice of the target type
		var targetScrumRoles []common.SCRUM_ROLE_TYPES

		// Iterate through userModel.Content.ScrumRoles and convert each element
		for _, role := range userModel.Content.ScrumRoles {
			// Perform the conversion and append to the target slice
			targetRole := common.SCRUM_ROLE_TYPES(role) // Assuming there's a valid conversion
			targetScrumRoles = append(targetScrumRoles, targetRole)
		}

		listUsersResponse.Users = append(listUsersResponse.Users, &user.User{
			ID:   userModel.ID,
			UUID: userModel.UUID,
			Metadata: &user.UserMetadata{
				Name: userModel.Metadata.Name,
				Dates: &common.CommonDate{
					CreatedAt:   userModel.Metadata.Dates.CreatedAt,
					CreatedBy:   userModel.Metadata.Dates.CreatedBy,
					UpdatedAt:   userModel.Metadata.Dates.UpdatedAt,
					UpdatedBy:   userModel.Metadata.Dates.UpdatedBy,
					StartDate:   &userModel.Metadata.Dates.StartDate,
					EndDate:     &userModel.Metadata.Dates.EndDate,
					StartedAt:   &userModel.Metadata.Dates.StartedAt,
					StartedBy:   &userModel.Metadata.Dates.StartedBy,
					CompletedAt: &userModel.Metadata.Dates.CompletedAt,
					CompletedBy: &userModel.Metadata.Dates.CompletedBy,
				},
			},
			Content: &user.UserContent{
				Email:        userModel.Content.Email,
				Phone:        userModel.Content.Phone,
				LastName:     userModel.Content.LastName,
				FirstName:    userModel.Content.FirstName,
				ProjectRoles: targetProjectRoles,
				ScrumRoles:   targetScrumRoles,
				Password:     userModel.Content.Password,
			},
		})
	}

	uss.logger.Info(listUsersResponse)
	return listUsersResponse, nil
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 user.UserService/ListUsersMetadata
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/ListUsersMetadata
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/ListUsersMetadata
func (uss *UserServiceServer) ListUsersMetadata(ctx context.Context, in *emptypb.Empty) (*user.ListUsersMetadataResponse, error) {
	listUsersMetadata, err := uss.userService.ListUsersMetadata()
	if err != nil {
		return nil, err
	}

	// convert []*model.UserMetadata to user.ListUsersMetadataResponse
	listUsersMetadataResponse := &user.ListUsersMetadataResponse{}
	for _, userModel := range listUsersMetadata.UserMetadataResponses {
		listUsersMetadataResponse.UserMetadataResponses = append(listUsersMetadataResponse.UserMetadataResponses, &user.UserMetadataResponse{
			ID:   userModel.ID,
			UUID: userModel.UUID,
			Metadata: &user.UserMetadata{
				Name: userModel.Metadata.Name,
				Dates: &common.CommonDate{
					CreatedAt:   userModel.Metadata.Dates.CreatedAt,
					CreatedBy:   userModel.Metadata.Dates.CreatedBy,
					UpdatedAt:   userModel.Metadata.Dates.UpdatedAt,
					UpdatedBy:   userModel.Metadata.Dates.UpdatedBy,
					StartDate:   &userModel.Metadata.Dates.StartDate,
					EndDate:     &userModel.Metadata.Dates.EndDate,
					StartedAt:   &userModel.Metadata.Dates.StartedAt,
					StartedBy:   &userModel.Metadata.Dates.StartedBy,
					CompletedAt: &userModel.Metadata.Dates.CompletedAt,
					CompletedBy: &userModel.Metadata.Dates.CompletedBy,
				},
			},
		})
	}

	uss.logger.Info(listUsersMetadataResponse)
	return listUsersMetadataResponse, nil
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 user.UserService/ListUsersContent
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/ListUsersContent
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/ListUsersContent
func (uss *UserServiceServer) ListUsersContent(ctx context.Context, in *emptypb.Empty) (*user.ListUsersContentResponse, error) {
	listUsersContent, err := uss.userService.ListUsersContent()
	if err != nil {
		return nil, err
	}

	// convert []*model.UserContent to user.ListUsersContentResponse
	listUsersContentResponse := &user.ListUsersContentResponse{}
	for _, userModel := range listUsersContent.UserContentResponses {

		// Assuming you have an empty slice of the target type
		var targetProjectRoles []common.PROJECT_ROLE_TYPES

		// Iterate through userModel.Content.ProjectRoles and convert each element
		for _, role := range userModel.Content.ProjectRoles {
			// Perform the conversion and append to the target slice
			targetRole := common.PROJECT_ROLE_TYPES(role) // Assuming there's a valid conversion
			targetProjectRoles = append(targetProjectRoles, targetRole)
		}

		// Assuming you have an empty slice of the target type
		var targetScrumRoles []common.SCRUM_ROLE_TYPES

		// Iterate through userModel.Content.ScrumRoles and convert each element
		for _, role := range userModel.Content.ScrumRoles {
			// Perform the conversion and append to the target slice
			targetRole := common.SCRUM_ROLE_TYPES(role) // Assuming there's a valid conversion
			targetScrumRoles = append(targetScrumRoles, targetRole)
		}

		listUsersContentResponse.UserContentResponses = append(listUsersContentResponse.UserContentResponses, &user.UserContentResponse{
			ID:   userModel.ID,
			UUID: userModel.UUID,
			Content: &user.UserContent{
				Email:        userModel.Content.Email,
				Phone:        userModel.Content.Phone,
				LastName:     userModel.Content.LastName,
				FirstName:    userModel.Content.FirstName,
				ProjectRoles: targetProjectRoles,
				ScrumRoles:   targetScrumRoles,
				Password:     userModel.Content.Password,
			},
		})
	}

	uss.logger.Info(listUsersContentResponse)
	return listUsersContentResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 user.UserService/GetUser
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/GetUser
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/GetUser
func (uss *UserServiceServer) GetUser(ctx context.Context, in *user.GetUserByUuidRequest) (*user.User, error) {
	userModel, err := uss.userService.GetUser(in.UUID)
	if err != nil {
		return nil, err
	}

	return uss.convertUser(userModel)
}

// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' 0.0.0.0:50051 user.UserService/CreateUser
// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/CreateUser
func (uss *UserServiceServer) CreateUser(ctx context.Context, in *user.CreateUserRequest) (*user.User, error) {
	var projectRoles []model.PROJECT_ROLE_TYPES
	for _, role := range in.Content.ProjectRoles {
		projectRoles = append(projectRoles, model.PROJECT_ROLE_TYPES(role))
	}

	var scrumRoles []model.SCRUM_ROLE_TYPES
	for _, role := range in.Content.ScrumRoles {
		scrumRoles = append(scrumRoles, model.SCRUM_ROLE_TYPES(role))
	}

	createUserRequest := model.CreateUserRequest{
		ID:   in.ID,
		UUID: in.UUID,
		Metadata: &model.UserMetadata{
			Name: in.Metadata.Name,
			Dates: &model.CommonDate{
				CreatedAt:   in.Metadata.Dates.CreatedAt,
				CreatedBy:   in.Metadata.Dates.CreatedBy,
				UpdatedAt:   in.Metadata.Dates.UpdatedAt,
				UpdatedBy:   in.Metadata.Dates.UpdatedBy,
				StartDate:   *in.Metadata.Dates.StartDate,
				EndDate:     *in.Metadata.Dates.EndDate,
				StartedAt:   *in.Metadata.Dates.StartedAt,
				StartedBy:   *in.Metadata.Dates.StartedBy,
				CompletedAt: *in.Metadata.Dates.CompletedAt,
				CompletedBy: *in.Metadata.Dates.CompletedBy,
			},
		},
		Content: &model.UserContent{
			Email:        in.Content.Email,
			Phone:        in.Content.Phone,
			LastName:     in.Content.LastName,
			FirstName:    in.Content.FirstName,
			ProjectRoles: projectRoles,
			ScrumRoles:   scrumRoles,
			Password:     in.Content.Password,
		},
	}

	userModel, err := uss.userService.CreateUser(createUserRequest)
	if err != nil {
		return nil, err
	}

	return uss.convertUser(userModel)
}

// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' 0.0.0.0:50051 user.UserService/UpdateUser
// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/UpdateUser
func (uss *UserServiceServer) UpdateUser(ctx context.Context, in *user.UpdateUserRequest) (*user.User, error) {
	var projectRoles []model.PROJECT_ROLE_TYPES
	for _, role := range in.Content.ProjectRoles {
		projectRoles = append(projectRoles, model.PROJECT_ROLE_TYPES(role))
	}

	var scrumRoles []model.SCRUM_ROLE_TYPES
	for _, role := range in.Content.ScrumRoles {
		scrumRoles = append(scrumRoles, model.SCRUM_ROLE_TYPES(role))
	}

	updateUserRequest := model.UpdateUserRequest{
		// ID:   in.ID,
		UUID: in.UUID,
		Metadata: &model.UserMetadata{
			Name: in.Metadata.Name,
			Dates: &model.CommonDate{
				CreatedAt:   in.Metadata.Dates.CreatedAt,
				CreatedBy:   in.Metadata.Dates.CreatedBy,
				UpdatedAt:   in.Metadata.Dates.UpdatedAt,
				UpdatedBy:   in.Metadata.Dates.UpdatedBy,
				StartDate:   *in.Metadata.Dates.StartDate,
				EndDate:     *in.Metadata.Dates.EndDate,
				StartedAt:   *in.Metadata.Dates.StartedAt,
				StartedBy:   *in.Metadata.Dates.StartedBy,
				CompletedAt: *in.Metadata.Dates.CompletedAt,
				CompletedBy: *in.Metadata.Dates.CompletedBy,
			},
		},
		Content: &model.UserContent{
			Email:        in.Content.Email,
			Phone:        in.Content.Phone,
			LastName:     in.Content.LastName,
			FirstName:    in.Content.FirstName,
			ProjectRoles: projectRoles,
			ScrumRoles:   scrumRoles,
			Password:     in.Content.Password,
		},
	}

	userModel, err := uss.userService.UpdateUser(in.UUID, updateUserRequest)
	if err != nil {
		return nil, err
	}

	return uss.convertUser(userModel)
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 user.UserService/DeleteUser
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/DeleteUser
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/DeleteUser
func (uss *UserServiceServer) DeleteUser(ctx context.Context, in *user.GetUserByUuidRequest) (*user.User, error) {
	userModel, err := uss.userService.DeleteUser(in.UUID)
	if err != nil {
		return nil, err
	}

	return uss.convertUser(userModel)
}

// grpcurl -plaintext -d '{"ID": "john.doe"}' 0.0.0.0:50051 user.UserService/GetUserById
// grpcurl -plaintext -d '{"ID": "john.doe"}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/GetUserById
// grpcurl -plaintext -d '{"ID": "john.doe"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/GetUserById
func (uss *UserServiceServer) GetUserById(ctx context.Context, in *user.GetUserByIdRequest) (*user.User, error) {
	userModel, err := uss.userService.GetUserByID(in.ID)
	if err != nil {
		return nil, err
	}

	return uss.convertUser(userModel)
}

// grpcurl -plaintext -d '{"name": "John Doe"}' 0.0.0.0:50051 user.UserService/GetUserByName
// grpcurl -plaintext -d '{"name": "John Doe"}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/GetUserByName
// grpcurl -plaintext -d '{"name": "John Doe"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/GetUserByName
func (uss *UserServiceServer) GetUserByName(ctx context.Context, in *user.GetUserByUsernameRequest) (*user.User, error) {
	userModel, err := uss.userService.GetUserByName(in.Username)
	if err != nil {
		return nil, err
	}

	return uss.convertUser(userModel)
}

// grpcurl -plaintext -d '{"email": "john.doe@mail.com"}' 0.0.0.0:50051 user.UserService/GetUserByEmail
// grpcurl -plaintext -d '{"email": "john.doe@mail.com"}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/GetUserByEmail
// grpcurl -plaintext -d '{"email": "john.doe@mail.com"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/GetUserByEmail
func (uss *UserServiceServer) GetUserByEmail(ctx context.Context, in *user.GetUserByEmailRequest) (*user.User, error) {
	userModel, err := uss.userService.GetUserByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	return uss.convertUser(userModel)
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}}' 0.0.0.0:50051 user.UserService/UpdateUserMetadata
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/UpdateUserMetadata
func (uss *UserServiceServer) UpdateUserMetadata(ctx context.Context, in *user.UpdateUserMetadataRequest) (*user.UserMetadataResponse, error) {
	updateUserMetadataRequest := model.UpdateUserMetadataRequest{
		// ID:   in.ID,
		UUID: in.UUID,
		Metadata: &model.UserMetadata{
			Name: in.Metadata.Name,
			Dates: &model.CommonDate{
				CreatedAt:   in.Metadata.Dates.CreatedAt,
				CreatedBy:   in.Metadata.Dates.CreatedBy,
				UpdatedAt:   in.Metadata.Dates.UpdatedAt,
				UpdatedBy:   in.Metadata.Dates.UpdatedBy,
				StartDate:   *in.Metadata.Dates.StartDate,
				EndDate:     *in.Metadata.Dates.EndDate,
				StartedAt:   *in.Metadata.Dates.StartedAt,
				StartedBy:   *in.Metadata.Dates.StartedBy,
				CompletedAt: *in.Metadata.Dates.CompletedAt,
				CompletedBy: *in.Metadata.Dates.CompletedBy,
			},
		},
	}

	userMetadata, err := uss.userService.UpdateUserMetadata(in.UUID, updateUserMetadataRequest)
	if err != nil {
		return nil, err
	}

	userMetadataResponse := &user.UserMetadataResponse{
		Metadata: &user.UserMetadata{
			Name: userMetadata.Name,
			Dates: &common.CommonDate{
				CreatedAt:   userMetadata.Dates.CreatedAt,
				CreatedBy:   userMetadata.Dates.CreatedBy,
				UpdatedAt:   userMetadata.Dates.UpdatedAt,
				UpdatedBy:   userMetadata.Dates.UpdatedBy,
				StartDate:   &userMetadata.Dates.StartDate,
				EndDate:     &userMetadata.Dates.EndDate,
				StartedAt:   &userMetadata.Dates.StartedAt,
				StartedBy:   &userMetadata.Dates.StartedBy,
				CompletedAt: &userMetadata.Dates.CompletedAt,
				CompletedBy: &userMetadata.Dates.CompletedBy,
			},
		},
	}

	uss.logger.Info(userMetadataResponse)
	return userMetadataResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Content": {"Email": "john.doe@mail.com"}}' 0.0.0.0:50051 user.UserService/UpdateUserContent
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Content": {"Email": "john.doe@mail.com"}}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/UpdateUserContent
func (uss *UserServiceServer) UpdateUserContent(ctx context.Context, in *user.UpdateUserContentRequest) (*user.UserContentResponse, error) {
	var projectRoles []model.PROJECT_ROLE_TYPES
	for _, role := range in.Content.ProjectRoles {
		projectRoles = append(projectRoles, model.PROJECT_ROLE_TYPES(role))
	}

	var scrumRoles []model.SCRUM_ROLE_TYPES
	for _, role := range in.Content.ScrumRoles {
		scrumRoles = append(scrumRoles, model.SCRUM_ROLE_TYPES(role))
	}

	updateUserContentRequest := model.UpdateUserContentRequest{
		// ID:   in.ID,
		UUID: in.UUID,
		Content: &model.UserContent{
			Email:        in.Content.Email,
			Phone:        in.Content.Phone,
			LastName:     in.Content.LastName,
			FirstName:    in.Content.FirstName,
			ProjectRoles: projectRoles,
			ScrumRoles:   scrumRoles,
			Password:     in.Content.Password,
		},
	}

	userContent, err := uss.userService.UpdateUserContent(in.UUID, updateUserContentRequest)
	if err != nil {
		return nil, err
	}

	// Assuming you have an empty slice of the target type
	var targetProjectRoles []common.PROJECT_ROLE_TYPES

	// Iterate through userContent.ProjectRoles and convert each element
	for _, role := range userContent.ProjectRoles {
		// Perform the conversion and append to the target slice
		targetRole := common.PROJECT_ROLE_TYPES(role) // Assuming there's a valid conversion
		targetProjectRoles = append(targetProjectRoles, targetRole)
	}

	// Assuming you have an empty slice of the target type
	var targetScrumRoles []common.SCRUM_ROLE_TYPES

	// Iterate through userContent.ScrumRoles and convert each element
	for _, role := range userContent.ScrumRoles {
		// Perform the conversion and append to the target slice
		targetRole := common.SCRUM_ROLE_TYPES(role) // Assuming there's a valid conversion
		targetScrumRoles = append(targetScrumRoles, targetRole)
	}

	userContentResponse := &user.UserContentResponse{
		Content: &user.UserContent{
			Email:        userContent.Email,
			Phone:        userContent.Phone,
			LastName:     userContent.LastName,
			FirstName:    userContent.FirstName,
			ProjectRoles: targetProjectRoles,
			ScrumRoles:   targetScrumRoles,
			Password:     userContent.Password,
		},
	}

	uss.logger.Info(userContentResponse)
	return userContentResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 user.UserService/GetUserMetadata
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/GetUserMetadata
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/GetUserMetadata
func (uss *UserServiceServer) GetUserMetadata(ctx context.Context, in *user.GetUserByUuidRequest) (*user.UserMetadataResponse, error) {
	userMetadata, err := uss.userService.GetUserMetadata(in.UUID)
	if err != nil {
		return nil, err
	}

	userMetadataResponse := &user.UserMetadataResponse{
		Metadata: &user.UserMetadata{
			Name: userMetadata.Name,
			Dates: &common.CommonDate{
				CreatedAt:   userMetadata.Dates.CreatedAt,
				CreatedBy:   userMetadata.Dates.CreatedBy,
				UpdatedAt:   userMetadata.Dates.UpdatedAt,
				UpdatedBy:   userMetadata.Dates.UpdatedBy,
				StartDate:   &userMetadata.Dates.StartDate,
				EndDate:     &userMetadata.Dates.EndDate,
				StartedAt:   &userMetadata.Dates.StartedAt,
				StartedBy:   &userMetadata.Dates.StartedBy,
				CompletedAt: &userMetadata.Dates.CompletedAt,
				CompletedBy: &userMetadata.Dates.CompletedBy,
			},
		},
	}

	uss.logger.Info(userMetadataResponse)
	return userMetadataResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 user.UserService/GetUserContent
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/user.proto -import-path api/grpc/proto 0.0.0.0:50051 user.UserService/GetUserContent
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/user.proto -import-path api/grpc/proto localhost:50051 user.UserService/GetUserContent
func (uss *UserServiceServer) GetUserContent(ctx context.Context, in *user.GetUserByUuidRequest) (*user.UserContentResponse, error) {
	userContent, err := uss.userService.GetUserContent(in.UUID)
	if err != nil {
		return nil, err
	}

	// Assuming you have an empty slice of the target type
	var targetProjectRoles []common.PROJECT_ROLE_TYPES

	// Iterate through userContent.ProjectRoles and convert each element
	for _, role := range userContent.ProjectRoles {
		// Perform the conversion and append to the target slice
		targetRole := common.PROJECT_ROLE_TYPES(role) // Assuming there's a valid conversion
		targetProjectRoles = append(targetProjectRoles, targetRole)
	}

	// Assuming you have an empty slice of the target type
	var targetScrumRoles []common.SCRUM_ROLE_TYPES

	// Iterate through userContent.ScrumRoles and convert each element
	for _, role := range userContent.ScrumRoles {
		// Perform the conversion and append to the target slice
		targetRole := common.SCRUM_ROLE_TYPES(role) // Assuming there's a valid conversion
		targetScrumRoles = append(targetScrumRoles, targetRole)
	}

	userContentResponse := &user.UserContentResponse{
		Content: &user.UserContent{
			Email:        userContent.Email,
			Phone:        userContent.Phone,
			LastName:     userContent.LastName,
			FirstName:    userContent.FirstName,
			ProjectRoles: targetProjectRoles,
			ScrumRoles:   targetScrumRoles,
			Password:     userContent.Password,
		},
	}

	uss.logger.Info(userContentResponse)
	return userContentResponse, nil
}

func (uss *UserServiceServer) convertUser(userModel *model.User) (*user.User, error) {
	var targetProjectRoles []common.PROJECT_ROLE_TYPES

	for _, role := range userModel.Content.ProjectRoles {

		targetRole := common.PROJECT_ROLE_TYPES(role)
		targetProjectRoles = append(targetProjectRoles, targetRole)
	}

	var targetScrumRoles []common.SCRUM_ROLE_TYPES

	for _, role := range userModel.Content.ScrumRoles {

		targetRole := common.SCRUM_ROLE_TYPES(role)
		targetScrumRoles = append(targetScrumRoles, targetRole)
	}

	userResponse := &user.User{
		ID:   userModel.ID,
		UUID: userModel.UUID,
		Metadata: &user.UserMetadata{
			Name: userModel.Metadata.Name,
			Dates: &common.CommonDate{
				CreatedAt:   userModel.Metadata.Dates.CreatedAt,
				CreatedBy:   userModel.Metadata.Dates.CreatedBy,
				UpdatedAt:   userModel.Metadata.Dates.UpdatedAt,
				UpdatedBy:   userModel.Metadata.Dates.UpdatedBy,
				StartDate:   &userModel.Metadata.Dates.StartDate,
				EndDate:     &userModel.Metadata.Dates.EndDate,
				StartedAt:   &userModel.Metadata.Dates.StartedAt,
				StartedBy:   &userModel.Metadata.Dates.StartedBy,
				CompletedAt: &userModel.Metadata.Dates.CompletedAt,
				CompletedBy: &userModel.Metadata.Dates.CompletedBy,
			},
		},
		Content: &user.UserContent{
			Email:        userModel.Content.Email,
			Phone:        userModel.Content.Phone,
			LastName:     userModel.Content.LastName,
			FirstName:    userModel.Content.FirstName,
			ProjectRoles: targetProjectRoles,
			ScrumRoles:   targetScrumRoles,
			Password:     userModel.Content.Password,
		},
	}

	uss.logger.Info(userResponse)
	return userResponse, nil
}
