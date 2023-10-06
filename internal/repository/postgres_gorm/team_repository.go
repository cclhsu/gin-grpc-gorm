// Path: internal/repository/team/memory_team_repository.go
// DESC: This is the memory implementation of the team repository.
package postgres_gorm

import (
	"context"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

// PostgresTeamRepository implements TeamRepositoryInterface for PostgreSQL.
type PostgresTeamRepository struct {
	ctx	   context.Context
	logger *logrus.Logger
	db	   *gorm.DB
}

// NewPostgresTeamRepository creates a new instance of PostgresTeamRepository.
func NewPostgresTeamRepository(ctx context.Context, logger *logrus.Logger, db *gorm.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{
		ctx:	ctx,
		logger: logger,
		db:		db,
	}
}

// // Error messages
// var (
//	ErrNotFound		= errors.New("record not found")
//	ErrDatabase		= errors.New("database error")
//	ErrAlreadyExist = errors.New("record already exists")
// )

// // handleErr is a helper function to handle errors consistently.
// func handleErr(err error) error {
//	if err == nil {
//		return nil
//	}

//	if errors.Is(err, gorm.ErrRecordNotFound) {
//		return ErrNotFound
//	}

//	return ErrDatabase
// }

// ListTeamIdsAndUUIDs returns a list of team IDs and UUIDs.
func (repository *PostgresTeamRepository) ListTeamIdsAndUUIDs() ([]model.IdUuid, error) {
	var idUuids []model.IdUuid
	if err := repository.db.WithContext(repository.ctx).Find(&idUuids).Error; err != nil {
		return nil, handleErr(err)
	}
	return idUuids, nil
}

// ListTeams returns a list of teams.
func (repository *PostgresTeamRepository) ListTeams() ([]model.Team, error) {
	var teams []model.Team
	if err := repository.db.WithContext(repository.ctx).Find(&teams).Error; err != nil {
		return nil, handleErr(err)
	}
	return teams, nil
}

// GetTeam retrieves a team by UUID.
func (repository *PostgresTeamRepository) GetTeam(UUID string) (*model.Team, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return &team, nil
}

// CreateTeam creates a new team.
func (repository *PostgresTeamRepository) CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error) {
	team := model.Team{
		ID:	  createTeamRequest.ID,
		UUID: createTeamRequest.UUID,
		Metadata: &model.TeamMetadata{
			Name: createTeamRequest.Metadata.Name,
			Dates: &model.CommonDate{
				CreatedAt:	 createTeamRequest.Metadata.Dates.CreatedAt,
				CreatedBy:	 createTeamRequest.Metadata.Dates.CreatedBy,
				UpdatedAt:	 createTeamRequest.Metadata.Dates.UpdatedAt,
				UpdatedBy:	 createTeamRequest.Metadata.Dates.UpdatedBy,
				StartDate:	 createTeamRequest.Metadata.Dates.StartDate,
				EndDate:	 createTeamRequest.Metadata.Dates.EndDate,
				StartedAt:	 createTeamRequest.Metadata.Dates.StartedAt,
				StartedBy:	 createTeamRequest.Metadata.Dates.StartedBy,
				CompletedAt: createTeamRequest.Metadata.Dates.CompletedAt,
				CompletedBy: createTeamRequest.Metadata.Dates.CompletedBy,
			},
		},
		Content: &model.TeamContent{
			Email:		  createTeamRequest.Content.Email,
			ProductOwner: createTeamRequest.Content.ProductOwner,
			ScrumMaster:  createTeamRequest.Content.ScrumMaster,
			Members:	  createTeamRequest.Content.Members,
		},
	}

	if err := repository.db.WithContext(repository.ctx).Create(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return &team, nil
}

// UpdateTeam updates a team by UUID.
func (repository *PostgresTeamRepository) UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}

	team.Metadata.Name = updateTeamRequest.Metadata.Name
	team.Metadata.Dates.UpdatedAt = updateTeamRequest.Metadata.Dates.UpdatedAt
	team.Metadata.Dates.UpdatedBy = updateTeamRequest.Metadata.Dates.UpdatedBy
	team.Metadata.Dates.StartDate = updateTeamRequest.Metadata.Dates.StartDate
	team.Metadata.Dates.EndDate = updateTeamRequest.Metadata.Dates.EndDate
	team.Metadata.Dates.StartedAt = updateTeamRequest.Metadata.Dates.StartedAt
	team.Metadata.Dates.StartedBy = updateTeamRequest.Metadata.Dates.StartedBy
	team.Metadata.Dates.CompletedAt = updateTeamRequest.Metadata.Dates.CompletedAt
	team.Metadata.Dates.CompletedBy = updateTeamRequest.Metadata.Dates.CompletedBy
	team.Content.Email = updateTeamRequest.Content.Email
	team.Content.ProductOwner = updateTeamRequest.Content.ProductOwner
	team.Content.ScrumMaster = updateTeamRequest.Content.ScrumMaster
	team.Content.Members = updateTeamRequest.Content.Members

	if err := repository.db.WithContext(repository.ctx).Save(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return &team, nil
}

