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

// UserServiceInterface defines the interface for the user service.
type UserServiceInterface interface {
	// Define your repository methods here...
	ListUserIdsAndUUIDs() ([]model.IdUuid, error)
	ListUsers() ([]*model.User, error)
	GetUser(uuid string) (*model.User, error)
	CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error)
	UpdateUser(uuid string, updateUserRequest model.UpdateUserRequest) (*model.User, error)
	DeleteUser(uuid string) (*model.User, error)
	GetUserByID(ID string) (*model.User, error)
	GetUserByName(name string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	ListUsersMetadata() (*model.ListUsersMetadataResponse, error)
	ListUsersContent() (*model.ListUsersContentResponse, error)
	UpdateUserMetadata(uuid string, updateUserMetadataRequest model.UpdateUserMetadataRequest) (*model.UserMetadata, error)
	UpdateUserContent(uuid string, updateUserContentRequest model.UpdateUserContentRequest) (*model.UserContent, error)
	GetUserMetadata(uuid string) (*model.UserMetadata, error)
	GetUserContent(uuid string) (*model.UserContent, error)
	IsUserExist(name, email, ID, UUID string) (bool, error)
	IsNoUserExist(name, email, ID, UUID string) (bool, error)
	IsExactlyOneUserExist(name, email, ID, UUID string) (bool, error)
	IsAtLeastOneUserExist(name, email, ID, UUID string) (bool, error)
}

// ErrUserNotFound is a custom error for "User not found" cases.
var ErrUserNotFound = errors.New("User not found")

// UserService represents the UserService type.
type UserService struct {
	ctx			   context.Context
	logger		   *logrus.Logger
	userRepository repository.UserRepositoryInterface
	cacheManager   *cache.Cache
	redisCache	   *redis_cache.RedisCache
	isCacheEnabled bool
}

// NewUserService creates a new instance of UserService.
func NewUserService(ctx context.Context, logger *logrus.Logger, userRepository repository.UserRepositoryInterface, redisCache *redis_cache.RedisCache) *UserService {
	cacheManager := cache.New(cache.NoExpiration, cache.NoExpiration)
	isCacheEnabled := os.Getenv("CACHE_ENABLED") == "true"

	return &UserService{
		ctx:			ctx,
		logger:			logger,
		userRepository: userRepository,
		cacheManager:	cacheManager,
		redisCache:		redisCache,
		isCacheEnabled: isCacheEnabled,
	}
}

// ListUserIdsAndUUIDs lists user IDs and UUIDs.
func (us *UserService) ListUserIdsAndUUIDs() ([]model.IdUuid, error) {
	if us.isCacheEnabled {
		if users, found := us.cacheManager.Get("usersIdsAndUUIDs"); found {
			return users.([]model.IdUuid), nil
		}
	}

	users, err := us.userRepository.ListUserIdsAndUUIDs()
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		us.cacheManager.Set("usersIdsAndUUIDs", users, cache.DefaultExpiration)
	}

	us.logger.Printf("Users: %+v\n", users)
	return users, nil
}

// ListUsers lists users.
func (us *UserService) ListUsers() ([]*model.User, error) {
	if us.isCacheEnabled {
		if users, found := us.cacheManager.Get("users"); found {
			return users.([]*model.User), nil
		}
	}

	users, err := us.userRepository.ListUsers()
	if err != nil {
		return nil, err
	}

	usersArray := []*model.User{}
	for _, user := range users {
		usersArray = append(usersArray, &model.User{
			ID:	  user.ID,
			UUID: user.UUID,
			Metadata: &model.UserMetadata{
				Name: user.Metadata.Name,
				Dates: &model.CommonDate{
					CreatedAt:	 user.Metadata.Dates.CreatedAt,
					CreatedBy:	 user.Metadata.Dates.CreatedBy,
					UpdatedAt:	 user.Metadata.Dates.UpdatedAt,
					UpdatedBy:	 user.Metadata.Dates.UpdatedBy,
					StartDate:	 user.Metadata.Dates.StartDate,
					EndDate:	 user.Metadata.Dates.EndDate,
					StartedAt:	 user.Metadata.Dates.StartedAt,
					StartedBy:	 user.Metadata.Dates.StartedBy,
					CompletedAt: user.Metadata.Dates.CompletedAt,
					CompletedBy: user.Metadata.Dates.CompletedBy,
				},
			},
			Content: &model.UserContent{
				Email:		  user.Content.Email,
				Phone:		  user.Content.Phone,
				LastName:	  user.Content.LastName,
				FirstName:	  user.Content.FirstName,
				ProjectRoles: user.Content.ProjectRoles,
				ScrumRoles:	  user.Content.ScrumRoles,
				Password:	  user.Content.Password,
			},
		})
	}

	if us.isCacheEnabled {
		us.cacheManager.Set("users", usersArray, cache.DefaultExpiration)
	}

	us.logger.Printf("Users: %+v\n", usersArray)
	return usersArray, nil
}

