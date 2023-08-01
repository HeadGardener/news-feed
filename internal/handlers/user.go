package handlers

import (
	"encoding/json"
	"github.com/HeadGardener/news-feed/internal/models"
	"github.com/go-chi/chi/v5"
	"net/http"
	"strconv"
)

func (h *Handler) signUp(w http.ResponseWriter, r *http.Request) {
	var userInput models.UserInput

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding userInput", err)
		return
	}

	if err := userInput.Validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed userInput validation", err)
		return
	}

	id, err := h.userService.Create(r.Context(), userInput)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "user creation failed", err)
		return
	}

	newResponse(w, http.StatusCreated, map[string]interface{}{
		"id": id,
	})
}

func (h *Handler) signIn(w http.ResponseWriter, r *http.Request) {
	var userInput models.UserInput

	if err := json.NewDecoder(r.Body).Decode(&userInput); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while decoding userInput", err)
		return
	}

	if err := userInput.Validate(); err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed userInput validation", err)
		return
	}

	token, err := h.tokenService.GenerateToken(r.Context(), userInput)
	if err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while generating token", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func (h *Handler) updateSendFlag(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "failed while getting userID", err)
		return
	}

	sendFlag, err := strconv.Atoi(chi.URLParam(r, "send_flag"))
	if err != nil {
		newErrResponse(w, http.StatusBadRequest, "invalid send_flag param", err)
		return
	}

	if err := h.userService.UpdateSendFlag(r.Context(), userID, sendFlag); err != nil {
		newErrResponse(w, http.StatusInternalServerError, "failed while updating send_flag", err)
		return
	}

	newResponse(w, http.StatusOK, map[string]string{
		"status": "updated",
	})
}
