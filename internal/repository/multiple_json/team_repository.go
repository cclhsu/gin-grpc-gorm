// Path: internal/repository/team/multiple_json_team_repository.go
// DESC: This is the multiple json implementation of the team repository.
package multiple_json

// import (
//	"github.com/cclhsu/gin-grpc-gorm/internal/model"
// )

// // TeamRepositoryInterface defines the interface for the team repository.
// type TeamRepositoryInterface interface {
//	ListTeamIdsAndUUIDs() ([]model.IdUuid, error)
//	ListTeams() ([]model.Team, error)
//	// ListTeamsMetadata() ([]model.TeamMetadata, error)
//	// ListTeamsContent() ([]model.TeamContent, error)
//	GetTeam(UUID string) (*model.Team, error)
//	CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error)
//	UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error)
//	DeleteTeam(UUID string) (*model.Team, error)
//	GetTeamByID(ID string) (*model.Team, error)
//	GetTeamByName(name string) (*model.Team, error)
//	GetTeamByEmail(email string) (*model.Team, error)
//	UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error)
//	UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error)
//	GetTeamMetadata(UUID string) (*model.TeamMetadata, error)
//	GetTeamContent(UUID string) (*model.TeamContent, error)
// }

// type MultipleJsonTeamRepository struct {
// }

// func NewMultipleJsonTeamRepository() *MultipleJsonTeamRepository {
//	return &MultipleJsonTeamRepository{}
// }

// func (repository *MultipleJsonTeamRepository) ListTeamIdsAndUUIDs() ([]model.IdUuid, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) ListTeams() ([]model.Team, error) {
//	return nil, nil
// }

// // func (repository *MultipleJsonTeamRepository) ListTeamsMetadata() ([]model.TeamMetadata, error) {
// //	return nil, nil
// // }

// // func (repository *MultipleJsonTeamRepository) ListTeamsContent() ([]model.TeamContent, error) {
// //	return nil, nil
// // }

