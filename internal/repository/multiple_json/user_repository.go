// Path: internal/repository/user/multiple_json_user_repository.go
// DESC: This is the multiple json implementation of the user repository.
package multiple_json

// import (
//	"github.com/cclhsu/gin-grpc-gorm/internal/model"
// )

// // UserRepositoryInterface defines the interface for the user repository.
// type UserRepositoryInterface interface {
//	ListUserIdsAndUUIDs() ([]model.IdUuid, error)
//	ListUsers() ([]model.User, error)
//	// ListUsersMetadata() ([]model.UserMetadata, error)
//	// ListUsersContent() ([]model.UserContent, error)
//	GetUser(UUID string) (*model.User, error)
//	CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error)
//	UpdateUser(UUID string, updateUserRequest model.UpdateUserRequest) (*model.User, error)
//	DeleteUser(UUID string) (*model.User, error)
//	GetUserByID(ID string) (*model.User, error)
//	GetUserByName(name string) (*model.User, error)
//	GetUserByEmail(email string) (*model.User, error)
//	UpdateUserMetadata(UUID string, updatedMetadata model.UserMetadata) (*model.UserMetadata, error)
//	UpdateUserContent(UUID string, updatedContent model.UserContent) (*model.UserContent, error)
//	GetUserMetadata(UUID string) (*model.UserMetadata, error)
//	GetUserContent(UUID string) (*model.UserContent, error)
// }

// type MultipleJsonUserRepository struct {
// }

// func NewMultipleJsonUserRepository() *MultipleJsonUserRepository {
//	return &MultipleJsonUserRepository{}
// }

// func (repository *MultipleJsonUserRepository) ListUserIdsAndUUIDs() ([]model.IdUuid, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) ListUsers() ([]model.User, error) {
//	return nil, nil
// }

// // func (repository *MultipleJsonUserRepository) ListUsersMetadata() ([]model.UserMetadata, error) {
// //	return nil, nil
// // }

// // func (repository *MultipleJsonUserRepository) ListUsersContent() ([]model.UserContent, error) {
// //	return nil, nil
// // }

// func (repository *MultipleJsonUserRepository) GetUser(UUID string) (*model.User, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) UpdateUser(UUID string, updateUserRequest model.UpdateUserRequest) (*model.User, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) DeleteUser(UUID string) (*model.User, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) GetUserByID(ID string) (*model.User, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) GetUserByName(name string) (*model.User, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) GetUserByEmail(email string) (*model.User, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) UpdateUserMetadata(UUID string, updatedMetadata model.UserMetadata) (*model.UserMetadata, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) UpdateUserContent(UUID string, updatedContent model.UserContent) (*model.UserContent, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) GetUserMetadata(UUID string) (*model.UserMetadata, error) {
//	return nil, nil
// }

// func (repository *MultipleJsonUserRepository) GetUserContent(UUID string) (*model.UserContent, error) {
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

