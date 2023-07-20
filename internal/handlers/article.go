package handlers

import (
	"net/http"
	"strconv"
)

func (h *Handler) articles(w http.ResponseWriter, r *http.Request) {
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "invalid page query param", err)
		return
	}

	if page == 0 {
		page = 1
	}

	articles, err := h.articleService.GetAll(r.Context(), page)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed to get articles", err)
		return
	}

	newResponse(w, http.StatusOK, articles)
}
