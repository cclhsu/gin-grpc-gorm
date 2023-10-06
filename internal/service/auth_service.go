package service

import (
	"context"

	cache "github.com/cclhsu/gin-grpc-gorm/internal/cache/redis"
	repository "github.com/cclhsu/gin-grpc-gorm/internal/repository/postgres_gorm"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	ctx			   context.Context
	logger		   *logrus.Logger
	userRepository repository.UserRepositoryInterface
	redisCache	   *cache.RedisCache
}

func NewAuthService(ctx context.Context, logger *logrus.Logger, userRepository repository.UserRepositoryInterface, redisCache *cache.RedisCache) *AuthService {
	return &AuthService{
		ctx:			ctx,
		logger:			logger,
		userRepository: userRepository,
		redisCache:		redisCache,
	}
}

// func (as *AuthService) Login(username, password string) (string, error) {
//	// Perform authentication logic
//	// Generate and return JWT access token

//	// Example implementation (replace with your actual logic):
//	// Find the user by username in the repository
//	user, err := as.userRepository.GetByUsername(username)
//	if err != nil {
//		return "", err
//	}

//	// Check if the user exists and the password is correct
//	if user == nil || user.Password != password {
//		return "", errors.New("invalid username or password")
//	}

//	// Generate and return a JWT access token
//	accessToken, err := generateAccessToken(user)
//	return accessToken, nil
// }

// func (as *AuthService) Logout(userID string) error {
//	// Perform logout logic

//	// Example implementation (replace with your actual logic):
//	// Perform any necessary actions to log out the user
//	// such as invalidating the access token or session
//	// This can be specific to your application's requirements

//	// For this example, we don't perform any action
//	return nil
// }

// func (as *AuthService) GetUserProfile(userID string) (*model.User, error) {
//	// Retrieve user profile based on userID
//	// Return the user profile

//	// Example implementation (replace with your actual logic):
//	// Retrieve the user from the repository by ID
//	user, err := as.userRepository.GetByID(userID)
//	if err != nil {
//		return nil, err
//	}

//	// Check if the user exists
//	if user == nil {
//		return nil, errors.New("user not found")
//	}

//	// Return the user profile
//	return user, nil
// }

// // Helper functions for token generation and verification
// // Replace these functions with your actual token library or implementation

// func generateAccessToken(user *model.User) (string, error) {
//	expirationTime := time.Now().Add(5 * time.Minute)
//	claims := &model.Claims{
//		Username: user.ID,
//		Sub:	  user.UUID,
//		StandardClaims: jwt.StandardClaims{
//			ExpiresAt: expirationTime.Unix(),
//		},
//	}

//	jwtSecret := []byte(os.Getenv("JWT_SECRET"))

//	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//	tokenString, err := token.SignedString(jwtSecret)
//	if err != nil {
//		return "", err
//	}

//	return tokenString, nil
// }
