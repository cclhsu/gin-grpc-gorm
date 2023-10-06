package controller

import (
	// "net/http"

	"net/http"

	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserControllerInterface interface {
	ListUserIdsAndUUIDs(c *gin.Context)
	ListUsers(c *gin.Context)
	ListUsersMetadata(c *gin.Context)
	ListUsersContent(c *gin.Context)
	GetUser(c *gin.Context)
	UpdateUser(c *gin.Context)
	DeleteUser(c *gin.Context)
	GetUserByID(c *gin.Context)
	GetUserByName(c *gin.Context)
	GetUserByEmail(c *gin.Context)
	UpdateUserMetadata(c *gin.Context)
	UpdateUserContent(c *gin.Context)
	GetUserMetadata(c *gin.Context)
	GetUserContent(c *gin.Context)
}

type UserController struct {
	// ctx	  context.Context
	logger		*logrus.Logger
	userService *service.UserService
}

func NewUserController(logger *logrus.Logger, userService *service.UserService) *UserController {
	return &UserController{
		// ctx:	   ctx,
		logger:		 logger,
		userService: userService,
	}
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/ids-and-uuids | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/ids-and-uuids -H "Authorization: Bearer <token>" | jq
// @Summary Get User IDs and UUIDs
// @Description Get User IDs and UUIDs
// @Tags User
// @Produce json
// @Success 200 {object} model.ListUserIdUuid
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/ids-and-uuids [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) ListUserIdsAndUUIDs(c *gin.Context) {
	idUuids, err := uc.userService.ListUserIdsAndUUIDs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	idUuidArray := []*model.IdUuid{}
	for _, idUuid := range idUuids {
		idUuidArray = append(idUuidArray, &model.IdUuid{
			ID:	  idUuid.ID,
			UUID: idUuid.UUID,
		})
	}

	uc.logger.Info(idUuidArray)
	c.JSON(http.StatusOK, idUuidArray)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users -H "Authorization: Bearer <token>" | jq
// @Summary Get Users
// @Description Get Users
// @Tags User
// @Produce json
// @Success 200 {object} []model.User
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) ListUsers(c *gin.Context) {
	users, err := uc.userService.ListUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(users)
	c.JSON(http.StatusOK, users)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/metadata | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/metadata -H "Authorization: Bearer <token>" | jq
// @Summary Get Users Metadata
// @Description Get Users Metadata
// @Tags User
// @Produce json
// @Success 200 {object} model.ListUsersMetadataResponse
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/metadata [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) ListUsersMetadata(c *gin.Context) {
	users, err := uc.userService.ListUsersMetadata()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(users)
	c.JSON(http.StatusOK, users)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/content | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/content -H "Authorization: Bearer <token>" | jq
// @Summary Get Users Content
// @Description Get Users Content
// @Tags User
// @Produce json
// @Success 200 {object} model.ListUsersContentResponse
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/content [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) ListUsersContent(c *gin.Context) {
	users, err := uc.userService.ListUsersContent()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(users)
	c.JSON(http.StatusOK, users)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000 | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000 -H "Authorization: Bearer <token>" | jq
// @Summary Get User by UUID
// @Description Get User by UUID
// @Tags User
// @Produce json
// @Param UUID path string true "User UUID"
// @Success 200 {object} model.User
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{UUID} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) GetUser(c *gin.Context) {
	UUID := c.Param("UUID")
	user, err := uc.userService.GetUser(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(user)
	c.JSON(http.StatusOK, user)
}

// curl -s -X POST -H 'Content-Type: application/json' http://0.0.0.0:3001/users -d '{"ID": "john.doe", "UUID": "00000000-0000-0000-0000-000000000000", "metadata": {"name": "John Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "john.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "john.doe"}}, "content": {"email": "john.doe@mail.com", "phone": "0912345678", "lastName": "Doe", "firstName": "John", "projectRoles": [1], "scrumRoles": [1], "password": "P@ssw0rd!234" }}' | jq
// curl -s -X POST -H 'Content-Type: application/json' http://0.0.0.0:3001/users -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "UUID": "00000000-0000-0000-0000-000000000000", "metadata": {"name": "Jane Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "jane.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "jane.doe"}}, "content": {"email": "jane.doe@mail.com", "phone": "0987654321", "lastName": "Doe", "firstName": "Jane", "projectRoles": [1], "scrumRoles": [1], "password": "P@ssw0rd!234" }}' | jq
// @Summary Create User
// @Description Create User
// @Tags User
// @Accept json
// @Produce json
// @Param body body model.CreateUserRequest true "New user details"
// @Success 200 {object} string "User created"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users [post]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) CreateUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc.logger.Info(user)

	createUserRequest := model.CreateUserRequest{
		ID:		  user.ID,
		UUID:	  user.UUID,
		Metadata: user.Metadata,
		Content:  user.Content,
	}

	createdUser, err := uc.userService.CreateUser(createUserRequest)
	uc.logger.Info(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, createdUser)
}

// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000 -d '{"ID": "john.doe", "metadata": {"name": "John Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "john.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "john.doe"}}, "content": {"email": "john.doe@example.com", "password": "P@ssw0rd!234"}}' | jq
// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000 -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "metadata": {"name": "Jane Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "jane.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "jane.doe"}}, "content": {"email": "jane.doe@example.com", "password": "P@ssw0rd!234"}}' | jq
// @Summary Update User by UUID
// @Description Update User by UUID
// @Tags User
// @Produce json
// @Param UUID path string true "User UUID"
// @Param body body model.UpdateUserRequest true "Updated user details"
// @Success 200 {object} string "User updated"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{UUID} [put]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) UpdateUser(c *gin.Context) {
	UUID := c.Param("UUID")
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc.logger.Info(user)

	updateUserRequest := model.UpdateUserRequest{
		UUID:	  user.UUID,
		Metadata: user.Metadata,
		Content:  user.Content,
	}

	updatedUser, err := uc.userService.UpdateUser(UUID, updateUserRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// curl -s -X DELETE -H 'Content-Type: application/json' http://0.0.0.0:3001/users/71fe55bc-c9f8-40e8-9031-1683b87aa5f4 | jq
// curl -s -X DELETE -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000 -H "Authorization: Bearer <token>" | jq
// @Summary Delete User by UUID
// @Description Delete User by UUID
// @Tags User
// @Produce json
// @Param UUID path string true "User UUID"
// @Success 200 {object} string "User deleted"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{UUID} [delete]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) DeleteUser(c *gin.Context) {
	UUID := c.Param("UUID")

	_, err := uc.userService.DeleteUser(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/id/john.doe | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/id/john.doe -H "Authorization: Bearer <token>" | jq
// @Summary Get User by ID
// @Description Get User by ID
// @Tags User
// @Produce json
// @Param ID path string true "User ID"
// @Success 200 {object} string "User found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/id/{ID} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) GetUserByID(c *gin.Context) {
	ID := c.Param("id")
	uc.logger.Info(ID)
	user, err := uc.userService.GetUserByID(ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(user)
	c.JSON(http.StatusOK, user)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/name/John%20Doe | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/name/John%20Doe -H "Authorization: Bearer <token>" | jq
// @Summary Get User by Name
// @Description Get User by Name
// @Tags User
// @Produce json
// @Param name path string true "User name"
// @Success 200 {object} string "User found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/name/{name} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) GetUserByName(c *gin.Context) {
	name := c.Param("name")
	user, err := uc.userService.GetUserByName(name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(user)
	c.JSON(http.StatusOK, user)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/email/john.doe@mail.com | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/email/john.doe@mail.com -H "Authorization: Bearer <token>" | jq
// @Summary Get User by Email
// @Description Get User by Email
// @Tags User
// @Produce json
// @Param email path string true "User email"
// @Success 200 {object} string "User found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/email/{email} [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) GetUserByEmail(c *gin.Context) {
	email := c.Param("email")
	user, err := uc.userService.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(user)
	c.JSON(http.StatusOK, user)
}

// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/metadata -d '{"ID": "john.doe", "metadata": {"name": "John Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "john.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "john.doe"}}}' | jq
// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/metadata -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "metadata": {"name": "Jane Doe", "dates": {"createdAt": "2021-01-01T00:00:00.000Z", "createdBy": "jane.doe", "updatedAt": "2021-01-01T00:00:00.000Z", "updatedBy": "jane.doe"}}}' | jq
// @Summary Update User Metadata by UUID
// @Description Update User Metadata by UUID
// @Tags User
// @Produce json
// @Param UUID path string true "User UUID"
// @Param body body model.UpdateUserMetadataRequest true "Updated user metadata"
// @Success 200 {object} string "User metadata updated"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{UUID}/metadata [put]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) UpdateUserMetadata(c *gin.Context) {
	UUID := c.Param("UUID")
	var updateUserMetadataRequest model.UpdateUserMetadataRequest

	if err := c.ShouldBindJSON(&updateUserMetadataRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc.logger.Info(updateUserMetadataRequest)

	updatedUserMetadataResponse, err := uc.userService.UpdateUserMetadata(UUID, updateUserMetadataRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUserMetadataResponse)
}

// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/content -d '{"ID": "john.doe", "content": {"email": "john.doe@example.com", "phone": "0912345678", "password": "P@ssw0rd!234"}}' | jq
// curl -s -X PUT -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/content -H "Authorization: Bearer <token>" -d '{"ID": "jane.doe", "content": {"email": "jane.doe@example.com", "phone": "0987654321", "password": "P@ssw0rd!234"}}' | jq
// @Summary Update User Content by UUID
// @Description Update User Content by UUID
// @Tags User
// @Produce json
// @Param UUID path string true "User UUID"
// @Param body body model.UpdateUserContentRequest true "Updated user content"
// @Success 200 {object} string "User content updated"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{UUID}/content [put]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) UpdateUserContent(c *gin.Context) {
	UUID := c.Param("UUID")
	var updateUserContentRequest model.UpdateUserContentRequest

	if err := c.ShouldBindJSON(&updateUserContentRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	uc.logger.Info(updateUserContentRequest)

	updatedUserContent, err := uc.userService.UpdateUserContent(UUID, updateUserContentRequest)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedUserContent)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/metadata | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/metadata -H "Authorization: Bearer <token>" | jq
// @Summary Get User Metadata by UUID
// @Description Get User Metadata by UUID
// @Tags User
// @Produce json
// @Param UUID path string true "User UUID"
// @Success 200 {object} string "User metadata found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{UUID}/metadata [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) GetUserMetadata(c *gin.Context) {
	UUID := c.Param("UUID")
	userMetadata, err := uc.userService.GetUserMetadata(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(userMetadata)
	c.JSON(http.StatusOK, userMetadata)
}

// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/content | jq
// curl -s -X GET -H 'Content-Type: application/json' http://0.0.0.0:3001/users/00000000-0000-0000-0000-000000000000/content -H "Authorization: Bearer <token>" | jq
// @Summary Get User Content by UUID
// @Description Get User Content by UUID
// @Tags User
// @Produce json
// @Param UUID path string true "User UUID"
// @Success 200 {object} string "User content found"
// @Failure 400 {object} string "Invalid request"
// @Failure 401 {object} string "Unauthorized"
// @Failure 404 {object} string "User not found"
// @Failure 500 {object} string "Internal Server Error"
// @Router /users/{UUID}/content [get]
// @Security ApiKeyAuth
// @Security BearerAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
func (uc *UserController) GetUserContent(c *gin.Context) {
	UUID := c.Param("UUID")
	userContent, err := uc.userService.GetUserContent(UUID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	uc.logger.Info(userContent)
	c.JSON(http.StatusOK, userContent)
}