// func (repository *MultipleJsonTeamRepository) GetTeam(UUID string) (*model.Team, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) DeleteTeam(UUID string) (*model.Team, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) GetTeamByID(ID string) (*model.Team, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) GetTeamByName(name string) (*model.Team, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) GetTeamByEmail(email string) (*model.Team, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) GetTeamMetadata(UUID string) (*model.TeamMetadata, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonTeamRepository) GetTeamContent(UUID string) (*model.TeamContent, error) {
//	return nil, nil
// }

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// TeamRepositoryInterface defines the interface for the team repository.
type TeamRepositoryInterface interface {
	ListTeamIdsAndUUIDs() ([]model.IdUuid, error)
	ListTeams() ([]model.Team, error)
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

// MultipleJsonTeamRepository implements TeamRepositoryInterface using multiple JSON files.
type MultipleJsonTeamRepository struct {
	ctx			  context.Context
	logger		  *logrus.Logger
	dirPath		  string
	fileExtension string
	lock		  sync.RWMutex
}

// NewMultipleJsonTeamRepository creates a new instance of MultipleJsonTeamRepository.
func NewMultipleJsonTeamRepository(ctx context.Context, logger *logrus.Logger, dirPath string) TeamRepositoryInterface {
	return &MultipleJsonTeamRepository{
		ctx:		   ctx,
		logger:		   logger,
		dirPath:	   dirPath,
		fileExtension: ".json",
	}
}

func (r *MultipleJsonTeamRepository) listFiles() ([]string, error) {
	files, err := fileutils.ListFiles(r.dirPath, r.fileExtension)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// ListTeamIdsAndUUIDs returns a list of team IDs and UUIDs.
func (r *MultipleJsonTeamRepository) ListTeamIdsAndUUIDs() ([]model.IdUuid, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	teamIdsAndUUIDs := make([]model.IdUuid, 0)

	for _, filePath := range files {
		teamID := strings.TrimSuffix(filepath.Base(filePath), r.fileExtension)
		teamUUID, err := uuid.Parse(teamID)
		if err != nil {
			return nil, err
		}

		teamIdsAndUUIDs = append(teamIdsAndUUIDs, model.IdUuid{
			ID:	  teamID,
			UUID: teamUUID.String(),
		})
	}

	return teamIdsAndUUIDs, nil
}

// ListTeams returns a list of teams.
func (r *MultipleJsonTeamRepository) ListTeams() ([]model.Team, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	teams := make([]model.Team, 0)

	for _, filePath := range files {
		team, err := r.readTeamFromFile(filePath)
		if err != nil {
			return nil, err
		}

		teams = append(teams, team)
	}

	return teams, nil
}

// GetTeam retrieves a team by UUID.
func (r *MultipleJsonTeamRepository) GetTeam(UUID string) (*model.Team, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.readTeamFromFile(filePath)
}

func (r *MultipleJsonTeamRepository) readTeamFromFile(filePath string) (*model.Team, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var team model.Team
	if err := json.Unmarshal(fileData, &team); err != nil {
		return nil, err
	}

	return &team, nil
}

// CreateTeam creates a new team.
func (r *MultipleJsonTeamRepository) CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error) {
	newUUID := uuid.New()
	newTeam := model.NewTeam(newUUID, createTeamRequest.Metadata, createTeamRequest.Content)
	filePath := filepath.Join(r.dirPath, newUUID.String()+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.writeTeamToFile(filePath, newTeam); err != nil {
		return nil, err
	}

	return newTeam, nil
}

func (r *MultipleJsonTeamRepository) writeTeamToFile(filePath string, team *model.Team) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	teamJSON, err := json.Marshal(team)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filePath, teamJSON, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// UpdateTeam updates a team by UUID.
func (r *MultipleJsonTeamRepository) UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	existingTeam, err := r.readTeamFromFile(filePath)
	if err != nil {
		return nil, err
	}

	existingTeam.UpdateMetadata(updateTeamRequest.Metadata)
	existingTeam.UpdateContent(updateTeamRequest.Content)

	if err := r.writeTeamToFile(filePath, existingTeam); err != nil {
		return nil, err
	}

	return existingTeam, nil
}

// DeleteTeam deletes a team by UUID.
func (r *MultipleJsonTeamRepository) DeleteTeam(UUID string) (*model.Team, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	deletedTeam, err := r.readTeamFromFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := os.Remove(filePath); err != nil {
		return nil, err
	}

	return deletedTeam, nil
}

// GetTeamByID retrieves a team by ID.
func (r *MultipleJsonTeamRepository) GetTeamByID(ID string) (*model.Team, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, filePath := range files {
		team, err := r.readTeamFromFile(filePath)
		if err != nil {
			return nil, err
		}

		if team.ID == ID {
			return team, nil
		}
	}

	return nil, fmt.Errorf("Team with ID %s not found", ID)
}

// GetTeamByName retrieves a team by name.
func (r *MultipleJsonTeamRepository) GetTeamByName(name string) (*model.Team, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, filePath := range files {
		team, err := r.readTeamFromFile(filePath)
		if err != nil {
			return nil, err
		}

		if team.Metadata.Name == name {
			return team, nil
		}
	}

	return nil, fmt.Errorf("Team with name %s not found", name)
}

// GetTeamByEmail retrieves a team by email.
func (r *MultipleJsonTeamRepository) GetTeamByEmail(email string) (*model.Team, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, filePath := range files {
		team, err := r.readTeamFromFile(filePath)
		if err != nil {
			return nil, err
		}

		if team.Content.Email == email {
			return team, nil
		}
	}

	return nil, fmt.Errorf("Team with email %s not found", email)
}

// UpdateTeamMetadata updates team metadata by UUID.
func (r *MultipleJsonTeamRepository) UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	existingTeam, err := r.readTeamFromFile(filePath)
	if err != nil {
		return nil, err
	}

	existingTeam.UpdateMetadata(updatedMetadata)

	if err := r.writeTeamToFile(filePath, existingTeam); err != nil {
		return nil, err
	}

	return existingTeam.Metadata, nil
}

// UpdateTeamContent updates team content by UUID.
func (r *MultipleJsonTeamRepository) UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	existingTeam, err := r.readTeamFromFile(filePath)
	if err != nil {
		return nil, err
	}

	existingTeam.UpdateContent(updatedContent)

	if err := r.writeTeamToFile(filePath, existingTeam); err != nil {
		return nil, err
	}

	return existingTeam.Content, nil
}

// GetTeamMetadata retrieves team metadata by UUID.
func (r *MultipleJsonTeamRepository) GetTeamMetadata(UUID string) (*model.TeamMetadata, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.RLock()
	defer r.lock.RUnlock()

	existingTeam, err := r.readTeamFromFile(filePath)
	if err != nil {
		return nil, err
	}

	return existingTeam.Metadata, nil
}

// GetTeamContent retrieves team content by UUID.
func (r *MultipleJsonTeamRepository) GetTeamContent(UUID string) (*model.TeamContent, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.RLock()
	defer r.lock.RUnlock()

	existingTeam, err := r.readTeamFromFile(filePath)
	if err != nil {
		return nil, err
	}

	return existingTeam.Content, nil
}
