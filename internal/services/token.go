package services

import (
	"context"
	"errors"
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
	UserByInput(ctx context.Context, userInput models.UserInput) (models.User, error)
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
	Role   string `json:"role"`
}

func (s *TokenService) GenerateToken(ctx context.Context, userInput models.UserInput) (string, error) {
	userInput.Password = getPasswordHash(userInput.Password)

	user, err := s.userGetter.UserByInput(ctx, userInput)
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
		user.Role,
	})

	return token.SignedString([]byte(secretKey))
}

func (s *TokenService) ParseToken(accessToken string) (models.UserAttributes, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secretKey), nil
	})

	if err != nil {
		return models.UserAttributes{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return models.UserAttributes{}, errors.New("token claims are not of type *tokenClaims")
	}

	return models.UserAttributes{
		ID:    claims.UserID,
		Email: claims.Email,
	}, nil
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
