package service

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	redis_cache "github.com/cclhsu/gin-grpc-gorm/internal/cache/redis"
	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	repository "github.com/cclhsu/gin-grpc-gorm/internal/repository/postgres_gorm"
	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
	"github.com/sirupsen/logrus"
)

// TeamServiceInterface defines the interface for the team service.
type TeamServiceInterface interface {
	// Define your repository methods here...
	ListTeamIdsAndUUIDs() ([]model.IdUuid, error)
	ListTeams() ([]*model.Team, error)
	GetTeam(uuid string) (*model.Team, error)
	CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error)
	UpdateTeam(uuid string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error)
	DeleteTeam(uuid string) (*model.Team, error)
	GetTeamByID(ID string) (*model.Team, error)
	GetTeamByName(name string) (*model.Team, error)
	GetTeamByEmail(email string) (*model.Team, error)
	ListTeamsMetadata() (*model.ListTeamsMetadataResponse, error)
	ListTeamsContent() (*model.ListTeamsContentResponse, error)
	UpdateTeamMetadata(uuid string, updateTeamMetadataRequest model.UpdateTeamMetadataRequest) (*model.TeamMetadata, error)
	UpdateTeamContent(uuid string, updateTeamContentRequest model.UpdateTeamContentRequest) (*model.TeamContent, error)
	GetTeamMetadata(uuid string) (*model.TeamMetadata, error)
	GetTeamContent(uuid string) (*model.TeamContent, error)
	IsTeamExist(name, email, ID, UUID string) (bool, error)
	IsNoTeamExist(name, email, ID, UUID string) (bool, error)
	IsExactlyOneTeamExist(name, email, ID, UUID string) (bool, error)
	IsAtLeastOneTeamExist(name, email, ID, UUID string) (bool, error)
}

// ErrTeamNotFound is a custom error for "Team not found" cases.
var ErrTeamNotFound = errors.New("User not found")

// TeamService represents the TeamService type.
type TeamService struct {
	ctx			   context.Context
	logger		   *logrus.Logger
	teamRepository repository.TeamRepositoryInterface
	cacheManager   *cache.Cache
	redisCache	   *redis_cache.RedisCache
	isCacheEnabled bool
}

// NewTeamService creates a new instance of TeamService.
func NewTeamService(ctx context.Context, logger *logrus.Logger, teamRepository repository.TeamRepositoryInterface, redisCache *redis_cache.RedisCache) *TeamService {
	cacheManager := cache.New(cache.NoExpiration, cache.NoExpiration)
	isCacheEnabled := os.Getenv("CACHE_ENABLED") == "true"

	return &TeamService{
		ctx:			ctx,
		logger:			logger,
		teamRepository: teamRepository,
		cacheManager:	cacheManager,
		redisCache:		redisCache,
		isCacheEnabled: isCacheEnabled,
	}
}

// ListTeamIdsAndUUIDs lists team IDs and UUIDs.
func (ts *TeamService) ListTeamIdsAndUUIDs() ([]model.IdUuid, error) {
	if ts.isCacheEnabled {
		if teams, found := ts.cacheManager.Get("teamsIdsAndUUIDs"); found {
			return teams.([]model.IdUuid), nil
		}
	}

	teams, err := ts.teamRepository.ListTeamIdsAndUUIDs()
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set("teamsIdsAndUUIDs", teams, cache.DefaultExpiration)
	}

	ts.logger.Printf("Teams: %+v\n", teams)
	return teams, nil
}

