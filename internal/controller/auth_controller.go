package controller

import (
	// "net/http"

	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/sirupsen/logrus"
	// "github.com/gin-gonic/gin"
)

type AuthController struct {
	// ctx	  context.Context
	logger		*logrus.Logger
	authService *service.AuthService
}

func NewAuthController(logger *logrus.Logger, authService *service.AuthService) *AuthController {
	return &AuthController{
		// ctx:	   ctx,
		logger:		 logger,
		authService: authService,
	}
}

// // Retrieve the username from the request context
// func GetUsernameFromContext(c *gin.Context) string {
//	if val, ok := c.Get("username"); ok {
//		if username, ok := val.(string); ok {
//			return username
//		}
//	}
//	return ""
// }

// // Retrieve the UUID from the request context
// func getIdFromContext(c *gin.Context) string {
//	if val, ok := c.Get("UUID"); ok {
//		if UUID, ok := val.(string); ok {
//			return UUID
//		}
//	}
//	return ""
// }

// // User login Handler
// // @Summary User login
// // @Description User login
// // @Tags Authentication
// // @Accept json
// // @Produce json
// // @Param body body model.LoginRequest true "User credentials"
// // @Success 200 {object} model.LoginResponse
// // @Failure 400 {object} string "Invalid request"
// // @Failure 500 {object} string "Internal Server Error"
// // @Router /auth/login [post]
// func (ac *AuthController) Login(c *gin.Context) {
//	// username := c.PostRequest("username")
//	// password := c.PostRequest("password")

//	var user model.User
//	// var userRequest model.LoginRequest
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
//		return
//	}

//	token, err := ac.authService.Login(user.Username, user.Password)
//	if err != nil {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
//		return
//	}

//	c.JSON(http.StatusOK, gin.H{"token": token})
// }

// // Protected route handler
// // @Summary Protected Route
// // @Description Protected Route
// // @Tags Authentication
// // @Accept json
// // @Produce json
// // @Success 200 {object} model.SecuredResponse
// // @Failure 400 {object} string "Invalid request"
// // @Failure 500 {object} string "Internal Server Error"
// // @Router /auth/protected [post]
// // @Security ApiKeyAuth
// // @Security BearerAuth
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// func (ac *AuthController) ProtectedRoute(c *gin.Context) {
//	// var securedRequest model.SecuredRequest
//	c.JSON(http.StatusOK, gin.H{"message": "This is a protected route"})
// }

// // User Profile Handler
// // @Summary User Profile
// // @Description User Profile
// // @Tags Authentication
// // @Accept json
// // @Produce json
// // @Success 200 {object} model.User
// // @Failure 400 {object} string "Invalid request"
// // @Failure 500 {object} string "Internal Server Error"
// // @Router /auth/profile [get]
// // @Security ApiKeyAuth
// // @Security BearerAuth
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// func (ac *AuthController) GetUserProfile(c *gin.Context) {
//	// username := GetUsernameFromContext(c)
//	userID := getIdFromContext(c)
//	// var securedRequest model.SecuredRequest

//	user, err := ac.authService.GetUserProfile(userID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}

//	if user == nil {
//		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
//		return
//	}

//	c.JSON(http.StatusOK, user)
// }

// // User Logout Handler
// // @Summary User Logout
// // @Description User Logout
// // @Tags Authentication
// // @Accept json
// // @Produce json
// // @Param UUID path string true "User ID"
// // @Success 200 {object} model.SecuredResponse
// // @Failure 400 {object} string "Invalid request"
// // @Failure 500 {object} string "Internal Server Error"
// // @Router /auth/logout [post]
// // @Security ApiKeyAuth
// // @Security BearerAuth
// // @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// func (ac *AuthController) Logout(c *gin.Context) {
//	userID := c.Param("UUID")
//	// var securedRequest model.SecuredRequest

//	err := ac.authService.Logout(userID)
//	if err != nil {
//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
//		return
//	}

//	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
// }
