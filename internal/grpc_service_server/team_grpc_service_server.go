package grpc_service_server

import (
	"context"

	"github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/common"
	"github.com/cclhsu/gin-grpc-gorm/generated/grpc/pb/team"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"

	// "github.com/cclhsu/gin-grpc-gorm/internal/model"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
)

// type TeamServiceClient interface {
// 	ListTeamIdsAndUUIDs(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTeamIdUuid, error)
// 	ListTeams(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTeamsResponse, error)
// 	ListTeamsMetadata(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTeamsMetadataResponse, error)
// 	ListTeamsContent(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*ListTeamsContentResponse, error)
// 	GetTeam(ctx context.Context, in *GetTeamByUuidRequest, opts ...grpc.CallOption) (*Team, error)
// 	CreateTeam(ctx context.Context, in *CreateTeamRequest, opts ...grpc.CallOption) (*Team, error)
// 	UpdateTeam(ctx context.Context, in *UpdateTeamRequest, opts ...grpc.CallOption) (*Team, error)
// 	DeleteTeam(ctx context.Context, in *GetTeamByUuidRequest, opts ...grpc.CallOption) (*Team, error)
// 	GetTeamById(ctx context.Context, in *GetTeamByIdRequest, opts ...grpc.CallOption) (*Team, error)
// 	GetTeamByName(ctx context.Context, in *GetTeamByTeamnameRequest, opts ...grpc.CallOption) (*Team, error)
// 	GetTeamByEmail(ctx context.Context, in *GetTeamByEmailRequest, opts ...grpc.CallOption) (*Team, error)
// 	UpdateTeamMetadata(ctx context.Context, in *UpdateTeamMetadataRequest, opts ...grpc.CallOption) (*TeamMetadataResponse, error)
// 	UpdateTeamContent(ctx context.Context, in *UpdateTeamContentRequest, opts ...grpc.CallOption) (*TeamContentResponse, error)
// 	GetTeamMetadata(ctx context.Context, in *GetTeamByUuidRequest, opts ...grpc.CallOption) (*TeamMetadataResponse, error)
// 	GetTeamContent(ctx context.Context, in *GetTeamByUuidRequest, opts ...grpc.CallOption) (*TeamContentResponse, error)
// }

type TeamServiceServer struct {
	ctx    context.Context
	logger *logrus.Logger
	team.UnimplementedTeamServiceServer
	teamService *service.TeamService
}