// UserRepositoryInterface defines the interface for the user repository.
type UserRepositoryInterface interface {
	ListUserIdsAndUUIDs() ([]model.IdUuid, error)
	ListUsers() ([]model.User, error)
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

// MultipleJsonUserRepository implements UserRepositoryInterface using multiple JSON files.
type MultipleJsonUserRepository struct {
	ctx			  context.Context
	logger		  *logrus.Logger
	dirPath		  string
	fileExtension string
	lock		  sync.RWMutex
}

// NewMultipleJsonUserRepository creates a new instance of MultipleJsonUserRepository.
func NewMultipleJsonUserRepository(ctx context.Context, logger *logrus.Logger, dirPath string) UserRepositoryInterface {
	return &MultipleJsonUserRepository{
		ctx:		   ctx,
		logger:		   logger,
		dirPath:	   dirPath,
		fileExtension: ".json",
	}
}

func (r *MultipleJsonUserRepository) listFiles() ([]string, error) {
	files, err := fileutils.ListFiles(r.dirPath, r.fileExtension)
	if err != nil {
		return nil, err
	}

	return files, nil
}

// ListUserIdsAndUUIDs returns a list of user IDs and UUIDs.
func (r *MultipleJsonUserRepository) ListUserIdsAndUUIDs() ([]model.IdUuid, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	userIdsAndUUIDs := make([]model.IdUuid, 0)

	for _, filePath := range files {
		userID := strings.TrimSuffix(filepath.Base(filePath), r.fileExtension)
		userUUID, err := uuid.Parse(userID)
		if err != nil {
			return nil, err
		}

		userIdsAndUUIDs = append(userIdsAndUUIDs, model.IdUuid{
			ID:	  userID,
			UUID: userUUID.String(),
		})
	}

	return userIdsAndUUIDs, nil
}

// ListUsers returns a list of users.
func (r *MultipleJsonUserRepository) ListUsers() ([]model.User, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	users := make([]model.User, 0)

	for _, filePath := range files {
		user, err := r.readUserFromFile(filePath)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

// GetUser retrieves a user by UUID.
func (r *MultipleJsonUserRepository) GetUser(UUID string) (*model.User, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.RLock()
	defer r.lock.RUnlock()

	return r.readUserFromFile(filePath)
}

func (r *MultipleJsonUserRepository) readUserFromFile(filePath string) (*model.User, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var user model.User
	if err := json.Unmarshal(fileData, &user); err != nil {
		return nil, err
	}

	return &user, nil
}

// CreateUser creates a new user.
func (r *MultipleJsonUserRepository) CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error) {
	newUUID := uuid.New()
	newUser := model.NewUser(newUUID, createUserRequest.Metadata, createUserRequest.Content)
	filePath := filepath.Join(r.dirPath, newUUID.String()+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	if err := r.writeUserToFile(filePath, newUser); err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *MultipleJsonUserRepository) writeUserToFile(filePath string, user *model.User) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(filePath, userJSON, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// UpdateUser updates a user by UUID.
func (r *MultipleJsonUserRepository) UpdateUser(UUID string, updateUserRequest model.UpdateUserRequest) (*model.User, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	existingUser, err := r.readUserFromFile(filePath)
	if err != nil {
		return nil, err
	}

	existingUser.UpdateMetadata(updateUserRequest.Metadata)
	existingUser.UpdateContent(updateUserRequest.Content)

	if err := r.writeUserToFile(filePath, existingUser); err != nil {
		return nil, err
	}

	return existingUser, nil
}

// DeleteUser deletes a user by UUID.
func (r *MultipleJsonUserRepository) DeleteUser(UUID string) (*model.User, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	deletedUser, err := r.readUserFromFile(filePath)
	if err != nil {
		return nil, err
	}

	if err := os.Remove(filePath); err != nil {
		return nil, err
	}

	return deletedUser, nil
}

// GetUserByID retrieves a user by ID.
func (r *MultipleJsonUserRepository) GetUserByID(ID string) (*model.User, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, filePath := range files {
		user, err := r.readUserFromFile(filePath)
		if err != nil {
			return nil, err
		}

		if user.ID == ID {
			return user, nil
		}
	}

	return nil, fmt.Errorf("User with ID %s not found", ID)
}

// GetUserByName retrieves a user by name.
func (r *MultipleJsonUserRepository) GetUserByName(name string) (*model.User, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, filePath := range files {
		user, err := r.readUserFromFile(filePath)
		if err != nil {
			return nil, err
		}

		if user.Metadata.Name == name {
			return user, nil
		}
	}

	return nil, fmt.Errorf("User with name %s not found", name)
}

// GetUserByEmail retrieves a user by email.
func (r *MultipleJsonUserRepository) GetUserByEmail(email string) (*model.User, error) {
	files, err := r.listFiles()
	if err != nil {
		return nil, err
	}

	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, filePath := range files {
		user, err := r.readUserFromFile(filePath)
		if err != nil {
			return nil, err
		}

		if user.Content.Email == email {
			return user, nil
		}
	}

	return nil, fmt.Errorf("User with email %s not found", email)
}

// UpdateUserMetadata updates user metadata by UUID.
func (r *MultipleJsonUserRepository) UpdateUserMetadata(UUID string, updatedMetadata model.UserMetadata) (*model.UserMetadata, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	existingUser, err := r.readUserFromFile(filePath)
	if err != nil {
		return nil, err
	}

	existingUser.UpdateMetadata(updatedMetadata)

	if err := r.writeUserToFile(filePath, existingUser); err != nil {
		return nil, err
	}

	return existingUser.Metadata, nil
}

// UpdateUserContent updates user content by UUID.
func (r *MultipleJsonUserRepository) UpdateUserContent(UUID string, updatedContent model.UserContent) (*model.UserContent, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.Lock()
	defer r.lock.Unlock()

	existingUser, err := r.readUserFromFile(filePath)
	if err != nil {
		return nil, err
	}

	existingUser.UpdateContent(updatedContent)

	if err := r.writeUserToFile(filePath, existingUser); err != nil {
		return nil, err
	}

	return existingUser.Content, nil
}

// GetUserMetadata retrieves user metadata by UUID.
func (r *MultipleJsonUserRepository) GetUserMetadata(UUID string) (*model.UserMetadata, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.RLock()
	defer r.lock.RUnlock()

	existingUser, err := r.readUserFromFile(filePath)
	if err != nil {
		return nil, err
	}

	return existingUser.Metadata, nil
}

// GetUserContent retrieves user content by UUID.
func (r *MultipleJsonUserRepository) GetUserContent(UUID string) (*model.UserContent, error) {
	filePath := filepath.Join(r.dirPath, UUID+r.fileExtension)

	r.lock.RLock()
	defer r.lock.RUnlock()

	existingUser, err := r.readUserFromFile(filePath)
	if err != nil {
		return nil, err
	}

	return existingUser.Content, nil
}
