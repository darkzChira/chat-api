package service

import (
	"chat-app/internal/models"
	"chat-app/internal/repository"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
	"log"
	"time"
)

type UserService struct {
	userRepo  *repository.UserRepository
	jwtSecret []byte
}

func NewUserService(userRepo *repository.UserRepository, jwtSecret []byte) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) RegisterUser(ctx context.Context, user *models.User) error {
	log.Println("uesrseive- register user")
	existingUser, err := s.userRepo.FindUserByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return errors.New("username already taken")
	}

	if err != nil && err != mongo.ErrNoDocuments {
		return err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.OnlineStatus = false

	return s.userRepo.SaveUser(ctx, user)
}

func (s *UserService) LoginUser(ctx context.Context, username, password string) (string, *models.User, error) {
	user, err := s.userRepo.FindUserByUsername(ctx, username)
	if err != nil {
		log.Println("User not found:", err)
		return "", nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Password mismatch for user:", username)
		return "", nil, errors.New("invalid password")
	}

	user.OnlineStatus = true
	err = s.userRepo.UpdateUser(ctx, user)
	if err != nil {
		log.Println("Failed to update user online status:", err)
		return "", nil, err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		log.Println("Failed to sign JWT token:", err)
		return "", nil, err
	}

	log.Println("JWT token generated for user:", username)

	return tokenString, user, nil
}

func (s *UserService) LogoutUser(ctx context.Context, userID string) error {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	user.OnlineStatus = false
	return s.userRepo.UpdateUser(ctx, user)
}

func (s *UserService) GetUser(ctx context.Context, userID string) (*models.User, error) {
	user, err := s.userRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) GetOtherUsers(ctx context.Context, currentUserID string) ([]models.User, error) {
	return s.userRepo.FindAllExceptUserID(ctx, currentUserID)
}