// GetUser gets a user by UUID.
func (us *UserService) GetUser(uuid string) (*model.User, error) {
	if us.isCacheEnabled {
		if user, found := us.cacheManager.Get(fmt.Sprintf("user:%s", uuid)); found {
			return user.(*model.User), nil
		}
	}

	user, err := us.userRepository.GetUser(uuid)
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		us.cacheManager.Set(fmt.Sprintf("user:%s", uuid), user, cache.DefaultExpiration)
	}

	us.logger.Printf("User: %+v\n", user)
	return user, nil
}

// CreateUser creates a user.
func (us *UserService) CreateUser(createUserRequest model.CreateUserRequest) (*model.User, error) {
	// check if UUID is empty string, undefined, null, and '00000000-0000-0000-0000-000000000000', generate a new UUID if so
	if createUserRequest.UUID == "" || createUserRequest.UUID == "undefined" || createUserRequest.UUID == "null" || createUserRequest.UUID == "00000000-0000-0000-0000-000000000000" {
		createUserRequest.UUID = uuid.New().String()
	}

	// Check if a user with the same name, email, ID, or UUID already exists
	userExists, err := us.IsUserExist(createUserRequest.Metadata.Name, createUserRequest.Content.Email, createUserRequest.ID, createUserRequest.UUID)
	us.logger.Printf("userExists: %+v\n", userExists)
	us.logger.Printf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "user not found" {
			// Expected error
		} else {
			// Handle other errors
			return nil, err
		}
	}

	if userExists {
		return nil, fmt.Errorf("user with the same name, email, ID, or UUID already exists")
	}

	// // Validate DTO metadata and content
	// if err := createUserRequest.Validate(); err != nil {
	//	return nil, err
	// }

	// Create a User instance from the CreateUserRequest
	user, err := us.userRepository.CreateUser(createUserRequest)
	if err != nil {
		return nil, err
	}

	// Create a User instance from the CreateUserRequest
	if us.isCacheEnabled {
		us.cacheManager.Set(fmt.Sprintf("user:%s", user.UUID), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByName:%s", user.Metadata.Name), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByEmail:%s", user.Content.Email), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByID:%s", user.ID), user, cache.DefaultExpiration)
	}

	us.logger.Printf("User: %+v\n", user)
	return user, nil
}

// UpdateUser updates a user.
func (us *UserService) UpdateUser(uuid string, updateUserRequest model.UpdateUserRequest) (*model.User, error) {
	// Check if user exists, and retrieve the user
	user, err := us.GetUser(uuid)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update the dates values UpdatedAt and UpdatedBy
	updateUserRequest.Metadata.Dates = &model.CommonDate{
		// CreatedAt:	user.Metadata.Dates.CreatedAt,
		// CreatedBy:	user.Metadata.Dates.CreatedBy,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedBy: updateUserRequest.Metadata.Dates.UpdatedBy,
		// StartDate:	user.Metadata.Dates.StartDate,
		// EndDate:		user.Metadata.Dates.EndDate,
		// StartedAt:	user.Metadata.Dates.StartedAt,
		// StartedBy:	user.Metadata.Dates.StartedBy,
		// CompletedAt: user.Metadata.Dates.CompletedAt,
		// CompletedBy: user.Metadata.Dates.CompletedBy,
	}

	// // Validate DTO metadata and content
	// if err := updateUserRequest.Validate(); err != nil {
	//	return nil, err
	// }

	// Update the user
	user, err = us.userRepository.UpdateUser(uuid, updateUserRequest)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	// Update the user in the cache
	if us.isCacheEnabled {
		us.cacheManager.Set(fmt.Sprintf("user:%s", user.UUID), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByName:%s", user.Metadata.Name), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByEmail:%s", user.Content.Email), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByID:%s", user.ID), user, cache.DefaultExpiration)
	}

	us.logger.Printf("User: %+v\n", user)
	return user, nil
}