func NewTeamServiceServer(ctx context.Context, logger *logrus.Logger, hs *service.TeamService) *TeamServiceServer {
	return &TeamServiceServer{
		ctx:         ctx,
		logger:      logger,
		teamService: hs,
	}
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 team.TeamService/ListTeamIdsAndUUIDs
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/ListTeamIdsAndUUIDs
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/ListTeamIdsAndUUIDs
func (uss *TeamServiceServer) ListTeamIdsAndUUIDs(ctx context.Context, in *emptypb.Empty) (*team.ListTeamIdUuid, error) {
	idUuids, err := uss.teamService.ListTeamIdsAndUUIDs()
	if err != nil {
		return nil, err
	}

	// convert []*model.IdUuid to team.ListTeamIdUuid
	listTeamIdUuid := &team.ListTeamIdUuid{}
	for _, idUuid := range idUuids {
		listTeamIdUuid.TeamIdUuids = append(listTeamIdUuid.TeamIdUuids, &common.IdUuid{
			ID:   idUuid.ID,
			UUID: idUuid.UUID,
		})
	}

	uss.logger.Info(listTeamIdUuid)
	return listTeamIdUuid, nil
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 team.TeamService/ListTeams
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/ListTeams
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/ListTeams
func (uss *TeamServiceServer) ListTeams(ctx context.Context, in *emptypb.Empty) (*team.ListTeamsResponse, error) {
	teams, err := uss.teamService.ListTeams()
	if err != nil {
		return nil, err
	}

	// convert []*model.TeamResponse to team.ListTeamsResponse
	listTeamsResponse := &team.ListTeamsResponse{}
	for _, teamModel := range teams {

		productOwner := &common.IdUuid{
			ID:   teamModel.Content.ProductOwner.ID,
			UUID: teamModel.Content.ProductOwner.UUID,
		}

		scrumMaster := &common.IdUuid{
			ID:   teamModel.Content.ScrumMaster.ID,
			UUID: teamModel.Content.ScrumMaster.UUID,
		}

		var members []*common.IdUuid
		for _, member := range teamModel.Content.Members {
			members = append(members, &common.IdUuid{
				ID:   member.ID,
				UUID: member.UUID,
			})
		}

		listTeamsResponse.Teams = append(listTeamsResponse.Teams, &team.Team{
			ID:   teamModel.ID,
			UUID: teamModel.UUID,
			Metadata: &team.TeamMetadata{
				Name: teamModel.Metadata.Name,
				Dates: &common.CommonDate{
					CreatedAt:   teamModel.Metadata.Dates.CreatedAt,
					CreatedBy:   teamModel.Metadata.Dates.CreatedBy,
					UpdatedAt:   teamModel.Metadata.Dates.UpdatedAt,
					UpdatedBy:   teamModel.Metadata.Dates.UpdatedBy,
					StartDate:   &teamModel.Metadata.Dates.StartDate,
					EndDate:     &teamModel.Metadata.Dates.EndDate,
					StartedAt:   &teamModel.Metadata.Dates.StartedAt,
					StartedBy:   &teamModel.Metadata.Dates.StartedBy,
					CompletedAt: &teamModel.Metadata.Dates.CompletedAt,
					CompletedBy: &teamModel.Metadata.Dates.CompletedBy,
				},
			},
			Content: &team.TeamContent{
				Email:        teamModel.Content.Email,
				ProductOwner: productOwner,
				ScrumMaster:  scrumMaster,
				Members:      members,
			},
		})
	}

	uss.logger.Info(listTeamsResponse)
	return listTeamsResponse, nil
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 team.TeamService/ListTeamsMetadata
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/ListTeamsMetadata
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/ListTeamsMetadata
func (uss *TeamServiceServer) ListTeamsMetadata(ctx context.Context, in *emptypb.Empty) (*team.ListTeamsMetadataResponse, error) {
	listTeamsMetadata, err := uss.teamService.ListTeamsMetadata()
	if err != nil {
		return nil, err
	}

	// convert []*model.TeamMetadata to team.ListTeamsMetadataResponse
	listTeamsMetadataResponse := &team.ListTeamsMetadataResponse{}
	for _, teamModel := range listTeamsMetadata.TeamMetadataResponses {
		listTeamsMetadataResponse.TeamMetadataResponses = append(listTeamsMetadataResponse.TeamMetadataResponses, &team.TeamMetadataResponse{
			ID:   teamModel.ID,
			UUID: teamModel.UUID,
			Metadata: &team.TeamMetadata{
				Name: teamModel.Metadata.Name,
				Dates: &common.CommonDate{
					CreatedAt:   teamModel.Metadata.Dates.CreatedAt,
					CreatedBy:   teamModel.Metadata.Dates.CreatedBy,
					UpdatedAt:   teamModel.Metadata.Dates.UpdatedAt,
					UpdatedBy:   teamModel.Metadata.Dates.UpdatedBy,
					StartDate:   &teamModel.Metadata.Dates.StartDate,
					EndDate:     &teamModel.Metadata.Dates.EndDate,
					StartedAt:   &teamModel.Metadata.Dates.StartedAt,
					StartedBy:   &teamModel.Metadata.Dates.StartedBy,
					CompletedAt: &teamModel.Metadata.Dates.CompletedAt,
					CompletedBy: &teamModel.Metadata.Dates.CompletedBy,
				},
			},
		})
	}

	uss.logger.Info(listTeamsMetadataResponse)
	return listTeamsMetadataResponse, nil
}

// grpcurl -plaintext -d '{}' 0.0.0.0:50051 team.TeamService/ListTeamsContent
// grpcurl -plaintext -d '{}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/ListTeamsContent
// grpcurl -plaintext -d '{}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/ListTeamsContent
func (uss *TeamServiceServer) ListTeamsContent(ctx context.Context, in *emptypb.Empty) (*team.ListTeamsContentResponse, error) {
	listTeamsContent, err := uss.teamService.ListTeamsContent()
	if err != nil {
		return nil, err
	}

	// convert []*model.TeamContent to team.ListTeamsContentResponse
	listTeamsContentResponse := &team.ListTeamsContentResponse{}
	for _, teamModel := range listTeamsContent.TeamContentResponses {

		productOwner := &common.IdUuid{
			ID:   teamModel.Content.ProductOwner.ID,
			UUID: teamModel.Content.ProductOwner.UUID,
		}

		scrumMaster := &common.IdUuid{
			ID:   teamModel.Content.ScrumMaster.ID,
			UUID: teamModel.Content.ScrumMaster.UUID,
		}

		var members []*common.IdUuid
		for _, member := range teamModel.Content.Members {
			members = append(members, &common.IdUuid{
				ID:   member.ID,
				UUID: member.UUID,
			})
		}

		listTeamsContentResponse.TeamContentResponses = append(listTeamsContentResponse.TeamContentResponses, &team.TeamContentResponse{
			ID:   teamModel.ID,
			UUID: teamModel.UUID,
			Content: &team.TeamContent{
				Email:        teamModel.Content.Email,
				ProductOwner: productOwner,
				ScrumMaster:  scrumMaster,
				Members:      members,
			},
		})
	}

	uss.logger.Info(listTeamsContentResponse)
	return listTeamsContentResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 team.TeamService/GetTeam
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/GetTeam
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/GetTeam
func (uss *TeamServiceServer) GetTeam(ctx context.Context, in *team.GetTeamByUuidRequest) (*team.Team, error) {
	teamModel, err := uss.teamService.GetTeam(in.UUID)
	if err != nil {
		return nil, err
	}

	return uss.convertTeam(teamModel)
}

// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' 0.0.0.0:50051 team.TeamService/CreateTeam
// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/CreateTeam
func (uss *TeamServiceServer) CreateTeam(ctx context.Context, in *team.CreateTeamRequest) (*team.Team, error) {

	productOwner := &model.IdUuid{
		ID:   in.Content.ProductOwner.ID,
		UUID: in.Content.ProductOwner.UUID,
	}

	scrumMaster := &model.IdUuid{
		ID:   in.Content.ScrumMaster.ID,
		UUID: in.Content.ScrumMaster.UUID,
	}

	var members []*model.IdUuid
	for _, member := range in.Content.Members {
		members = append(members, &model.IdUuid{
			ID:   member.ID,
			UUID: member.UUID,
		})
	}

	createTeamRequest := model.CreateTeamRequest{
		ID:   in.ID,
		UUID: in.UUID,
		Metadata: &model.TeamMetadata{
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
		Content: &model.TeamContent{
			Email:        in.Content.Email,
			ProductOwner: productOwner,
			ScrumMaster:  scrumMaster,
			Members:      members,
		},
	}

	teamModel, err := uss.teamService.CreateTeam(createTeamRequest)
	if err != nil {
		return nil, err
	}

	return uss.convertTeam(teamModel)
}

// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' 0.0.0.0:50051 team.TeamService/UpdateTeam
// grpcurl -plaintext -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}, "Content": {"Email": "john.doe@mail.com", "Phone": "0987-654-321", "LastName": "Doe", "FirstName": "John", "ProjectRoles": ["PROJECT_ROLE_TYPES_PROJECT_MANAGER"], "ScrumRoles": ["SCRUM_ROLE_TYPES_SCRUM_MASTER"], "Password": "password"}}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/UpdateTeam
func (uss *TeamServiceServer) UpdateTeam(ctx context.Context, in *team.UpdateTeamRequest) (*team.Team, error) {
	productOwner := &model.IdUuid{
		ID:   in.Content.ProductOwner.ID,
		UUID: in.Content.ProductOwner.UUID,
	}

	scrumMaster := &model.IdUuid{
		ID:   in.Content.ScrumMaster.ID,
		UUID: in.Content.ScrumMaster.UUID,
	}

	var members []*model.IdUuid
	for _, member := range in.Content.Members {
		members = append(members, &model.IdUuid{
			ID:   member.ID,
			UUID: member.UUID,
		})
	}

	updateTeamRequest := model.UpdateTeamRequest{
		// ID:   in.ID,
		UUID: in.UUID,
		Metadata: &model.TeamMetadata{
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
		Content: &model.TeamContent{
			Email:        in.Content.Email,
			ProductOwner: productOwner,
			ScrumMaster:  scrumMaster,
			Members:      members,
		},
	}

	teamModel, err := uss.teamService.UpdateTeam(in.UUID, updateTeamRequest)
	if err != nil {
		return nil, err
	}

	return uss.convertTeam(teamModel)
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 team.TeamService/DeleteTeam
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/DeleteTeam
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/DeleteTeam
func (uss *TeamServiceServer) DeleteTeam(ctx context.Context, in *team.GetTeamByUuidRequest) (*team.Team, error) {
	teamModel, err := uss.teamService.DeleteTeam(in.UUID)
	if err != nil {
		return nil, err
	}

	return uss.convertTeam(teamModel)
}

// grpcurl -plaintext -d '{"ID": "john.doe"}' 0.0.0.0:50051 team.TeamService/GetTeamById
// grpcurl -plaintext -d '{"ID": "john.doe"}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/GetTeamById
// grpcurl -plaintext -d '{"ID": "john.doe"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/GetTeamById
func (uss *TeamServiceServer) GetTeamById(ctx context.Context, in *team.GetTeamByIdRequest) (*team.Team, error) {
	teamModel, err := uss.teamService.GetTeamByID(in.ID)
	if err != nil {
		return nil, err
	}

	return uss.convertTeam(teamModel)
}

// grpcurl -plaintext -d '{"name": "John Doe"}' 0.0.0.0:50051 team.TeamService/GetTeamByName
// grpcurl -plaintext -d '{"name": "John Doe"}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/GetTeamByName
// grpcurl -plaintext -d '{"name": "John Doe"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/GetTeamByName
func (uss *TeamServiceServer) GetTeamByName(ctx context.Context, in *team.GetTeamByNameRequest) (*team.Team, error) {
	teamModel, err := uss.teamService.GetTeamByName(in.Name)
	if err != nil {
		return nil, err
	}

	return uss.convertTeam(teamModel)
}

// grpcurl -plaintext -d '{"email": "john.doe@mail.com"}' 0.0.0.0:50051 team.TeamService/GetTeamByEmail
// grpcurl -plaintext -d '{"email": "john.doe@mail.com"}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/GetTeamByEmail
// grpcurl -plaintext -d '{"email": "john.doe@mail.com"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/GetTeamByEmail
func (uss *TeamServiceServer) GetTeamByEmail(ctx context.Context, in *team.GetTeamByEmailRequest) (*team.Team, error) {
	teamModel, err := uss.teamService.GetTeamByEmail(in.Email)
	if err != nil {
		return nil, err
	}

	return uss.convertTeam(teamModel)
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}}' 0.0.0.0:50051 team.TeamService/UpdateTeamMetadata
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Metadata": {"Name": "John Doe", "Dates": {"CreatedAt": "2021-01-01T00:00:00Z", "CreatedBy": "John Doe", "UpdatedAt": "2021-01-01T00:00:00Z", "UpdatedBy": "John Doe", "StartDate": "2021-01-01T00:00:00Z", "EndDate": "2021-01-01T00:00:00Z", "StartedAt": "2021-01-01T00:00:00Z", "StartedBy": "John Doe", "CompletedAt": "2021-01-01T00:00:00Z", "CompletedBy": "John Doe"}}}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/UpdateTeamMetadata
func (uss *TeamServiceServer) UpdateTeamMetadata(ctx context.Context, in *team.UpdateTeamMetadataRequest) (*team.TeamMetadataResponse, error) {
	updateTeamMetadataRequest := model.UpdateTeamMetadataRequest{
		// ID:   in.ID,
		UUID: in.UUID,
		Metadata: &model.TeamMetadata{
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

	teamMetadata, err := uss.teamService.UpdateTeamMetadata(in.UUID, updateTeamMetadataRequest)
	if err != nil {
		return nil, err
	}

	teamMetadataResponse := &team.TeamMetadataResponse{
		Metadata: &team.TeamMetadata{
			Name: teamMetadata.Name,
			Dates: &common.CommonDate{
				CreatedAt:   teamMetadata.Dates.CreatedAt,
				CreatedBy:   teamMetadata.Dates.CreatedBy,
				UpdatedAt:   teamMetadata.Dates.UpdatedAt,
				UpdatedBy:   teamMetadata.Dates.UpdatedBy,
				StartDate:   &teamMetadata.Dates.StartDate,
				EndDate:     &teamMetadata.Dates.EndDate,
				StartedAt:   &teamMetadata.Dates.StartedAt,
				StartedBy:   &teamMetadata.Dates.StartedBy,
				CompletedAt: &teamMetadata.Dates.CompletedAt,
				CompletedBy: &teamMetadata.Dates.CompletedBy,
			},
		},
	}

	uss.logger.Info(teamMetadataResponse)
	return teamMetadataResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Content": {"Email": "john.doe@mail.com"}}' 0.0.0.0:50051 team.TeamService/UpdateTeamContent
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000", "Content": {"Email": "john.doe@mail.com"}}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/UpdateTeamContent
func (uss *TeamServiceServer) UpdateTeamContent(ctx context.Context, in *team.UpdateTeamContentRequest) (*team.TeamContentResponse, error) {
	productOwner := &model.IdUuid{
		ID:   in.Content.ProductOwner.ID,
		UUID: in.Content.ProductOwner.UUID,
	}

	scrumMaster := &model.IdUuid{
		ID:   in.Content.ScrumMaster.ID,
		UUID: in.Content.ScrumMaster.UUID,
	}

	var members []*model.IdUuid
	for _, member := range in.Content.Members {
		members = append(members, &model.IdUuid{
			ID:   member.ID,
			UUID: member.UUID,
		})
	}

	updateTeamContentRequest := model.UpdateTeamContentRequest{
		// ID:   in.ID,
		UUID: in.UUID,
		Content: &model.TeamContent{
			Email:        in.Content.Email,
			ProductOwner: productOwner,
			ScrumMaster:  scrumMaster,
			Members:      members,
		},
	}

	teamContent, err := uss.teamService.UpdateTeamContent(in.UUID, updateTeamContentRequest)
	if err != nil {
		return nil, err
	}

	// productOwner = common.IdUuid{
	// 	ID:   teamContent.ProductOwner.ID,
	// 	UUID: teamContent.ProductOwner.UUID,
	// }

	// scrumMaster = common.IdUuid{
	// 	ID:   teamContent.ScrumMaster.ID,
	// 	UUID: teamContent.ScrumMaster.UUID,
	// }

	// var members []*common.IdUuid
	// for _, member := range teamContent.Members {
	// 	members = append(members, &common.IdUuid{
	// 		ID:   member.ID,
	// 		UUID: member.UUID,
	// 	})
	// }

	teamContentResponse := &team.TeamContentResponse{
		Content: &team.TeamContent{
			Email: teamContent.Email,
			// ProductOwner: productOwner,
			// ScrumMaster:  scrumMaster,
			// Members:      members,
		},
	}

	uss.logger.Info(teamContentResponse)
	return teamContentResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 team.TeamService/GetTeamMetadata
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/GetTeamMetadata
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/GetTeamMetadata
func (uss *TeamServiceServer) GetTeamMetadata(ctx context.Context, in *team.GetTeamByUuidRequest) (*team.TeamMetadataResponse, error) {
	teamMetadata, err := uss.teamService.GetTeamMetadata(in.UUID)
	if err != nil {
		return nil, err
	}

	teamMetadataResponse := &team.TeamMetadataResponse{
		Metadata: &team.TeamMetadata{
			Name: teamMetadata.Name,
			Dates: &common.CommonDate{
				CreatedAt:   teamMetadata.Dates.CreatedAt,
				CreatedBy:   teamMetadata.Dates.CreatedBy,
				UpdatedAt:   teamMetadata.Dates.UpdatedAt,
				UpdatedBy:   teamMetadata.Dates.UpdatedBy,
				StartDate:   &teamMetadata.Dates.StartDate,
				EndDate:     &teamMetadata.Dates.EndDate,
				StartedAt:   &teamMetadata.Dates.StartedAt,
				StartedBy:   &teamMetadata.Dates.StartedBy,
				CompletedAt: &teamMetadata.Dates.CompletedAt,
				CompletedBy: &teamMetadata.Dates.CompletedBy,
			},
		},
	}

	uss.logger.Info(teamMetadataResponse)
	return teamMetadataResponse, nil
}

// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' 0.0.0.0:50051 team.TeamService/GetTeamContent
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -proto api/grpc/proto/team.proto -import-path api/grpc/proto 0.0.0.0:50051 team.TeamService/GetTeamContent
// grpcurl -plaintext -d '{"UUID": "00000000-0000-0000-0000-000000000000"}' -H "Authorization: Bearer YOUR_TOKEN" -proto api/grpc/proto/team.proto -import-path api/grpc/proto localhost:50051 team.TeamService/GetTeamContent
func (uss *TeamServiceServer) GetTeamContent(ctx context.Context, in *team.GetTeamByUuidRequest) (*team.TeamContentResponse, error) {
	teamContent, err := uss.teamService.GetTeamContent(in.UUID)
	if err != nil {
		return nil, err
	}

	productOwner := &common.IdUuid{
		ID:   teamContent.ProductOwner.ID,
		UUID: teamContent.ProductOwner.UUID,
	}

	scrumMaster := &common.IdUuid{
		ID:   teamContent.ScrumMaster.ID,
		UUID: teamContent.ScrumMaster.UUID,
	}

	var members []*common.IdUuid
	for _, member := range teamContent.Members {
		members = append(members, &common.IdUuid{
			ID:   member.ID,
			UUID: member.UUID,
		})
	}

	teamContentResponse := &team.TeamContentResponse{
		Content: &team.TeamContent{
			Email:        teamContent.Email,
			ProductOwner: productOwner,
			ScrumMaster:  scrumMaster,
		},
	}

	uss.logger.Info(teamContentResponse)
	return teamContentResponse, nil
}

func (uss *TeamServiceServer) convertTeam(teamModel *model.Team) (*team.Team, error) {

	productOwner := &common.IdUuid{
		ID:   teamModel.Content.ProductOwner.ID,
		UUID: teamModel.Content.ProductOwner.UUID,
	}

	scrumMaster := &common.IdUuid{
		ID:   teamModel.Content.ScrumMaster.ID,
		UUID: teamModel.Content.ScrumMaster.UUID,
	}

	var members []*common.IdUuid
	for _, member := range teamModel.Content.Members {
		members = append(members, &common.IdUuid{
			ID:   member.ID,
			UUID: member.UUID,
		})
	}

	teamResponse := &team.Team{
		ID:   teamModel.ID,
		UUID: teamModel.UUID,
		Metadata: &team.TeamMetadata{
			Name: teamModel.Metadata.Name,
			Dates: &common.CommonDate{
				CreatedAt:   teamModel.Metadata.Dates.CreatedAt,
				CreatedBy:   teamModel.Metadata.Dates.CreatedBy,
				UpdatedAt:   teamModel.Metadata.Dates.UpdatedAt,
				UpdatedBy:   teamModel.Metadata.Dates.UpdatedBy,
				StartDate:   &teamModel.Metadata.Dates.StartDate,
				EndDate:     &teamModel.Metadata.Dates.EndDate,
				StartedAt:   &teamModel.Metadata.Dates.StartedAt,
				StartedBy:   &teamModel.Metadata.Dates.StartedBy,
				CompletedAt: &teamModel.Metadata.Dates.CompletedAt,
				CompletedBy: &teamModel.Metadata.Dates.CompletedBy,
			},
		},
		Content: &team.TeamContent{
			Email:        teamModel.Content.Email,
			ProductOwner: productOwner,
			ScrumMaster:  scrumMaster,
			Members:      members,
		},
	}

	uss.logger.Info(teamResponse)
	return teamResponse, nil
}
