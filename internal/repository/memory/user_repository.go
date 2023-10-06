// Path: ineternal/repository/user/memory_user_repository.go
// DESC: This is the memory implementation of the user repository.
package memory

import (
	"context"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/sirupsen/logrus"
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

type MemoryUserRepository struct {
	ctx	   context.Context
	logger *logrus.Logger
}

func NewMemoryUserRepository(ctx context.Context, logger *logrus.Logger) *MemoryUserRepository {
	return &MemoryUserRepository{
		ctx:	ctx,
		logger: logger,
	}
}

func (repository *MemoryUserRepository) ListUserIdsAndUUIDs() ([]model.IdUuid, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) ListUsers() ([]model.User, error) {
	return nil, nil
}

// func (repository *MemoryUserRepository) ListUsersMetadata() ([]model.UserMetadata, error) {
//	return nil, nil
// }

// func (repository *MemoryUserRepository) ListUsersContent() ([]model.UserContent, error) {
//	return nil, nil
// }

func (repository *MemoryUserRepository) GetUser(UUID string) (*model.User, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) UpdateUser(UUID string, updateUserRequest model.UpdateUserRequest) (*model.User, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) DeleteUser(UUID string) (*model.User, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) GetUserByID(ID string) (*model.User, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) GetUserByName(name string) (*model.User, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) GetUserByEmail(email string) (*model.User, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) UpdateUserMetadata(UUID string, updatedMetadata model.UserMetadata) (*model.UserMetadata, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) UpdateUserContent(UUID string, updatedContent model.UserContent) (*model.UserContent, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) GetUserMetadata(UUID string) (*model.UserMetadata, error) {
	return nil, nil
}

func (repository *MemoryUserRepository) GetUserContent(UUID string) (*model.UserContent, error) {
	return nil, nil
}
