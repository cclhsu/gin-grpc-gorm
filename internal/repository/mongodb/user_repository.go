// Path: internal/repository/user/memory_user_repository.go
// DESC: This is the memory implementation of the user repository.
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

// UserRepositoryInterface defines the interface for the user repository.
type UserRepositoryInterface interface {
	ListUserIdsAndUUIDs() ([]model.IdUuid, error)
	ListUsers() ([]model.User, error)
	// ListUsersMetadata() ([]model.UserMetadata, error)
	// ListUsersContent() ([]model.UserContent, error)
	GetUser(UUID string) (*model.User, error)
	CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error)
	UpdateUser(UUID string, updateUserRequest model.UpdateUserRequest) (*model.User, error)
	DeleteUser(UUID string) (*model.User, error)
	GetUserByID(ID string) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	UpdateUserMetadata(UUID string, updatedMetadata model.UserMetadata) (*model.UserMetadata, error)
	UpdateUserContent(UUID string, updatedContent model.UserContent) (*model.UserContent, error)
	GetUserMetadata(UUID string) (*model.UserMetadata, error)
	GetUserContent(UUID string) (*model.UserContent, error)
}

type MongoDBUserRepository struct {
	ctx		   context.Context
	logger	   *logrus.Logger
	db		   *mongo.Client
	collection *mongo.Collection
}

func NewMongoDBUserRepository(ctx context.Context, logger *logrus.Logger, mongoDBClient *mongo.Client) *MongoDBUserRepository {
	mongoDBName := os.Getenv("MONGO_DB")
	return &MongoDBUserRepository{
		ctx:		ctx,
		logger:		logger,
		db:			mongoDBClient,
		collection: mongoDBClient.Database(mongoDBName).Collection("users"),
	}
}

func (repository *MongoDBUserRepository) ListUserIdsAndUUIDs() ([]model.IdUuid, error) {
	userIdUuids := []model.IdUuid{}
	repository.logger.Debug("ListUserIdsAndUUIDs")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(repository.ctx, 5*time.Second)
	defer cancel()

	// Perform the find operation
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	// Iterate over the cursor and decode results into users
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		userIdUuids = append(userIdUuids, model.IdUuid{ID: user.ID, UUID: user.UUID})
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return userIdUuids, nil
}

