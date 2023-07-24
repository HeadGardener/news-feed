package handlers

import (
	"context"
	"errors"
	"fmt"
	"github.com/HeadGardener/news-feed/internal/models"
	"net/http"
	"strings"
)

const (
	userCtx = "userAtr"
)

func (h *Handler) identifyUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := r.Header.Get("Authorization")

		if header == "" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("empty auth header"))
			return
		}

		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("invalid auth header, must be like `Bearer token`"))
			return
		}

		if headerParts[0] != "Bearer" {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New(fmt.Sprintf("invalid auth header %s, must be Bearer", headerParts[0])))
			return
		}

		if len(headerParts[1]) == 0 {
			newErrResponse(w, http.StatusUnauthorized, "failed while identifying user",
				errors.New("jwt token is empty"))
			return
		}

		userAttributes, err := h.tokenService.ParseToken(headerParts[1])
		if err != nil {
			newErrResponse(w, http.StatusUnauthorized, "failed while parsing token", err)
			return
		}

		ctx := context.WithValue(r.Context(), userCtx, userAttributes)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserID(r *http.Request) (int, error) {
	userCtxValue := r.Context().Value(userCtx)
	userAttributes, ok := userCtxValue.(models.UserAttributes)
	if !ok {
		return 0, errors.New("userCtx value is not of type UserAttributes")
	}

	return userAttributes.ID, nil
}
