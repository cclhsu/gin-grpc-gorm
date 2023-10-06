// Path: internal/repository/user/memory_user_repository.go
// DESC: This is the memory implementation of the user repository.
package postgres_gorm

import (
	"context"
	"errors"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
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

// PostgresUserRepository implements UserRepositoryInterface for PostgreSQL.
type PostgresUserRepository struct {
	ctx	   context.Context
	logger *logrus.Logger
	db	   *gorm.DB
}

// NewPostgresUserRepository creates a new instance of PostgresUserRepository.
func NewPostgresUserRepository(ctx context.Context, logger *logrus.Logger, db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		ctx:	ctx,
		logger: logger,
		db:		db,
	}
}

// Error messages
var (
	ErrNotFound		= errors.New("record not found")
	ErrDatabase		= errors.New("database error")
	ErrAlreadyExist = errors.New("record already exists")
)

// handleErr is a helper function to handle errors consistently.
func handleErr(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}

	return ErrDatabase
}

// ListUserIdsAndUUIDs returns a list of user IDs and UUIDs.
func (repository *PostgresUserRepository) ListUserIdsAndUUIDs() ([]model.IdUuid, error) {
	var idUuids []model.IdUuid
	if err := repository.db.WithContext(repository.ctx).Find(&idUuids).Error; err != nil {
		return nil, handleErr(err)
	}
	return idUuids, nil
}

// ListUsers returns a list of users.
func (repository *PostgresUserRepository) ListUsers() ([]model.User, error) {
	var users []model.User
	if err := repository.db.WithContext(repository.ctx).Find(&users).Error; err != nil {
		return nil, handleErr(err)
	}
	return users, nil
}

// GetUser retrieves a user by UUID.
func (repository *PostgresUserRepository) GetUser(UUID string) (*model.User, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return &user, nil
}

