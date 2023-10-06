// Path: internal/repository/team/memory_team_repository.go
// DESC: This is the memory implementation of the team repository.
package postgres_sql

import (
	"context"
	"database/sql"

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

// PostgresTeamRepository implements TeamRepositoryInterface for PostgreSQL.
type PostgresTeamRepository struct {
	ctx	   context.Context
	logger *logrus.Logger
	db	   *sql.DB
}

// NewPostgresTeamRepository creates a new instance of PostgresTeamRepository.
func NewPostgresTeamRepository(ctx context.Context, logger *logrus.Logger, db *sql.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{
		ctx:	ctx,
		logger: logger,
		db:		db,
	}
}

// ListTeamIdsAndUUIDs returns a list of team IDs and UUIDs.
func (repository *PostgresTeamRepository) ListTeamIdsAndUUIDs() ([]model.IdUuid, error) {
	// Implement your PostgreSQL query here to fetch team IDs and UUIDs.
	return nil, nil
	//	query := `
	//		SELECT id, uuid
	//		FROM teams
	//	`

	//	rows, err := repository.db.Query(query)
	//	if err != nil {
	//		return nil, err
	//	}

	//	defer rows.Close()

	//	var idUuids []model.IdUuid

	//	for rows.Next() {
	//		var idUuid model.IdUuid

	//		err := rows.Scan(&idUuid.ID, &idUuid.UUID)
	//		if err != nil {
	//			return nil, err
	//		}

	//		idUuids = append(idUuids, idUuid)
	//	}

	//	if err := rows.Err(); err != nil {
	//		return nil, err
	//	}

	//	return idUuids, nil
	// }

	// // ListTeams returns a list of teams.
	// func (repository *PostgresTeamRepository) ListTeams() ([]model.Team, error) {
	//	// Implement your PostgreSQL query here to fetch teams.
	//	query := `
	//		SELECT t.id, t.uuid,
	//			   tm.name,
	//			   tm.created_at AS team_created_at,
	//			   tm.created_by AS team_created_by,
	//			   tm.updated_at AS team_updated_at,
	//			   tm.updated_by AS team_updated_by,
	//			   tm.started_at AS team_started_at,
	//			   tm.started_by AS team_started_by,
	//			   tm.start_date AS team_start_date,
	//			   tm.end_date AS team_end_date,
	//			   tm.completed_at AS team_completed_at,
	//			   tm.completed_by AS team_completed_by,
	//			   tc.email,
	//			   json_agg(json_build_object('ID', m.id, 'UUID', m.uuid)) AS members,
	//			   json_build_object('ID', po.id, 'UUID', po.uuid) AS product_owner,
	//			   json_build_object('ID', sm.id, 'UUID', sm.uuid) AS scrum_master
	//		FROM teams t
	//		JOIN team_metadata tm ON t.id = tm.team_id
	//		JOIN team_content tc ON t.id = tc.team_id
	//		LEFT JOIN members m ON t.id = m.team_id
	//		LEFT JOIN users po ON t.product_owner_id = po.id
	//		LEFT JOIN users sm ON t.scrum_master_id = sm.id
	//		GROUP BY t.id, t.uuid, tm.name, tm.created_at, tm.updated_at, tc.email, po.id, po.uuid, sm.id, sm.uuid
	//	`

	//	rows, err := repository.db.Query(query)
	//	if err != nil {
	//		return nil, err
	//	}
	//	defer rows.Close()

	//	var teams []model.Team

	//	for rows.Next() {
	//		var team model.Team
	//		var membersJSON []byte
	//		var productOwnerJSON []byte
	//		var scrumMasterJSON []byte

	//		err := rows.Scan(
	//			&team.ID, &team.UUID,
	//			&team.Metadata.Name,
	//			&team.Metadata.Dates.CreatedAt,
	//			&team.Metadata.Dates.CompletedBy,
	//			&team.Metadata.Dates.UpdatedAt,
	//			&team.Metadata.Dates.UpdatedBy,
	//			&team.Metadata.Dates.StartedAt,
	//			&team.Metadata.Dates.StartedBy,
	//			&team.Metadata.Dates.StartDate,
	//			&team.Metadata.Dates.EndDate,
	//			&team.Metadata.Dates.CompletedAt,
	//			&team.Metadata.Dates.CompletedBy,
	//			&team.Content.Email,
	//			&membersJSON,
	//			&productOwnerJSON,
	//			&scrumMasterJSON,
	//		)
	//		if err != nil {
	//			return nil, err
	//		}

	//		// Unmarshal JSON data into slices or objects as needed.
	//		if err := json.Unmarshal(membersJSON, &team.Content.Members); err != nil {
	//			return nil, err
	//		}
	//		if err := json.Unmarshal(productOwnerJSON, &team.Content.ProductOwner); err != nil {
	//			return nil, err
	//		}
	//		if err := json.Unmarshal(scrumMasterJSON, &team.Content.ScrumMaster); err != nil {
	//			return nil, err
	//		}

	//		teams = append(teams, team)
	//	}

	//	if err := rows.Err(); err != nil {
	//		return nil, err
	//	}

	// return teams, nil
}

// GetTeam retrieves a team by UUID.
func (repository *PostgresTeamRepository) GetTeam(UUID string) (*model.Team, error) {
	// Implement your PostgreSQL query here to fetch a team by UUID.
	return nil, nil
	// query := `
	//	SELECT t.id, t.uuid,
	//	tm.name,
	//	tm.created_at AS team_created_at,
	//	tm.created_by AS team_created_by,
	//	tm.updated_at AS team_updated_at,
	//	tm.updated_by AS team_updated_by,
	//	tm.started_at AS team_started_at,
	//	tm.started_by AS team_started_by,
	//	tm.start_date AS team_start_date,
	//	tm.end_date AS team_end_date,
	//	tm.completed_at AS team_completed_at,
	//	tm.completed_by AS team_completed_by,
	//	tc.email,
	//	json_agg(json_build_object('ID', m.id, 'UUID', m.uuid)) AS members,
	//	json_build_object('ID', po.id, 'UUID', po.uuid) AS product_owner,
	//	json_build_object('ID', sm.id, 'UUID', sm.uuid) AS scrum_master
	//	FROM teams t
	//	JOIN team_metadata tm ON t.id = tm.team_id
	//	JOIN team_content tc ON t.id = tc.team_id
	//	LEFT JOIN members m ON t.id = m.team_id
	//	LEFT JOIN users po ON t.product_owner_id = po.id
	//	LEFT JOIN users sm ON t.scrum_master_id = sm.id
	//	WHERE t.uuid = $1
	//	GROUP BY t.id, t.uuid, tm.name, tm.created_at, tm.updated_at, tc.email, po.id, po.uuid, sm.id, sm.uuid
	// `

	// row := repository.db.QueryRow(query, UUID)

	// var team model.Team
	// var membersJSON []byte
	// var productOwnerJSON []byte
	// var scrumMasterJSON []byte

	// err := row.Scan(
	//	&team.ID, &team.UUID,
	//	&team.Metadata.Name,
	//	&team.Metadata.Dates.CreatedAt,
	//	&team.Metadata.Dates.CompletedBy,
	//	&team.Metadata.Dates.UpdatedAt,
	//	&team.Metadata.Dates.UpdatedBy,
	//	&team.Metadata.Dates.StartedAt,
	//	&team.Metadata.Dates.StartedBy,
	//	&team.Metadata.Dates.StartDate,
	//	&team.Metadata.Dates.EndDate,
	//	&team.Metadata.Dates.CompletedAt,
	//	&team.Metadata.Dates.CompletedBy,
	//	&team.Content.Email,
	//	&membersJSON,
	//	&productOwnerJSON,
	//	&scrumMasterJSON,
	// )
	// if err != nil {
	//	return nil, err
	// }

	// // Unmarshal JSON data into slices or objects as needed.
	// if err := json.Unmarshal(membersJSON, &team.Content.Members); err != nil {
	//	return nil, err
	// }

	// if err := json.Unmarshal(productOwnerJSON, &team.Content.ProductOwner); err != nil {
	//	return nil, err
	// }

	// if err := json.Unmarshal(scrumMasterJSON, &team.Content.ScrumMaster); err != nil {
	//	return nil, err
	// }

	// return &team, nil
}

// CreateTeam creates a new team.
func (repository *PostgresTeamRepository) CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error) {
	// Implement your PostgreSQL query here to create a new team.
	return nil, nil
}