// DeleteUser deletes a user.
func (us *UserService) DeleteUser(uuid string) (*model.User, error) {
	user, err := us.userRepository.DeleteUser(uuid)
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		us.cacheManager.Delete(fmt.Sprintf("user:%s", uuid))
		us.cacheManager.Delete(fmt.Sprintf("userByName:%s", user.Metadata.Name))
		us.cacheManager.Delete(fmt.Sprintf("userByEmail:%s", user.Content.Email))
		us.cacheManager.Delete(fmt.Sprintf("userByID:%s", user.ID))
	}

	us.logger.Printf("User: %+v\n", user)
	return user, nil
}

// GetUserByID gets a user by ID.
func (us *UserService) GetUserByID(ID string) (*model.User, error) {
	if us.isCacheEnabled {
		if user, found := us.cacheManager.Get(fmt.Sprintf("userByID:%s", ID)); found {
			return user.(*model.User), nil
		}
	}

	user, err := us.userRepository.GetUserByID(ID)
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		us.cacheManager.Set(fmt.Sprintf("userByID:%s", ID), user, cache.DefaultExpiration)
	}

	us.logger.Printf("User: %+v\n", user)
	return user, nil
}

// GetUserByName gets a user by name.
func (us *UserService) GetUserByName(name string) (*model.User, error) {
	if us.isCacheEnabled {
		if user, found := us.cacheManager.Get(fmt.Sprintf("userByName:%s", name)); found {
			return user.(*model.User), nil
		}
	}

	user, err := us.userRepository.GetUserByName(name)
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		us.cacheManager.Set(fmt.Sprintf("userByName:%s", name), user, cache.DefaultExpiration)
	}

	us.logger.Printf("User: %+v\n", user)
	return user, nil
}

// GetUserByEmail gets a user by email.
func (us *UserService) GetUserByEmail(email string) (*model.User, error) {
	if us.isCacheEnabled {
		if user, found := us.cacheManager.Get(fmt.Sprintf("userByEmail:%s", email)); found {
			return user.(*model.User), nil
		}
	}

	user, err := us.userRepository.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		us.cacheManager.Set(fmt.Sprintf("userByEmail:%s", email), user, cache.DefaultExpiration)
	}

	us.logger.Printf("User: %+v\n", user)
	return user, nil
}

// ListUsersMetadata lists users with metadata.
func (us *UserService) ListUsersMetadata() (*model.ListUsersMetadataResponse, error) {
	if us.isCacheEnabled {
		if users, found := us.cacheManager.Get("usersMetadata"); found {
			return users.(*model.ListUsersMetadataResponse), nil
		}
	}

	users, err := us.userRepository.ListUsers()
	if err != nil {
		return nil, err
	}

	userMetadataResponses := make([]*model.UserMetadataResponse, len(users))
	for i, user := range users {
		userMetadataResponses[i] = &model.UserMetadataResponse{
			UUID: user.UUID,
			ID:	  user.ID,
			Metadata: &model.UserMetadata{
				Name: user.Metadata.Name,
				Dates: &model.CommonDate{
					CreatedAt:	 user.Metadata.Dates.CreatedAt,
					CreatedBy:	 user.Metadata.Dates.CreatedBy,
					UpdatedAt:	 user.Metadata.Dates.UpdatedAt,
					UpdatedBy:	 user.Metadata.Dates.UpdatedBy,
					StartDate:	 user.Metadata.Dates.StartDate,
					EndDate:	 user.Metadata.Dates.EndDate,
					StartedAt:	 user.Metadata.Dates.StartedAt,
					StartedBy:	 user.Metadata.Dates.StartedBy,
					CompletedAt: user.Metadata.Dates.CompletedAt,
					CompletedBy: user.Metadata.Dates.CompletedBy,
				},
			},
		}
	}
	listUsersMetadataResponse := &model.ListUsersMetadataResponse{
		UserMetadataResponses: userMetadataResponses,
	}

	if us.isCacheEnabled {
		us.cacheManager.Set("listUsersMetadataResponse", listUsersMetadataResponse, cache.DefaultExpiration)
	}

	us.logger.Printf("Users: %+v\n", listUsersMetadataResponse)
	return listUsersMetadataResponse, nil
}

