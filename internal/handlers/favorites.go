package handlers

import (
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) addToFavorites(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting userID", err)
		return
	}

	articleID, err := strconv.Atoi(chi.URLParam(r, "article_id"))
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "invalid article_id param", err)
		return
	}

	if err := h.favoritesService.Add(r.Context(), userID, articleID); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while adding to favorites", err)
	}

	newResponse(w, http.StatusOK, map[string]string{
		"status": "successfully added",
	})
}

func (h *Handler) getFavorites(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting userID", err)
		return
	}

	articles, err := h.favoritesService.GetAll(r.Context(), userID)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while getting favorites", err)
		return
	}

	newResponse(w, http.StatusOK, articles)
}

func (h *Handler) deleteFromFavorites(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting userID", err)
		return
	}

	articleID, err := strconv.Atoi(chi.URLParam(r, "article_id"))
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "invalid article_id param", err)
		return
	}

	if err := h.favoritesService.Delete(r.Context(), userID, articleID); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while deleting from favorites", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]string{
		"status": "successfully deleted",
	})
}
