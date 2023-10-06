// Path: ineternal/repository/team/memory_team_repository.go
// DESC: This is the memory implementation of the team repository.
package memory

import (
	"context"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/sirupsen/logrus"
)

// TeamRepositoryInterface defines the interface for the team repository.
type TeamRepositoryInterface interface {
	ListTeamIdsAndUUIDs() ([]model.IdUuid, error)
	ListTeams() ([]model.Team, error)
	// ListTeamsMetadata() ([]model.TeamMetadata, error)
	// ListTeamsContent() ([]model.TeamContent, error)
	GetTeam(UUID string) (*model.Team, error)
	CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error)
	UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error)
	DeleteTeam(UUID string) (*model.Team, error)
	GetTeamByID(ID string) (*model.Team, error)
	GetTeamByName(name string) (*model.Team, error)
	GetTeamByEmail(email string) (*model.Team, error)
	UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error)
	UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error)
	GetTeamMetadata(UUID string) (*model.TeamMetadata, error)
	GetTeamContent(UUID string) (*model.TeamContent, error)
}

type MemoryTeamRepository struct {
	ctx	   context.Context
	logger *logrus.Logger
}

func NewMemoryTeamRepository(ctx context.Context, logger *logrus.Logger) *MemoryTeamRepository {
	return &MemoryTeamRepository{
		ctx:	ctx,
		logger: logger,
	}
}

func (repository *MemoryTeamRepository) ListTeamIdsAndUUIDs() ([]model.IdUuid, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) ListTeams() ([]model.Team, error) {
	return nil, nil
}

// func (repository *MemoryTeamRepository) ListTeamsMetadata() ([]model.TeamMetadata, error) {
//	return nil, nil
// }

// func (repository *MemoryTeamRepository) ListTeamsContent() ([]model.TeamContent, error) {
//	return nil, nil
// }

func (repository *MemoryTeamRepository) GetTeam(UUID string) (*model.Team, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) DeleteTeam(UUID string) (*model.Team, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) GetTeamByID(ID string) (*model.Team, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) GetTeamByName(name string) (*model.Team, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) GetTeamByEmail(email string) (*model.Team, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) GetTeamMetadata(UUID string) (*model.TeamMetadata, error) {
	return nil, nil
}

func (repository *MemoryTeamRepository) GetTeamContent(UUID string) (*model.TeamContent, error) {
	return nil, nil
}
