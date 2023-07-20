package handlers

import (
	"encoding/json"
	"github.com/HeadGardener/news-feed/internal/models"
	"net/http"
)

func (h *Handler) addSource(w http.ResponseWriter, r *http.Request) {
	var srcInput models.SourceInput

	if err := json.NewDecoder(r.Body).Decode(&srcInput); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding srcInput", err)
		return
	}

	if err := srcInput.Validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "srcInput validation failed", err)
		return
	}

	id, err := h.sourceService.Save(r.Context(), srcInput)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed source save", err)
		return
	}

	newResponse(w, http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}
