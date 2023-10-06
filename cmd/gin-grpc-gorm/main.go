package main

// go run ./cmd/gin-grpc-gorm/main.go
// swag init -g cmd/gin-grpc-gorm/main.go -o doc/openapi
// go build -o ./bin/gin-grpc-gorm ./cmd/gin-grpc-gorm
// ./bin/gin-grpc-gorm

// go get github.com/golang-jwt/jwt
// go get github.com/go-redis/redis/v8
// go get github.com/joho/godotenv
// go get github.com/gin-gonic/gin
// go get google.golang.org/grpc
// go get google.golang.org/grpc/codes
// go get google.golang.org/grpc/metadata
// go get google.golang.org/grpc/status
// go get github.com/swaggo/swag
// go get google.golang.org/protobuf/reflect/protoreflect
// go get google.golang.org/protobuf/runtime/protoimpl
// go get google.golang.org/protobuf/types/known/emptypb
// go get github.com/swaggo/files
// go get github.com/swaggo/gin-swagger
// go get google.golang.org/grpc/reflection
// go get github.com/patrickmn/go-cache

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"

	"net"
	"os"

	cache "github.com/cclhsu/gin-grpc-gorm/internal/cache/redis"
	"github.com/cclhsu/gin-grpc-gorm/internal/config"
	"github.com/cclhsu/gin-grpc-gorm/internal/model"
	mongodb_repository "github.com/cclhsu/gin-grpc-gorm/internal/repository/mongodb"
	postgres_repository "github.com/cclhsu/gin-grpc-gorm/internal/repository/postgres_gorm"
	"github.com/cclhsu/gin-grpc-gorm/internal/route"
	"github.com/cclhsu/gin-grpc-gorm/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	// "go.mongodb.org/mongo-driver/bson"
)

var (
	ctx    context.Context
	logger *logrus.Logger

	host     string
	port     string
	endpoint string
	router   *gin.Engine

	redisClient   *redis.Client
	redisCache    *cache.RedisCache
	redisHost     string
	redisPort     string
	redisPassword string
	redisEndpoint string

	mongoDBClient   *mongo.Client
	mongoDB         *gorm.DB
	mongoDBHost     string
	mongoDBPort     string
	mongoDBUser     string
	mongoDBPassword string
	mongoDBName     string
	mongoDBEndpoint string

	postgresDB       *gorm.DB
	postgresHost     string
	postgresPort     string
	postgresUser     string
	postgresPassword string
	postgresDBName   string
	postgresEndpoint string

	mongoDBUserRepository  *mongodb_repository.MongoDBUserRepository
	mongoDBTeamRepository  *mongodb_repository.MongoDBTeamRepository
	postgresUserRepository *postgres_repository.PostgresUserRepository
	postgresTeamRepository *postgres_repository.PostgresTeamRepository

	grpcHost     string
	grpcPort     string
	grpcEndpoint string
	grpcServer   *grpc.Server

	authService   *service.AuthService
	userService   *service.UserService
	teamService   *service.TeamService
	helloService  *service.HelloService
	healthService *service.HealthService
)

// CallerPrettyfier is a function that formats the caller information.
func CallerPrettyfier(f *runtime.Frame) (string, string) {
	// Split the file path by slashes
	parts := strings.Split(f.File, "/")
	// Get the last part of the split, which is the file name
	fileName := parts[len(parts)-1]
	// Format the caller information as "filename:line"
	return fmt.Sprintf("%5d ", f.Line), fmt.Sprintf("%20v ", fileName)
}

func setupLogger() {
	// Create a new logrus logger
	logger = logrus.New()
	logger.SetOutput(os.Stdout)

	// // Create a new loggly hook
	// hook := logrusly.NewLogglyHook("https://logs-01.loggly.com/inputs/0b0f0f1e-0b0b-4b0b-8b0b-0b0b0b0b0b0b/tag/http/", "gin-grpc-gorm", logrus.InfoLevel)

	// // Add the hook to the logger
	// logger.Hooks.Add(hook)

	// Set the logger formatter
	logger.SetReportCaller(true)

	// // Set the logger formatter
	// logger.SetFormatter(&logrus.JSONFormatter{})
	// logrus.SetFormatter(&logrus.JSONFormatter{
	//	DisableTimestamp: false,
	//	TimestampFormat:  "2006-01-02 15:04:05",
	//	FieldMap: logrus.FieldMap{
	//		logrus.FieldKeyTime:  "@timestamp",
	//		logrus.FieldKeyLevel: "@level",
	//		logrus.FieldKeyMsg:	  "@message",
	//	}})
	// logger.SetFormatter(&logrus.TextFormatter{
	//	FullTimestamp:			true,
	//	TimestampFormat:		"2006-01-02 15:04:05",
	//	DisableLevelTruncation: true,
	//	CallerPrettyfier:		CallerPrettyfier,
	//	// You can customize other formatting options here
	// })
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors:    false, // Disable colored output
		FullTimestamp:    true,  // Include the timestamp
		TimestampFormat:  time.RFC3339,
		CallerPrettyfier: CallerPrettyfier,
	})

	// Set the logger level: info, debug, warn, error, fatal
	logger.SetLevel(logrus.TraceLevel)

	// Log a message
	logger.Info("Hello world!")
}