// ListUsersContent lists users with content.
func (us *UserService) ListUsersContent() (*model.ListUsersContentResponse, error) {
	if us.isCacheEnabled {
		if users, found := us.cacheManager.Get("usersContent"); found {
			return users.(*model.ListUsersContentResponse), nil
		}
	}

	users, err := us.userRepository.ListUsers()
	if err != nil {
		return nil, err
	}

	userContentResponses := make([]*model.UserContentResponse, len(users))
	for i, user := range users {
		userContentResponses[i] = &model.UserContentResponse{
			UUID: user.UUID,
			ID:	  user.ID,
			Content: &model.UserContent{
				Email:		  user.Content.Email,
				Phone:		  user.Content.Phone,
				LastName:	  user.Content.LastName,
				FirstName:	  user.Content.FirstName,
				ProjectRoles: user.Content.ProjectRoles,
				ScrumRoles:	  user.Content.ScrumRoles,
			},
		}
	}
	listUsersContentResponse := &model.ListUsersContentResponse{
		UserContentResponses: userContentResponses,
	}

	if us.isCacheEnabled {
		us.cacheManager.Set("listUsersContentResponse", listUsersContentResponse, cache.DefaultExpiration)
	}

	us.logger.Printf("Users: %+v\n", listUsersContentResponse)
	return listUsersContentResponse, nil
}

// UpdateUserMetadata updates a user's metadata.
func (us *UserService) UpdateUserMetadata(uuid string, updateUserMetadataRequest model.UpdateUserMetadataRequest) (*model.UserMetadata, error) {
	newUserMetadata := model.UserMetadata{
		Name:  updateUserMetadataRequest.Metadata.Name,
		Dates: updateUserMetadataRequest.Metadata.Dates,
	}
	userMetadata, err := us.userRepository.UpdateUserMetadata(uuid, newUserMetadata)
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		user, _ := us.userRepository.GetUser(uuid)
		user.Metadata = userMetadata
		us.cacheManager.Set(fmt.Sprintf("user:%s", uuid), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByName:%s", user.Metadata.Name), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByEmail:%s", user.Content.Email), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByID:%s", user.ID), user, cache.DefaultExpiration)
	}

	us.logger.Printf("UserMetadata: %+v\n", userMetadata)
	return userMetadata, nil
}

// UpdateUserContent updates a user's content.
func (us *UserService) UpdateUserContent(uuid string, updateUserContentRequest model.UpdateUserContentRequest) (*model.UserContent, error) {
	newUserContent := model.UserContent{
		Email:		  updateUserContentRequest.Content.Email,
		Phone:		  updateUserContentRequest.Content.Phone,
		LastName:	  updateUserContentRequest.Content.LastName,
		FirstName:	  updateUserContentRequest.Content.FirstName,
		ProjectRoles: updateUserContentRequest.Content.ProjectRoles,
		ScrumRoles:	  updateUserContentRequest.Content.ScrumRoles,
		Password:	  updateUserContentRequest.Content.Password,
	}
	userContent, err := us.userRepository.UpdateUserContent(uuid, newUserContent)
	if err != nil {
		return nil, err
	}

	if us.isCacheEnabled {
		user, _ := us.userRepository.GetUser(uuid)
		user.Content = userContent
		us.cacheManager.Set(fmt.Sprintf("user:%s", uuid), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByName:%s", user.Metadata.Name), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByEmail:%s", user.Content.Email), user, cache.DefaultExpiration)
		us.cacheManager.Set(fmt.Sprintf("userByID:%s", user.ID), user, cache.DefaultExpiration)
	}

	us.logger.Printf("UserContent: %+v\n", userContent)
	return userContent, nil
}

