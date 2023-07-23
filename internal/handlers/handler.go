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
}

type Handler struct {
	sourceService  SourceService
	articleService ArticleService
	userService    UserService
	tokenService   TokenService
}

func NewHandler(srcSvc SourceService, artSvc ArticleService, uSvc UserService, tknSvc TokenService) *Handler {
	return &Handler{
		sourceService:  srcSvc,
		articleService: artSvc,
		userService:    uSvc,
		tokenService:   tknSvc,
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
	})

	return r
}
