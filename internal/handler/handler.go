package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"invitation-service/internal/domain"
	"invitation-service/internal/service"

	"github.com/go-chi/chi/v5"
)

type InvitationHandler struct {
	svc service.InvitationService
}

func NewInvitationHandler(svc service.InvitationService) *InvitationHandler {
	return &InvitationHandler{svc: svc}
}

type useInvitationRequest struct {
	Email string `json:"Email"`
}

func (h *InvitationHandler) UseInvitation(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	var req useInvitationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := h.svc.UseInvitation(r.Context(), code, req.Email)
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	if domain.IsClientError(err) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	log.Printf("unexpected error: %v", err)
	w.WriteHeader(http.StatusInternalServerError)
}