func createRedisCache() (*cache.RedisCache, *redis.Client, error) {
	// Cache
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisEndpoint := fmt.Sprintf("%s:%s", redisHost, redisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: redisPassword,
		DB:       0,
	})

	// Ping the Redis server to check the connection
	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		logger.Errorf("Failed to connect to Redis: %v\n", err)
		return nil, nil, err
	}

	// // Create a Redis cache instance
	// cacheManager := cache.New(cache.NoExpiration, cache.NoExpiration)

	redisCache := cache.NewRedisCache(ctx, logger, redisClient)

	// defer redisClient.Close()
	return redisCache, redisClient, nil
}

func closeMongoDBConnection(client *mongo.Client, ctx context.Context, cancel context.CancelFunc) {

	// CancelFunc to cancel to context
	defer cancel()

	// client provides a method to close
	// a mongoDB connection.
	defer func() {

		// client.Disconnect method also has deadline.
		// returns error if any,
		if err := client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()
}

func connectToMongoDatabase(ctx context.Context) (*mongo.Client, context.Context, context.CancelFunc, error) {
	// Database configuration
	mongoDBHost := os.Getenv("MONGO_HOST")
	mongoDBPort := os.Getenv("MONGO_PORT")
	mongoDBUser := os.Getenv("MONGO_USER")
	mongoDBPassword := os.Getenv("MONGO_PASS")
	// mongoDBName := os.Getenv("MONGO_DB")

	// Build the MongoDB connection URI
	// mongoDBURI := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s",
	//	mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBPort, mongoDBName)
	mongoDBURI := fmt.Sprintf("mongodb://%s:%s@%s:%s",
		mongoDBUser, mongoDBPassword, mongoDBHost, mongoDBPort)

	logger.Info("MongoDB URI: ", mongoDBURI)

	// Set client options
	clientOptions := options.Client().ApplyURI(mongoDBURI)

	// Set a 30-second timeout
	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	// Connect to MongoDB
	mongoDBClient, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("Failed to connect to MongoDB: %w", err)
	}

	// userCollection := mongoDBClient.Database(mongoDBName).Collection("users")
	// teamCollection := mongoDBClient.Database(mongoDBName).Collection("teams")

	return mongoDBClient, ctx, nil, nil
}

func PingMongoDB(ctx context.Context, mongoDBClient *mongo.Client) error {

	if err := mongoDBClient.Ping(ctx, readpref.Primary()); err != nil {
		return err
	}
	logger.Info("Ping to MongoDB successfully")
	return nil
}

func connectToPostgresDatabase() (*gorm.DB, error) {
	// Database
	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPassword := os.Getenv("POSTGRES_PASS")
	postgresDBName := os.Getenv("POSTGRES_DB")

	postgresEndpoint := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		postgresHost, postgresPort, postgresUser, postgresPassword, postgresDBName)

	// Initialize the GORM instance
	postgresDB, err := gorm.Open(postgres.Open(postgresEndpoint), &gorm.Config{})
	if err != nil {
		logger.Fatal("Failed to initialize GORM:", err)
	}

	// AutoMigrate model
	err = postgresDB.AutoMigrate(&model.User{})
	if err != nil {
		return nil, fmt.Errorf("Failed to auto migrate model: %w", err)
	}

	err = postgresDB.AutoMigrate(&model.Team{})
	if err != nil {
		return nil, fmt.Errorf("Failed to auto migrate model: %w", err)
	}

	return postgresDB, nil
}

func startGinServer() {
	host = os.Getenv("SERVICE_HOST")
	port = os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	endpoint = host + ":" + port

	// Set Gin to production mode
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Set up Gin server
	router := gin.Default()

	route.SetupRestfulRoutes(router, host, port, logger, authService, userService, teamService, helloService, healthService)
	// route.SetupGraphQLRoutes(router, host, port, logger, authService, userService, teamService, helloService, healthService)

	// Add redis client to gin context
	router.Use(func(c *gin.Context) {
		c.Set("redis", redisClient)
		c.Next()
	})

	fmt.Printf("Starting up on http://%s/\n", endpoint)
	fmt.Printf("Starting up on http://%s/doc/openapi/index.html\n", endpoint)
	// Start the Gin server
	err := router.Run(endpoint)
	if err != nil {
		logger.Fatalf("Failed to start the server: %v", err)
	}
}

func startGrpCServiceServer() {
	grpcHost = os.Getenv("GRPC_HOST")
	grpcPort = os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}
	grpcEndpoint = grpcHost + ":" + grpcPort

	// Set up gRPC server
	grpcServer = grpc.NewServer(
	// grpc.UnaryInterceptor(middleware.AuthInterceptor),
	)

	// Register reflection service on gRPC server.
	reflection.Register(grpcServer)

	// Register gRPC service implementations
	route.SetupGrpcRoutes(ctx, logger, grpcServer, authService, userService, teamService, helloService, healthService)

	fmt.Printf("Starting up on http://%s/\n", grpcEndpoint)
	// Start the gRPC server
	lis, err := net.Listen("tcp", grpcEndpoint)
	if err != nil {
		logger.Fatalf("Failed to listen: %v", err)
	}
	if err := grpcServer.Serve(lis); err != nil {
		logger.Fatalf("Failed to serve: %v", err)
	}
}

