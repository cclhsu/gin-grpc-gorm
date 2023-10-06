// Path: internal/repository/team/memory_team_repository.go
// DESC: This is the memory implementation of the team repository.
package mongodb

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

type MongoDBTeamRepository struct {
	ctx		   context.Context
	logger	   *logrus.Logger
	db		   *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBTeamRepository(ctx context.Context, logger *logrus.Logger, mongoDBClient *mongo.Client) *MongoDBTeamRepository {
	mongoDBName := os.Getenv("MONGO_DB")
	return &MongoDBTeamRepository{
		ctx:		ctx,
		logger:		logger,
		db:			mongoDBClient,
		collection: mongoDBClient.Database(mongoDBName).Collection("teams"),
	}
}

func (repository *MongoDBTeamRepository) ListTeamIdsAndUUIDs() ([]model.IdUuid, error) {
	teamIdUuids := []model.IdUuid{}
	repository.logger.Debug("ListTeamIdsAndUUIDs")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(repository.ctx, 5*time.Second)
	defer cancel()

	// Perform the find operation
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and decode results into teams
	for cursor.Next(ctx) {
		var team model.Team
		if err := cursor.Decode(&team); err != nil {
			return nil, err
		}
		teamIdUuids = append(teamIdUuids, model.IdUuid{ID: team.ID, UUID: team.UUID})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return teamIdUuids, nil
}

func (repository *MongoDBTeamRepository) ListTeams() ([]model.Team, error) {
	var teams []model.Team
	repository.logger.Debug("ListTeams")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(repository.ctx, 5*time.Second)
	defer cancel()

	// Perform the find operation
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode results into teams
	for cursor.Next(ctx) {
		var team model.Team
		if err := cursor.Decode(&team); err != nil {
			return nil, err
		}
		teams = append(teams, team)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return teams, nil
}

// func (repository *MongoDBTeamRepository) ListTeamsMetadata() ([]model.TeamMetadata, error) {
//	return nil, nil
// }

// func (repository *MongoDBTeamRepository) ListTeamsContent() ([]model.TeamContent, error) {
//	return nil, nil
// }

func (repository *MongoDBTeamRepository) GetTeam(UUID string) (*model.Team, error) {
	// Define a filter to find the team by UUID
	filter := bson.M{"uuid": UUID}

	// Create an empty Team struct to store the result
	var team model.Team

	// Find the team using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&team)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}

	return &team, nil
}

func (repository *MongoDBTeamRepository) CreateTeam(createTeamRequest model.CreateTeamRequest) (*model.Team, error) {
	// Create a Team instance from the CreateTeamRequest
	team := model.Team{
		UUID:	  uuid.New().String(),
		ID:		  createTeamRequest.ID,
		Metadata: createTeamRequest.Metadata,
		Content:  createTeamRequest.Content,
	}

	// Insert the team into the MongoDB collection
	insertResult, err := repository.collection.InsertOne(repository.ctx, team)
	if err != nil {
		return nil, err
	}
	if insertResult == nil {
		return nil, fmt.Errorf("failed to insert team")
	}

	// // Convert the _id to a string
	// insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	// if !ok {
	//	return nil, fmt.Errorf("failed to convert inserted ID to string")
	// }

	return &team, nil
}

func (repository *MongoDBTeamRepository) UpdateTeam(UUID string, updateTeamRequest model.UpdateTeamRequest) (*model.Team, error) {
	if !govalidator.IsUUIDv4(UUID) {
		return nil, fmt.Errorf("invalid UUID")
	}
	team, err := repository.GetTeam(UUID)
	if err != nil {
		return nil, err
	}

	// update team if updateTeamRequest has new value
	if updateTeamRequest.Metadata.Name != "" {
		team.Metadata.Name = updateTeamRequest.Metadata.Name
	}
	if updateTeamRequest.Metadata.Dates.UpdatedAt != "" {
		team.Metadata.Dates.UpdatedAt = updateTeamRequest.Metadata.Dates.UpdatedAt
	}
	if updateTeamRequest.Metadata.Dates.UpdatedBy != "" {
		team.Metadata.Dates.UpdatedBy = updateTeamRequest.Metadata.Dates.UpdatedBy
	}
	if updateTeamRequest.Metadata.Dates.StartDate != "" {
		team.Metadata.Dates.StartDate = updateTeamRequest.Metadata.Dates.StartDate
	}
	if updateTeamRequest.Metadata.Dates.EndDate != "" {
		team.Metadata.Dates.EndDate = updateTeamRequest.Metadata.Dates.EndDate
	}
	if updateTeamRequest.Metadata.Dates.StartedAt != "" {
		team.Metadata.Dates.StartedAt = updateTeamRequest.Metadata.Dates.StartedAt
	}
	if updateTeamRequest.Metadata.Dates.StartedBy != "" {
		team.Metadata.Dates.StartedBy = updateTeamRequest.Metadata.Dates.StartedBy
	}
	if updateTeamRequest.Metadata.Dates.CompletedAt != "" {
		team.Metadata.Dates.CompletedAt = updateTeamRequest.Metadata.Dates.CompletedAt
	}
	if updateTeamRequest.Metadata.Dates.CompletedBy != "" {
		team.Metadata.Dates.CompletedBy = updateTeamRequest.Metadata.Dates.CompletedBy
	}
	if updateTeamRequest.Content.Email != "" {
		team.Content.Email = updateTeamRequest.Content.Email
	}
	if updateTeamRequest.Content.ProductOwner != nil {
		team.Content.ProductOwner = updateTeamRequest.Content.ProductOwner
	}
	if updateTeamRequest.Content.ScrumMaster != nil {
		team.Content.ScrumMaster = updateTeamRequest.Content.ScrumMaster
	}
	if updateTeamRequest.Content.Members != nil {
		team.Content.Members = updateTeamRequest.Content.Members
	}

	filter := bson.M{"uuid": UUID}
	update := bson.M{"$set": team}
	_, err = repository.collection.UpdateOne(repository.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return team, nil
}

func (repository *MongoDBTeamRepository) DeleteTeam(UUID string) (*model.Team, error) {
	team := model.Team{
		UUID: UUID,
	}

	filter := bson.M{"uuid": UUID}
	_, err := repository.collection.DeleteOne(repository.ctx, filter)
	if err != nil {
		return nil, err
	}

	return &team, nil
}

func (repository *MongoDBTeamRepository) GetTeamByID(ID string) (*model.Team, error) {
	// Define a filter to find the team by ID
	filter := bson.M{"id": ID}

	// Create an empty Team struct to store the result
	var team model.Team

	// Find the team using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&team)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}

	return &team, nil
}

func (repository *MongoDBTeamRepository) GetTeamByName(name string) (*model.Team, error) {
	// Define a filter to find the team by name
	filter := bson.M{"metadata.name": name}

	// Create an empty Team struct to store the result
	var team model.Team

	// Find the team using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&team)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}

	return &team, nil
}

func (repository *MongoDBTeamRepository) GetTeamByEmail(email string) (*model.Team, error) {
	// Define a filter to find the team by email
	filter := bson.M{"content.email": email}

	// Create an empty Team struct to store the result
	var team model.Team

	// Find the team using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&team)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}

	return &team, nil
}

func (repository *MongoDBTeamRepository) UpdateTeamMetadata(UUID string, updatedMetadata model.TeamMetadata) (*model.TeamMetadata, error) {
	// get team by UUID
	team, err := repository.GetTeam(UUID)
	if err != nil {
		return nil, err
	}

	// update team metadata
	team.Metadata = &updatedMetadata

	filter := bson.M{"uuid": UUID}
	update := bson.M{"$set": team}
	_, err = repository.collection.UpdateOne(repository.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return team.Metadata, nil
}

func (repository *MongoDBTeamRepository) UpdateTeamContent(UUID string, updatedContent model.TeamContent) (*model.TeamContent, error) {
	// get team by UUID
	team, err := repository.GetTeam(UUID)
	if err != nil {
		return nil, err
	}

	// update team content
	team.Content = &updatedContent

	filter := bson.M{"uuid": UUID}
	update := bson.M{"$set": team}
	_, err = repository.collection.UpdateOne(repository.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return team.Content, nil
}

func (repository *MongoDBTeamRepository) GetTeamMetadata(UUID string) (*model.TeamMetadata, error) {
	// Define a filter to find the team by UUID
	filter := bson.M{"uuid": UUID}

	// Create an empty Team struct to store the result
	var team model.Team

	// Find the team using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&team)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}

	return team.Metadata, nil
}

func (repository *MongoDBTeamRepository) GetTeamContent(UUID string) (*model.TeamContent, error) {
	// Define a filter to find the team by UUID
	filter := bson.M{"uuid": UUID}

	// Create an empty Team struct to store the result
	var team model.Team

	// Find the team using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&team)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("team not found")
		}
		return nil, err
	}

	return team.Content, nil
}
