package services

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"log"
	"os"
	"time"
)

const path = "./config/vars.env"

var (
	salt      string
	secretKey string
	tokenTTL  = 15 * time.Minute
)

type UserGetter interface {
	GetUserByInput(ctx context.Context, userInput models.UserInput) (models.User, error)
}

type TokenService struct {
	userGetter UserGetter
}

func NewTokenService(userGetter UserGetter) *TokenService {
	initVars()
	return &TokenService{
		userGetter: userGetter,
	}
}

type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int    `json:"id"`
	Email  string `json:"email"`
}

func (s *TokenService) GenerateToken(ctx context.Context, userInput models.UserInput) (string, error) {
	userInput.Password = getPasswordHash(userInput.Password)

	user, err := s.userGetter.GetUserByInput(ctx, userInput)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		user.ID,
		user.Email,
	})

	return token.SignedString([]byte(secretKey))
}

func initVars() {
	err := godotenv.Load(path)
	if err != nil {
		log.Fatalf("[FATAL] invalid path: %s", path)
	}

	salt = os.Getenv("SALT")
	if salt == "" {
		log.Fatalf("[FATAL] SALT is empty")
	}

	secretKey = os.Getenv("TOKEN_SECRET_KEY")
	if secretKey == "" {
		log.Fatalf("[FATAL] TOKEN_SECRET_KEY is empty")
	}
}