// DeleteTeam deletes a team by UUID.
func (repository *PostgresTeamRepository) DeleteTeam(UUID string) (*model.Team, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}

	if err := repository.db.WithContext(repository.ctx).Delete(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return &team, nil
}

// GetTeamByID retrieves a team by ID.
func (repository *PostgresTeamRepository) GetTeamByID(ID string) (*model.Team, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("id = ?", ID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return &team, nil
}

// GetTeamByName retrieves a team by name.
func (repository *PostgresTeamRepository) GetTeamByName(name string) (*model.Team, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("name = ?", name).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return &team, nil
}

// GetTeamByEmail retrieves a team by email.
func (repository *PostgresTeamRepository) GetTeamByEmail(email string) (*model.Team, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("email = ?", email).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return &team, nil
}

// UpdateTeamMetadata updates team metadata by UUID.
func (repository *PostgresTeamRepository) UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}

	team.Metadata.Name = updatedMetadata.Name
	team.Metadata.Dates.UpdatedAt = updatedMetadata.Dates.UpdatedAt
	team.Metadata.Dates.UpdatedBy = updatedMetadata.Dates.UpdatedBy
	team.Metadata.Dates.StartDate = updatedMetadata.Dates.StartDate
	team.Metadata.Dates.EndDate = updatedMetadata.Dates.EndDate
	team.Metadata.Dates.StartedAt = updatedMetadata.Dates.StartedAt
	team.Metadata.Dates.StartedBy = updatedMetadata.Dates.StartedBy
	team.Metadata.Dates.CompletedAt = updatedMetadata.Dates.CompletedAt
	team.Metadata.Dates.CompletedBy = updatedMetadata.Dates.CompletedBy

	if err := repository.db.WithContext(repository.ctx).Save(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return team.Metadata, nil
}

// UpdateTeamContent updates team content by UUID.
func (repository *PostgresTeamRepository) UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}

	team.Content.Email = updatedContent.Email
	team.Content.ProductOwner = updatedContent.ProductOwner
	team.Content.ScrumMaster = updatedContent.ScrumMaster
	team.Content.Members = updatedContent.Members

	if err := repository.db.WithContext(repository.ctx).Save(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return team.Content, nil
}

// GetTeamMetadata retrieves team metadata by UUID.
func (repository *PostgresTeamRepository) GetTeamMetadata(UUID string) (*model.TeamMetadata, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return team.Metadata, nil
}

// GetTeamContent retrieves team content by UUID.
func (repository *PostgresTeamRepository) GetTeamContent(UUID string) (*model.TeamContent, error) {
	var team model.Team
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&team).Error; err != nil {
		return nil, handleErr(err)
	}
	return team.Content, nil
}