// ListTeams lists teams.
func (ts *TeamService) ListTeams() ([]*model.Team, error) {
	if ts.isCacheEnabled {
		if teams, found := ts.cacheManager.Get("teams"); found {
			return teams.([]*model.Team), nil
		}
	}

	teams, err := ts.teamRepository.ListTeams()
	if err != nil {
		return nil, err
	}

	teamsArray := []*model.Team{}
	for _, team := range teams {
		teamsArray = append(teamsArray, &model.Team{
			ID:	  team.ID,
			UUID: team.UUID,
			Metadata: &model.TeamMetadata{
				Name: team.Metadata.Name,
				Dates: &model.CommonDate{
					CreatedAt:	 team.Metadata.Dates.CreatedAt,
					CreatedBy:	 team.Metadata.Dates.CreatedBy,
					UpdatedAt:	 team.Metadata.Dates.UpdatedAt,
					UpdatedBy:	 team.Metadata.Dates.UpdatedBy,
					StartDate:	 team.Metadata.Dates.StartDate,
					EndDate:	 team.Metadata.Dates.EndDate,
					StartedAt:	 team.Metadata.Dates.StartedAt,
					StartedBy:	 team.Metadata.Dates.StartedBy,
					CompletedAt: team.Metadata.Dates.CompletedAt,
					CompletedBy: team.Metadata.Dates.CompletedBy,
				},
			},
			Content: &model.TeamContent{
				Email:		  team.Content.Email,
				ProductOwner: team.Content.ProductOwner,
				ScrumMaster:  team.Content.ScrumMaster,
				Members:	  team.Content.Members,
			},
		})
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set("teams", teamsArray, cache.DefaultExpiration)
	}

	ts.logger.Printf("Teams: %+v\n", teamsArray)
	return teamsArray, nil
}

// GetTeam gets a team by UUID.
func (ts *TeamService) GetTeam(uuid string) (*model.Team, error) {
	if ts.isCacheEnabled {
		if team, found := ts.cacheManager.Get(fmt.Sprintf("team:%s", uuid)); found {
			return team.(*model.Team), nil
		}
	}

	team, err := ts.teamRepository.GetTeam(uuid)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set(fmt.Sprintf("team:%s", uuid), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("Team: %+v\n", team)
	return team, nil
}

// CreateTeam creates a team.
func (ts *TeamService) CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error) {
	// check if UUID is empty string, undefined, null, and '00000000-0000-0000-0000-000000000000', generate a new UUID if so
	if createTeamRequest.UUID == "" || createTeamRequest.UUID == "undefined" || createTeamRequest.UUID == "null" || createTeamRequest.UUID == "00000000-0000-0000-0000-000000000000" {
		createTeamRequest.UUID = uuid.New().String()
	}

	// Check if a team with the same name, email, ID, or UUID already exists
	teamExists, err := ts.IsTeamExist(createTeamRequest.Metadata.Name, createTeamRequest.Content.Email, createTeamRequest.ID, createTeamRequest.UUID)
	ts.logger.Printf("teamExists: %+v\n", teamExists)
	ts.logger.Printf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "team not found" {
			// Expected error
		} else {
			// Handle other errors
			return nil, err
		}
	}

	if teamExists {
		return nil, fmt.Errorf("team with the same name, email, ID, or UUID already exists")
	}

	// // Validate DTO metadata and content
	// if err := createTeamRequest.Validate(); err != nil {
	//	return nil, err
	// }

	// Create a Team instance from the CreateTeamRequest
	team, err := ts.teamRepository.CreateTeam(createTeamRequest)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set(fmt.Sprintf("team:%s", team.UUID), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByName:%s", team.Metadata.Name), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByEmail:%s", team.Content.Email), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByID:%s", team.ID), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("Team: %+v\n", team)
	return team, nil
}

// UpdateTeam updates a team.
func (ts *TeamService) UpdateTeam(uuid string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error) {
	// Check if team exists, and retrieve the team
	team, err := ts.GetTeam(uuid)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}

	// Update the dates values UpdatedAt and UpdatedBy
	updateTeamRequest.Metadata.Dates = &model.CommonDate{
		// CreatedAt:	team.Metadata.Dates.CreatedAt,
		// CreatedBy:	team.Metadata.Dates.CreatedBy,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedBy: updateTeamRequest.Metadata.Dates.UpdatedBy,
		// StartDate:	team.Metadata.Dates.StartDate,
		// EndDate:		team.Metadata.Dates.EndDate,
		// StartedAt:	team.Metadata.Dates.StartedAt,
		// StartedBy:	team.Metadata.Dates.StartedBy,
		// CompletedAt: team.Metadata.Dates.CompletedAt,
		// CompletedBy: team.Metadata.Dates.CompletedBy,
	}

	// // Validate DTO metadata and content
	// if err := updateTeamRequest.Validate(); err != nil {
	//	return nil, err
	// }

	// Update the team
	team, err = ts.teamRepository.UpdateTeam(uuid, updateTeamRequest)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, ErrTeamNotFound
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set(fmt.Sprintf("team:%s", team.UUID), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByName:%s", team.Metadata.Name), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByEmail:%s", team.Content.Email), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByID:%s", team.ID), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("Team: %+v\n", team)
	return team, nil
}

