package route

import (
	"fmt"
	"net/http"

	openapiDocs "github.com/cclhsu/gin-grpc-gorm/doc/openapi"
	"github.com/cclhsu/gin-grpc-gorm/internal/controller"
	"github.com/sirupsen/logrus"

	// "github.com/cclhsu/gin-grpc-gorm/internal/middleware"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/gin-gonic/gin"

	// "github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// SetupRoutes sets up the API routes
func SetupRestfulRoutes(r *gin.Engine, host string, port string, logger *logrus.Logger, authService *service.AuthService, userService *service.UserService, teamService *service.TeamService, helloService *service.HelloService, healthService *service.HealthService) {

	// Create instances of the controller
	userController := controller.NewUserController(logger, userService)
	teamController := controller.NewTeamController(logger, teamService)
	// authController := controller.NewAuthController(logger, authService)
	helloController := controller.NewHelloController(logger, helloService)
	healthController := controller.NewHealthController(logger, healthService)

	// Enable CORS middleware
	r.Use(func(c *gin.Context) {
		origin := fmt.Sprintf("http://%s:%s", host, port)
		c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// cors.DefaultConfig()
		// corsConfig := cors.DefaultConfig()
		// corsConfig.AllowAllOrigins = true
		// corsConfig.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
		// corsConfig.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
		// corsConfig.AllowCredentials = true
		// corsConfig.AddAllowHeaders("Connection")

		// Handle preflight OPTIONS request
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	})

	// helloGroup := r.Group("/api/v1/hello")
	helloGroup := r.Group("/hello")
	{
		// Get hello world json
		helloGroup.GET("/json", helloController.GetHelloJson)

		// Get hello world string
		helloGroup.GET("/string", helloController.GetHelloString)
	}

	// docGroup := r.Group("/api/v1/doc")
	docGroup := r.Group("/doc")
	{
		openapiDocs.SwaggerInfo.BasePath = "/"

		// Serve Swagger documentation
		// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		docGroup.GET("/openapi/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	// // authGroup := r.Group("/api/v1/auth")
	// authGroup := r.Group("/auth")
	// {
	//	// Login endpoint
	//	authGroup.POST("/login", authController.Login)

	//	// Protected route
	//	protected := authGroup.Group("/")
	//	protected.Use(middleware.AuthMiddleware())
	//	{
	//		// Protected route
	//		protected.POST("/protected", authController.ProtectedRoute)

	//		// Get user profile
	//		protected.GET("/profile", authController.GetUserProfile)

	//		// User logout
	//		authGroup.POST("/logout", authController.Logout)
	//	}
	// }

	// userGroup := r.Group("/api/v1/users")
	userGroup := r.Group("/users")
	// userGroup.Use(middleware.AuthMiddleware())
	{
		// Get all users id adn uuid
		userGroup.GET("ids-and-uuids", userController.ListUserIdsAndUUIDs)

		// Get all users
		userGroup.GET("", userController.ListUsers)

		// Get all users metadata
		userGroup.GET("/metadata", userController.ListUsersMetadata)

		// Get all users content
		userGroup.GET("/content", userController.ListUsersContent)

		// Get user by UUID
		userGroup.GET("/:UUID", userController.GetUser)

		// Create a new user
		userGroup.POST("", userController.CreateUser)

		// Update user by ID
		userGroup.PUT("/:UUID", userController.UpdateUser)

		// Delete user by ID
		userGroup.DELETE("/:UUID", userController.DeleteUser)

		// Get user by ID
		userGroup.GET("/id/:id", userController.GetUserByID)

		// Get user by name
		userGroup.GET("/name/:name", userController.GetUserByName)

		// Get user by email
		userGroup.GET("/email/:email", userController.GetUserByEmail)

		// Update user metadata by UUID
		userGroup.PUT("/:UUID/metadata", userController.UpdateUserMetadata)

		// Update user content by UUID
		userGroup.PUT("/:UUID/content", userController.UpdateUserContent)

		// Get user metadata by UUID
		userGroup.GET("/:UUID/metadata", userController.GetUserMetadata)

		// Get user content by UUID
		userGroup.GET("/:UUID/content", userController.GetUserContent)
	}

	// teamGroup := r.Group("/api/v1/teams")
	teamGroup := r.Group("/teams")
	// teamGroup.Use(middleware.AuthMiddleware())
	{
		// Get all teams id adn uuid
		teamGroup.GET("ids-and-uuids", teamController.ListTeamIdsAndUUIDs)

		// Get all teams
		teamGroup.GET("", teamController.ListTeams)

		// Get all teams metadata
		teamGroup.GET("/metadata", teamController.ListTeamsMetadata)

		// Get all teams content
		teamGroup.GET("/content", teamController.ListTeamsContent)

		// Get team by UUID
		teamGroup.GET("/:UUID", teamController.GetTeam)

		// Create a new team
		teamGroup.POST("", teamController.CreateTeam)

		// Update team by ID
		teamGroup.PUT("/:UUID", teamController.UpdateTeam)

		// Delete team by ID
		teamGroup.DELETE("/:UUID", teamController.DeleteTeam)

		// Get team by ID
		teamGroup.GET("/id/:id", teamController.GetTeamByID)

		// Get team by name
		teamGroup.GET("/name/:name", teamController.GetTeamByName)

		// Get team by email
		teamGroup.GET("/email/:email", teamController.GetTeamByEmail)

		// Update team metadata by UUID
		teamGroup.PUT("/:UUID/metadata", teamController.UpdateTeamMetadata)

		// Update team content by UUID
		teamGroup.PUT("/:UUID/content", teamController.UpdateTeamContent)

		// Get team metadata by UUID
		teamGroup.GET("/:UUID/metadata", teamController.GetTeamMetadata)

		// Get team content by UUID
		teamGroup.GET("/:UUID/content", teamController.GetTeamContent)
	}

	// healthGroup := r.Group("/api/v1/health")
	healthGroup := r.Group("/health")
	{
		// Get health check
		healthGroup.GET("/healthy", healthController.IsHealthy)

		// Get health check
		healthGroup.GET("/live", healthController.IsALive)

		// Get health check
		healthGroup.GET("/ready", healthController.IsReady)
	}
}