// GetUserMetadata gets a user's metadata.
func (us *UserService) GetUserMetadata(uuid string) (*model.UserMetadata, error) {
	userMetadata, err := us.userRepository.GetUserMetadata(uuid)
	if err != nil {
		return nil, err
	}

	us.logger.Printf("UserMetadata: %+v\n", userMetadata)
	return userMetadata, nil
}

// GetUserContent gets a user's content.
func (us *UserService) GetUserContent(uuid string) (*model.UserContent, error) {
	userContent, err := us.userRepository.GetUserContent(uuid)
	if err != nil {
		return nil, err
	}

	us.logger.Printf("UserContent: %+v\n", userContent)
	return userContent, nil
}

// IsUserExist checks if a user with the specified attributes exists.
func (us *UserService) IsUserExist(name, email, ID, UUID string) (bool, error) {
	if !us.isCacheEnabled {
		return us.checkUserExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("userExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if userExists, found := us.cacheManager.Get(cacheKey); found {
		return userExists.(bool), nil
	}

	userExists, err := us.checkUserExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	us.cacheManager.Set(cacheKey, userExists, cache.DefaultExpiration)
	return userExists, nil
}

// IsNoUserExist checks if no user with the specified attributes exists.
func (us *UserService) IsNoUserExist(name, email, ID, UUID string) (bool, error) {
	if !us.isCacheEnabled {
		return us.checkUserExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("noUserExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if userExists, found := us.cacheManager.Get(cacheKey); found {
		return userExists.(bool), nil
	}

	userExists, err := us.checkUserExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	us.cacheManager.Set(cacheKey, userExists, cache.DefaultExpiration)
	return userExists, nil
}

// IsExactlyOneUserExist checks if exactly one user with the specified attributes exists.
func (us *UserService) IsExactlyOneUserExist(name, email, ID, UUID string) (bool, error) {
	if !us.isCacheEnabled {
		return us.checkUserExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("exactlyOneUserExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if userExists, found := us.cacheManager.Get(cacheKey); found {
		return userExists.(bool), nil
	}

	userExists, err := us.checkUserExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	us.cacheManager.Set(cacheKey, userExists, cache.DefaultExpiration)
	return userExists, nil
}

// IsAtLeastOneUserExist checks if at least one user with the specified attributes exists.
func (us *UserService) IsAtLeastOneUserExist(name, email, ID, UUID string) (bool, error) {
	if !us.isCacheEnabled {
		return us.checkUserExistence(name, email, ID, UUID)
	}

	cacheKey := fmt.Sprintf("atLeastOneUserExistence:%s-%s-%s-%s", name, email, ID, UUID)
	if userExists, found := us.cacheManager.Get(cacheKey); found {
		return userExists.(bool), nil
	}

	userExists, err := us.checkUserExistence(name, email, ID, UUID)
	if err != nil {
		return false, err
	}

	us.cacheManager.Set(cacheKey, userExists, cache.DefaultExpiration)
	return userExists, nil
}

// checkUserExistence checks if a user with the specified attributes exists.
func (us *UserService) checkUserExistence(name, email, ID, UUID string) (bool, error) {
	userByName, err := us.userRepository.GetUserByName(name)
	us.logger.Debugf("name: %+v\n", name)
	us.logger.Debugf("userByName: %+v\n", userByName)
	us.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "user not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if userByName != nil {
		return true, nil
	}

	userByEmail, err := us.userRepository.GetUserByEmail(email)
	us.logger.Debugf("email: %+v\n", email)
	us.logger.Debugf("userByEmail: %+v\n", userByEmail)
	us.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "user not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if userByEmail != nil {
		return true, nil
	}

	userByID, err := us.userRepository.GetUserByID(ID)
	us.logger.Debugf("ID: %+v\n", ID)
	us.logger.Debugf("userByID: %+v\n", userByID)
	us.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "user not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if userByID != nil {
		return true, nil
	}

	userByUUID, err := us.userRepository.GetUser(UUID)
	us.logger.Debugf("UUID: %+v\n", UUID)
	us.logger.Debugf("userByUUID: %+v\n", userByUUID)
	us.logger.Debugf("err: %+v\n", err)
	if err != nil {
		if err.Error() == "user not found" {
			// Expected error
		} else {
			// Handle other errors
			return false, err
		}
	}
	if userByUUID != nil {
		return true, nil
	}

	return false, nil
}