// DeleteTeam deletes a team.
func (ts *TeamService) DeleteTeam(uuid string) (*model.Team, error) {
	team, err := ts.teamRepository.DeleteTeam(uuid)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Delete(fmt.Sprintf("team:%s", uuid))
		ts.cacheManager.Delete(fmt.Sprintf("teamByName:%s", team.Metadata.Name))
		ts.cacheManager.Delete(fmt.Sprintf("teamByEmail:%s", team.Content.Email))
		ts.cacheManager.Delete(fmt.Sprintf("teamByID:%s", team.ID))
	}

	ts.logger.Printf("Team: %+v\n", team)
	return team, nil
}

// GetTeamByID gets a team by ID.
func (ts *TeamService) GetTeamByID(ID string) (*model.Team, error) {
	if ts.isCacheEnabled {
		if team, found := ts.cacheManager.Get(fmt.Sprintf("teamByID:%s", ID)); found {
			return team.(*model.Team), nil
		}
	}

	team, err := ts.teamRepository.GetTeamByID(ID)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set(fmt.Sprintf("teamByID:%s", ID), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("Team: %+v\n", team)
	return team, nil
}

// GetTeamByName gets a team by name.
func (ts *TeamService) GetTeamByName(name string) (*model.Team, error) {
	if ts.isCacheEnabled {
		if team, found := ts.cacheManager.Get(fmt.Sprintf("teamByName:%s", name)); found {
			return team.(*model.Team), nil
		}
	}

	team, err := ts.teamRepository.GetTeamByName(name)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set(fmt.Sprintf("teamByName:%s", name), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("Team: %+v\n", team)
	return team, nil
}

// GetTeamByEmail gets a team by email.
func (ts *TeamService) GetTeamByEmail(email string) (*model.Team, error) {
	if ts.isCacheEnabled {
		if team, found := ts.cacheManager.Get(fmt.Sprintf("teamByEmail:%s", email)); found {
			return team.(*model.Team), nil
		}
	}

	team, err := ts.teamRepository.GetTeamByEmail(email)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set(fmt.Sprintf("teamByEmail:%s", email), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("Team: %+v\n", team)
	return team, nil
}

// ListTeamsMetadata lists teams with metadata.
func (ts *TeamService) ListTeamsMetadata() (*model.ListTeamsMetadataResponse, error) {
	if ts.isCacheEnabled {
		if teams, found := ts.cacheManager.Get("teamsMetadata"); found {
			return teams.(*model.ListTeamsMetadataResponse), nil
		}
	}

	teams, err := ts.teamRepository.ListTeams()
	if err != nil {
		return nil, err
	}

	teamMetadataResponses := make([]*model.TeamMetadataResponse, len(teams))
	for i, team := range teams {
		teamMetadataResponses[i] = &model.TeamMetadataResponse{
			UUID: team.UUID,
			ID:	  team.ID,
			Metadata: &model.TeamMetadata{
				Name: team.Metadata.Name,
				Dates: &model.CommonDate{
					CreatedAt:	 team.Metadata.Dates.CreatedAt,
					CreatedBy:	 team.Metadata.Dates.CreatedBy,
					UpdatedAt:	 team.Metadata.Dates.UpdatedAt,
					UpdatedBy:	 team.Metadata.Dates.UpdatedBy,
					StartDate:	 team.Metadata.Dates.StartDate,
					EndDate:	 team.Metadata.Dates.EndDate,
					StartedAt:	 team.Metadata.Dates.StartedAt,
					StartedBy:	 team.Metadata.Dates.StartedBy,
					CompletedAt: team.Metadata.Dates.CompletedAt,
					CompletedBy: team.Metadata.Dates.CompletedBy,
				},
			},
		}
	}
	listTeamsMetadataResponse := &model.ListTeamsMetadataResponse{
		TeamMetadataResponses: teamMetadataResponses,
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set("listTeamsMetadataResponse", listTeamsMetadataResponse, cache.DefaultExpiration)
	}

	ts.logger.Printf("Teams: %+v\n", listTeamsMetadataResponse)
	return listTeamsMetadataResponse, nil
}

// ListTeamsContent lists teams with content.
func (ts *TeamService) ListTeamsContent() (*model.ListTeamsContentResponse, error) {
	if ts.isCacheEnabled {
		if teams, found := ts.cacheManager.Get("teamsContent"); found {
			return teams.(*model.ListTeamsContentResponse), nil
		}
	}

	teams, err := ts.teamRepository.ListTeams()
	if err != nil {
		return nil, err
	}

	teamContentResponses := make([]*model.TeamContentResponse, len(teams))
	for i, team := range teams {
		teamContentResponses[i] = &model.TeamContentResponse{
			UUID: team.UUID,
			ID:	  team.ID,
			Content: &model.TeamContent{
				Email:		  team.Content.Email,
				ProductOwner: team.Content.ProductOwner,
				ScrumMaster:  team.Content.ScrumMaster,
				Members:	  team.Content.Members,
			},
		}
	}
	listTeamsContentResponse := &model.ListTeamsContentResponse{
		TeamContentResponses: teamContentResponses,
	}

	if ts.isCacheEnabled {
		ts.cacheManager.Set("listTeamsContentResponse", listTeamsContentResponse, cache.DefaultExpiration)
	}

	ts.logger.Printf("Teams: %+v\n", listTeamsContentResponse)
	return listTeamsContentResponse, nil
}

// UpdateTeamMetadata updates a team's metadata.
func (ts *TeamService) UpdateTeamMetadata(uuid string, updateTeamMetadataRequest model.UpdateTeamMetadataRequest) (*model.TeamMetadata, error) {
	newTeamMetadata := model.TeamMetadata{
		Name:  updateTeamMetadataRequest.Metadata.Name,
		Dates: updateTeamMetadataRequest.Metadata.Dates,
	}
	teamMetadata, err := ts.teamRepository.UpdateTeamMetadata(uuid, newTeamMetadata)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		team, _ := ts.teamRepository.GetTeam(uuid)
		team.Metadata = teamMetadata
		ts.cacheManager.Set(fmt.Sprintf("team:%s", uuid), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByName:%s", team.Metadata.Name), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByEmail:%s", team.Content.Email), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByID:%s", team.ID), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("TeamMetadata: %+v\n", teamMetadata)
	return teamMetadata, nil
}

// UpdateTeamContent updates a team's content.
func (ts *TeamService) UpdateTeamContent(uuid string, updateTeamContentRequest model.UpdateTeamContentRequest) (*model.TeamContent, error) {
	newTeamContent := model.TeamContent{
		Email: updateTeamContentRequest.Content.Email,
		ProductOwner: &model.IdUuid{
			ID:	  updateTeamContentRequest.Content.ProductOwner.ID,
			UUID: updateTeamContentRequest.Content.ProductOwner.UUID,
		},
		ScrumMaster: &model.IdUuid{
			ID:	  updateTeamContentRequest.Content.ScrumMaster.ID,
			UUID: updateTeamContentRequest.Content.ScrumMaster.UUID,
		},
		Members: updateTeamContentRequest.Content.Members,
	}
	teamContent, err := ts.teamRepository.UpdateTeamContent(uuid, newTeamContent)
	if err != nil {
		return nil, err
	}

	if ts.isCacheEnabled {
		team, _ := ts.teamRepository.GetTeam(uuid)
		team.Content = teamContent
		ts.cacheManager.Set(fmt.Sprintf("team:%s", uuid), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByName:%s", team.Metadata.Name), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByEmail:%s", team.Content.Email), team, cache.DefaultExpiration)
		ts.cacheManager.Set(fmt.Sprintf("teamByID:%s", team.ID), team, cache.DefaultExpiration)
	}

	ts.logger.Printf("TeamContent: %+v\n", teamContent)
	return teamContent, nil
}

// GetTeamMetadata gets a team's metadata.
func (ts *TeamService) GetTeamMetadata(uuid string) (*model.TeamMetadata, error) {
	teamMetadata, err := ts.teamRepository.GetTeamMetadata(uuid)
	if err != nil {
		return nil, err
	}

	ts.logger.Printf("TeamMetadata: %+v\n", teamMetadata)
	return teamMetadata, nil
}

// GetTeamContent gets a team's content.
func (ts *TeamService) GetTeamContent(uuid string) (*model.TeamContent, error) {
	teamContent, err := ts.teamRepository.GetTeamContent(uuid)
	if err != nil {
		return nil, err
	}

	ts.logger.Printf("TeamContent: %+v\n", teamContent)
	return teamContent, nil
}

// IsTeamExist checks if a team with the specified attributes exists.
func (ts *TeamService) IsTeamExist(name, email, ID, UUID string) (bool, error) {
	if !ts.isCacheEnabled {
		return ts.checkTeamExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("teamExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if teamExists, found := ts.cacheManager.Get(cacheKey); found {
		return teamExists.(bool), nil
	}

	teamExists, err := ts.checkTeamExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	ts.cacheManager.Set(cacheKey, teamExists, cache.DefaultExpiration)
	return teamExists, nil
}

// IsNoTeamExist checks if no team with the specified attributes exists.
func (ts *TeamService) IsNoTeamExist(name, email, ID, UUID string) (bool, error) {
	if !ts.isCacheEnabled {
		return ts.checkTeamExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("noTeamExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if teamExists, found := ts.cacheManager.Get(cacheKey); found {
		return teamExists.(bool), nil
	}

	teamExists, err := ts.checkTeamExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	ts.cacheManager.Set(cacheKey, teamExists, cache.DefaultExpiration)
	return teamExists, nil
}

// IsExactlyOneTeamExist checks if exactly one team with the specified attributes exists.
func (ts *TeamService) IsExactlyOneTeamExist(name, email, ID, UUID string) (bool, error) {
	if !ts.isCacheEnabled {
		return ts.checkTeamExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("exactlyOneTeamExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if teamExists, found := ts.cacheManager.Get(cacheKey); found {
		return teamExists.(bool), nil
	}

	teamExists, err := ts.checkTeamExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	ts.cacheManager.Set(cacheKey, teamExists, cache.DefaultExpiration)
	return teamExists, nil
}

// IsAtLeastOneTeamExist checks if at least one team with the specified attributes exists.
func (ts *TeamService) IsAtLeastOneTeamExist(name, email, ID, UUID string) (bool, error) {
	if !ts.isCacheEnabled {
		return ts.checkTeamExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("atLeastOneTeamExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if teamExists, found := ts.cacheManager.Get(cacheKey); found {
		return teamExists.(bool), nil
	}

	teamExists, err := ts.checkTeamExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	ts.cacheManager.Set(cacheKey, teamExists, cache.DefaultExpiration)
	return teamExists, nil
}

// checkTeamExistence checks if a team with the specified attributes exists.
func (ts *TeamService) checkTeamExistence(name, email, ID, UUID string) (bool, error) {
	teamByName, err := ts.teamRepository.GetTeamByName(name)
	ts.logger.Debugf("name: %+v\n", name)
	ts.logger.Debugf("teamByName: %+v\n", teamByName)
	ts.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "team not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if teamByName != nil {
		return true, nil
	}

	teamByEmail, err := ts.teamRepository.GetTeamByEmail(email)
	ts.logger.Debugf("email: %+v\n", email)
	ts.logger.Debugf("teamByEmail: %+v\n", teamByEmail)
	ts.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "team not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if teamByEmail != nil {
		return true, nil
	}

	teamByID, err := ts.teamRepository.GetTeamByID(ID)
	ts.logger.Debugf("ID: %+v\n", ID)
	ts.logger.Debugf("teamByID: %+v\n", teamByID)
	ts.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "team not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if teamByID != nil {
		return true, nil
	}

	teamByUUID, err := ts.teamRepository.GetTeam(UUID)
	ts.logger.Debugf("UUID: %+v\n", UUID)
	ts.logger.Debugf("teamByUUID: %+v\n", teamByUUID)
	ts.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "team not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if teamByUUID != nil {
		return true, nil
	}

	return false, nil
}
