package handlers

import (
	"context"
	"github.com/HeadGardener/news-feed/internal/models"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"time"
)

type SourceService interface {
	Save(ctx context.Context, srcInput models.SourceInput) (int, error)
}

type ArticleService interface {
	GetAll(ctx context.Context, page int) ([]models.Article, error)
}

type UserService interface {
	Create(ctx context.Context, userInput models.UserInput) (int, error)
}

type TokenService interface {
	GenerateToken(ctx context.Context, userInput models.UserInput) (string, error)
	ParseToken(accessToken string) (models.UserAttributes, error)
}

type FavoritesService interface {
	Add(ctx context.Context, userID, articleID int) error
	GetAll(ctx context.Context, userID int) ([]models.Article, error)
	Delete(ctx context.Context, userID, articleID int) error
}

type Handler struct {
	sourceService    SourceService
	articleService   ArticleService
	userService      UserService
	tokenService     TokenService
	favoritesService FavoritesService
}

func NewHandler(sourceService SourceService,
	articleService ArticleService,
	userService UserService,
	tokenService TokenService,
	favoritesService FavoritesService) *Handler {
	return &Handler{
		sourceService:    sourceService,
		articleService:   articleService,
		userService:      userService,
		tokenService:     tokenService,
		favoritesService: favoritesService,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api", func(r chi.Router) {
		r.Route("/sources", func(r chi.Router) {
			r.Post("/", h.addSource)
		})

		r.Route("/articles", func(r chi.Router) {
			r.Get("/", h.articles)
		})

		r.Route("/users", func(r chi.Router) {
			r.Post("/sign-up", h.signUp)
			r.Post("/sign-in", h.signIn)
		})

		r.Route("/favorites", func(r chi.Router) {
			r.Use(h.identifyUser)
			r.Post("/{article_id}", h.addToFavorites)
			r.Get("/", h.getFavorites)
			r.Delete("/{article_id}", h.deleteFromFavorites)
		})
	})

	return r
}