// CreateUser creates a new user.
func (repository *PostgresUserRepository) CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error) {
	user := model.User{
		ID:	  createUserRequest.ID,
		UUID: createUserRequest.UUID,
		Metadata: &model.UserMetadata{
			Name: createUserRequest.Metadata.Name,
			Dates: &model.CommonDate{
				CreatedAt:	 createUserRequest.Metadata.Dates.CreatedAt,
				CreatedBy:	 createUserRequest.Metadata.Dates.CreatedBy,
				UpdatedAt:	 createUserRequest.Metadata.Dates.UpdatedAt,
				UpdatedBy:	 createUserRequest.Metadata.Dates.UpdatedBy,
				StartDate:	 createUserRequest.Metadata.Dates.StartDate,
				EndDate:	 createUserRequest.Metadata.Dates.EndDate,
				StartedAt:	 createUserRequest.Metadata.Dates.StartedAt,
				StartedBy:	 createUserRequest.Metadata.Dates.StartedBy,
				CompletedAt: createUserRequest.Metadata.Dates.CompletedAt,
				CompletedBy: createUserRequest.Metadata.Dates.CompletedBy,
			},
		},
		Content: &model.UserContent{
			Email:		  createUserRequest.Content.Email,
			Phone:		  createUserRequest.Content.Phone,
			LastName:	  createUserRequest.Content.LastName,
			FirstName:	  createUserRequest.Content.FirstName,
			ProjectRoles: createUserRequest.Content.ProjectRoles,
			ScrumRoles:	  createUserRequest.Content.ScrumRoles,
			Password:	  createUserRequest.Content.Password,
		},
	}

	if err := repository.db.WithContext(repository.ctx).Create(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return &user, nil
}

// UpdateUser updates a user by UUID.
func (repository *PostgresUserRepository) UpdateUser(UUID string, updateUserRequest model.UpdateUserRequest) (*model.User, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}

	user.Metadata.Name = updateUserRequest.Metadata.Name
	user.Metadata.Dates.UpdatedAt = updateUserRequest.Metadata.Dates.UpdatedAt
	user.Metadata.Dates.UpdatedBy = updateUserRequest.Metadata.Dates.UpdatedBy
	user.Metadata.Dates.StartDate = updateUserRequest.Metadata.Dates.StartDate
	user.Metadata.Dates.EndDate = updateUserRequest.Metadata.Dates.EndDate
	user.Metadata.Dates.StartedAt = updateUserRequest.Metadata.Dates.StartedAt
	user.Metadata.Dates.StartedBy = updateUserRequest.Metadata.Dates.StartedBy
	user.Metadata.Dates.CompletedAt = updateUserRequest.Metadata.Dates.CompletedAt
	user.Metadata.Dates.CompletedBy = updateUserRequest.Metadata.Dates.CompletedBy
	user.Content.Email = updateUserRequest.Content.Email
	user.Content.Phone = updateUserRequest.Content.Phone
	user.Content.LastName = updateUserRequest.Content.LastName
	user.Content.FirstName = updateUserRequest.Content.FirstName
	user.Content.ProjectRoles = updateUserRequest.Content.ProjectRoles
	user.Content.ScrumRoles = updateUserRequest.Content.ScrumRoles
	user.Content.Password = updateUserRequest.Content.Password

	if err := repository.db.WithContext(repository.ctx).Save(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return &user, nil
}

// DeleteUser deletes a user by UUID.
func (repository *PostgresUserRepository) DeleteUser(UUID string) (*model.User, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}

	if err := repository.db.WithContext(repository.ctx).Delete(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return &user, nil
}

// GetUserByID retrieves a user by ID.
func (repository *PostgresUserRepository) GetUserByID(ID string) (*model.User, error) {
	var user model.User
	repository.logger.Info(ID)
	if err := repository.db.WithContext(repository.ctx).Where("id = ?", ID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return &user, nil
}

// GetUserByName retrieves a user by name.
func (repository *PostgresUserRepository) GetUserByName(name string) (*model.User, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("name = ?", name).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by email.
func (repository *PostgresUserRepository) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("email = ?", email).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return &user, nil
}

// UpdateUserMetadata updates user metadata by UUID.
func (repository *PostgresUserRepository) UpdateUserMetadata(UUID string, updatedMetadata model.UserMetadata) (*model.UserMetadata, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}

	user.Metadata.Name = updatedMetadata.Name
	user.Metadata.Dates.UpdatedAt = updatedMetadata.Dates.UpdatedAt
	user.Metadata.Dates.UpdatedBy = updatedMetadata.Dates.UpdatedBy
	user.Metadata.Dates.StartDate = updatedMetadata.Dates.StartDate
	user.Metadata.Dates.EndDate = updatedMetadata.Dates.EndDate
	user.Metadata.Dates.StartedAt = updatedMetadata.Dates.StartedAt
	user.Metadata.Dates.StartedBy = updatedMetadata.Dates.StartedBy
	user.Metadata.Dates.CompletedAt = updatedMetadata.Dates.CompletedAt
	user.Metadata.Dates.CompletedBy = updatedMetadata.Dates.CompletedBy

	if err := repository.db.WithContext(repository.ctx).Save(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return user.Metadata, nil
}

// UpdateUserContent updates user content by UUID.
func (repository *PostgresUserRepository) UpdateUserContent(UUID string, updatedContent model.UserContent) (*model.UserContent, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}

	user.Content.Email = updatedContent.Email
	user.Content.Phone = updatedContent.Phone
	user.Content.LastName = updatedContent.LastName
	user.Content.FirstName = updatedContent.FirstName
	user.Content.ProjectRoles = updatedContent.ProjectRoles
	user.Content.ScrumRoles = updatedContent.ScrumRoles

	if err := repository.db.WithContext(repository.ctx).Save(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return user.Content, nil
}

// GetUserMetadata retrieves user metadata by UUID.
func (repository *PostgresUserRepository) GetUserMetadata(UUID string) (*model.UserMetadata, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return user.Metadata, nil
}

// GetUserContent retrieves user content by UUID.
func (repository *PostgresUserRepository) GetUserContent(UUID string) (*model.UserContent, error) {
	var user model.User
	if err := repository.db.WithContext(repository.ctx).Where("uuid = ?", UUID).First(&user).Error; err != nil {
		return nil, handleErr(err)
	}
	return user.Content, nil
}