// UpdateTeam updates a team by UUID.
func (repository *PostgresTeamRepository) UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error) {
	// Implement your PostgreSQL query here to update a team by UUID.
	return nil, nil
}

// DeleteTeam deletes a team by UUID.
func (repository *PostgresTeamRepository) DeleteTeam(UUID string) (*model.Team, error) {
	// Implement your PostgreSQL query here to delete a team by UUID.
	return nil, nil
}

// GetTeamByID retrieves a team by ID.
func (repository *PostgresTeamRepository) GetTeamByID(ID string) (*model.Team, error) {
	// Implement your PostgreSQL query here to fetch a team by ID.
	return nil, nil
}

// GetTeamByName retrieves a team by name.
func (repository *PostgresTeamRepository) GetTeamByName(name string) (*model.Team, error) {
	// Implement your PostgreSQL query here to fetch a team by name.
	return nil, nil
}

// GetTeamByEmail retrieves a team by email.
func (repository *PostgresTeamRepository) GetTeamByEmail(email string) (*model.Team, error) {
	// Implement your PostgreSQL query here to fetch a team by email.
	return nil, nil
}

// UpdateTeamMetadata updates team metadata by UUID.
func (repository *PostgresTeamRepository) UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error) {
	// Implement your PostgreSQL query here to update team metadata by UUID.
	return nil, nil
}

// UpdateTeamContent updates team content by UUID.
func (repository *PostgresTeamRepository) UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error) {
	// Implement your PostgreSQL query here to update team content by UUID.
	return nil, nil
}

// GetTeamMetadata retrieves team metadata by UUID.
func (repository *PostgresTeamRepository) GetTeamMetadata(UUID string) (*model.TeamMetadata, error) {
	// Implement your PostgreSQL query here to fetch team metadata by UUID.
	return nil, nil
}

// GetTeamContent retrieves team content by UUID.
func (repository *PostgresTeamRepository) GetTeamContent(UUID string) (*model.TeamContent, error) {
	// Implement your PostgreSQL query here to fetch team content by UUID.
	return nil, nil
}