// @title My API
// @description This is a sample API server using Gin and Swagger.
// @version 1.0
// @BasePath /
// @schemes http https
func main() {
	// Create a context
	ctx = context.Background()

	// Set up logger
	setupLogger()

	// Load environment variables
	err := config.LoadEnv()
	if err != nil {
		logger.Fatal("Failed to load environment variables")
	}

	// check REPOSITORY_DEVICE_TYPE: database, file
	repositoryDeviceType := os.Getenv("REPOSITORY_DEVICE_TYPE")

	// check DATABASE_REPOSITORY_DEVICE_TYPE: mongo, postgres
	databaseRepositoryDeviceType := os.Getenv("DATABASE_REPOSITORY_DEVICE_TYPE")

	// check FILE_REPOSITORY_DEVICE_TYPE: json, memory, s3, minio
	fileRepositoryDeviceType := os.Getenv("FILE_REPOSITORY_DEVICE_TYPE")

	// check if the REDIS_ENABLE is true
	redisEnable := os.Getenv("REDIS_ENABLE")

	// check if the MONGO_ENABLE is true
	mongoEnable := os.Getenv("MONGO_ENABLE")

	// check if the POSTGRES_ENABLE is true
	postgresEnable := os.Getenv("POSTGRES_ENABLE")

	if redisEnable == "true" {
		// Create Redis cache
		redisCache, redisClient, err = createRedisCache()
		if err != nil {
			logger.Errorf("Failed to create Redis cache: %v\n", err)
			// return // disable redis cache if failed to create
		}

		logger.Info("Connected to Redis!")
	} else {
		redisCache = nil
		redisClient = nil
	}

	if repositoryDeviceType == "database" {
		if databaseRepositoryDeviceType == "mongodb" && mongoEnable == "true" {
			// Connect to the database
			mongoDBClient, ctx, _, err = connectToMongoDatabase(ctx)
			if err != nil {
				logger.Fatalf("Failed to connect to the database: %v", err)
			}

			err = PingMongoDB(ctx, mongoDBClient)
			if err != nil {
				logger.Fatalf("Failed to ping the database: %v", err)
			}

			// Create the repository
			mongoDBUserRepository = mongodb_repository.NewMongoDBUserRepository(ctx, logger, mongoDBClient)
			mongoDBTeamRepository = mongodb_repository.NewMongoDBTeamRepository(ctx, logger, mongoDBClient)

			// Create the user service with cache and repository
			userService = service.NewUserService(ctx, logger, mongoDBUserRepository, redisCache)
			teamService = service.NewTeamService(ctx, logger, mongoDBTeamRepository, redisCache)

			// Create the auth service with cache and repository
			authService = service.NewAuthService(ctx, logger, mongoDBUserRepository, redisCache)

			logger.Info("Connected to MongoDB!")
		}

		if databaseRepositoryDeviceType == "postgres" && postgresEnable == "true" {
			// Connect to the database
			postgresDB, err = connectToPostgresDatabase()
			if err != nil {
				logger.Fatalf("Failed to connect to the database: %v", err)
			}

			// Create the repository
			postgresUserRepository = postgres_repository.NewPostgresUserRepository(ctx, logger, postgresDB)
			postgresTeamRepository = postgres_repository.NewPostgresTeamRepository(ctx, logger, postgresDB)

			// Create the user service with cache and repository
			userService = service.NewUserService(ctx, logger, postgresUserRepository, redisCache)
			teamService = service.NewTeamService(ctx, logger, postgresTeamRepository, redisCache)

			// Create the auth service with cache and repository
			authService = service.NewAuthService(ctx, logger, postgresUserRepository, redisCache)

			logger.Info("Connected to Postgres!")
		}
	} else if repositoryDeviceType == "file" {
		// TODO: create file repository
		if fileRepositoryDeviceType == "json" {
			// TODO: create json file repository

			logger.Info("Connected to JSON file!")
		} else if fileRepositoryDeviceType == "memory" {
			// TODO: create memory file repository

			logger.Info("Connected to memory file!")
		} else if fileRepositoryDeviceType == "s3" {
			//	TODO: create s3 file repository

			logger.Info("Connected to S3 file!")
		} else if fileRepositoryDeviceType == "minio" {
			// TODO: create minio file repository

			logger.Info("Connected to Minio file!")
		}
	}

	// Create the hello service
	helloService = service.NewHelloService(ctx, logger)

	// Create the health check service
	healthService = service.NewHealthService(ctx, logger)

	// // Start Gin server in a goroutine
	// go startGinServer()

	// Start gRPC server in a goroutine
	go startGrpCServiceServer()

	// Block the main goroutine to keep the servers running
	select {}
}