func (repository *MongoDBUserRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	repository.logger.Debug("ListUsers")

	// Create a context with a timeout
	ctx, cancel := context.WithTimeout(repository.ctx, 5*time.Second)
	defer cancel()

	// Perform the find operation
	cursor, err := repository.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	// Iterate over the cursor and decode results into users
	for cursor.Next(ctx) {
		var user model.User
		if err := cursor.Decode(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

// func (repository *MongoDBUserRepository) ListUsersMetadata() ([]model.UserMetadata, error) {
//	return nil, nil
// }

// func (repository *MongoDBUserRepository) ListUsersContent() ([]model.UserContent, error) {
//	return nil, nil
// }

func (repository *MongoDBUserRepository) GetUser(UUID string) (*model.User, error) {
	// Define a filter to find the user by UUID
	filter := bson.M{"uuid": UUID}

	// Create an empty User struct to store the result
	var user model.User

	// Find the user using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&user)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (repository *MongoDBUserRepository) CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error) {
	// Create a User instance from the CreateUserRequest
	user := model.User{
		UUID:	  uuid.New().String(),
		ID:		  createUserRequest.ID,
		Metadata: createUserRequest.Metadata,
		Content:  createUserRequest.Content,
	}

	// Insert the user into the MongoDB collection
	insertResult, err := repository.collection.InsertOne(repository.ctx, user)
	if err != nil {
		return nil, err
	}
	if insertResult == nil {
		return nil, fmt.Errorf("failed to insert user")
	}

	// // Convert the _id to a string
	// insertedID, ok := insertResult.InsertedID.(primitive.ObjectID)
	// if !ok {
	//	return nil, fmt.Errorf("failed to convert inserted ID to string")
	// }

	return &user, nil
}

func (repository *MongoDBUserRepository) UpdateUser(UUID string, updateUserRequest model.UpdateUserRequest) (*model.User, error) {
	if !govalidator.IsUUIDv4(UUID) {
		return nil, fmt.Errorf("invalid UUID")
	}
	user, err := repository.GetUser(UUID)
	if err != nil {
		return nil, err
	}

	// update user if updateUserRequest has new value
	if updateUserRequest.Metadata.Name != "" {
		user.Metadata.Name = updateUserRequest.Metadata.Name
	}
	if updateUserRequest.Metadata.Dates.UpdatedAt != "" {
		user.Metadata.Dates.UpdatedAt = updateUserRequest.Metadata.Dates.UpdatedAt
	}
	if updateUserRequest.Metadata.Dates.UpdatedBy != "" {
		user.Metadata.Dates.UpdatedBy = updateUserRequest.Metadata.Dates.UpdatedBy
	}
	if updateUserRequest.Metadata.Dates.StartDate != "" {
		user.Metadata.Dates.StartDate = updateUserRequest.Metadata.Dates.StartDate
	}
	if updateUserRequest.Metadata.Dates.EndDate != "" {
		user.Metadata.Dates.EndDate = updateUserRequest.Metadata.Dates.EndDate
	}
	if updateUserRequest.Metadata.Dates.StartedAt != "" {
		user.Metadata.Dates.StartedAt = updateUserRequest.Metadata.Dates.StartedAt
	}
	if updateUserRequest.Metadata.Dates.StartedBy != "" {
		user.Metadata.Dates.StartedBy = updateUserRequest.Metadata.Dates.StartedBy
	}
	if updateUserRequest.Metadata.Dates.CompletedAt != "" {
		user.Metadata.Dates.CompletedAt = updateUserRequest.Metadata.Dates.CompletedAt
	}
	if updateUserRequest.Metadata.Dates.CompletedBy != "" {
		user.Metadata.Dates.CompletedBy = updateUserRequest.Metadata.Dates.CompletedBy
	}
	if updateUserRequest.Content.Email != "" {
		user.Content.Email = updateUserRequest.Content.Email
	}
	if updateUserRequest.Content.Phone != "" {
		user.Content.Phone = updateUserRequest.Content.Phone
	}
	if updateUserRequest.Content.LastName != "" {
		user.Content.LastName = updateUserRequest.Content.LastName
	}
	if updateUserRequest.Content.FirstName != "" {
		user.Content.FirstName = updateUserRequest.Content.FirstName
	}
	if updateUserRequest.Content.ProjectRoles != nil {
		user.Content.ProjectRoles = updateUserRequest.Content.ProjectRoles
	}
	if updateUserRequest.Content.ScrumRoles != nil {
		user.Content.ScrumRoles = updateUserRequest.Content.ScrumRoles
	}
	if updateUserRequest.Content.Password != "" {
		user.Content.Password = updateUserRequest.Content.Password
	}

	filter := bson.M{"uuid": UUID}
	update := bson.M{"$set": user}
	_, err = repository.collection.UpdateOne(repository.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (repository *MongoDBUserRepository) DeleteUser(UUID string) (*model.User, error) {
	user := model.User{
		UUID: UUID,
	}

	filter := bson.M{"uuid": UUID}
	_, err := repository.collection.DeleteOne(repository.ctx, filter)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (repository *MongoDBUserRepository) GetUserByID(ID string) (*model.User, error) {
	// Define a filter to find the user by ID
	filter := bson.M{"id": ID}

	// Create an empty User struct to store the result
	var user model.User

	// Find the user using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&user)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (repository *MongoDBUserRepository) GetUserByName(name string) (*model.User, error) {
	// Define a filter to find the user by name
	filter := bson.M{"metadata.name": name}

	// Create an empty User struct to store the result
	var user model.User

	// Find the user using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&user)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (repository *MongoDBUserRepository) GetUserByEmail(email string) (*model.User, error) {
	// Define a filter to find the user by email
	filter := bson.M{"content.email": email}

	// Create an empty User struct to store the result
	var user model.User

	// Find the user using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&user)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return &user, nil
}

func (repository *MongoDBUserRepository) UpdateUserMetadata(UUID string, updatedMetadata model.UserMetadata) (*model.UserMetadata, error) {
	// get user by UUID
	user, err := repository.GetUser(UUID)
	if err != nil {
		return nil, err
	}

	// update user metadata
	user.Metadata = &updatedMetadata

	filter := bson.M{"uuid": UUID}
	update := bson.M{"$set": user}
	_, err = repository.collection.UpdateOne(repository.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return user.Metadata, nil
}

func (repository *MongoDBUserRepository) UpdateUserContent(UUID string, updatedContent model.UserContent) (*model.UserContent, error) {
	// get user by UUID
	user, err := repository.GetUser(UUID)
	if err != nil {
		return nil, err
	}

	// update user content
	user.Content = &updatedContent

	filter := bson.M{"uuid": UUID}
	update := bson.M{"$set": user}
	_, err = repository.collection.UpdateOne(repository.ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return user.Content, nil
}

func (repository *MongoDBUserRepository) GetUserMetadata(UUID string) (*model.UserMetadata, error) {
	// Define a filter to find the user by UUID
	filter := bson.M{"uuid": UUID}

	// Create an empty User struct to store the result
	var user model.User

	// Find the user using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&user)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user.Metadata, nil
}

func (repository *MongoDBUserRepository) GetUserContent(UUID string) (*model.UserContent, error) {
	// Define a filter to find the user by UUID
	filter := bson.M{"uuid": UUID}

	// Create an empty User struct to store the result
	var user model.User

	// Find the user using the filter
	err := repository.collection.FindOne(repository.ctx, filter).Decode(&user)
	if err != nil {
		// Handle the error, such as not found or other errors
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	return user.Content, nil
}
