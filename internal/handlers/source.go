package handlers

import (
	"encoding/json"
	"github.com/HeadGardener/news-feed/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
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

func (h *Handler) deleteSource(w http.ResponseWriter, r *http.Request) {
	sourceID, err := strconv.Atoi(chi.URLParam(r, "source_id"))
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "invalid source_id param", err)
		return
	}

	if err := h.sourceService.Delete(r.Context(), sourceID); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while deleting source", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]string{
		"status": "deleted",
	})
}
