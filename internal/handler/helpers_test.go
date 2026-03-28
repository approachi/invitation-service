package handler_test

import (
	"net/http"

	"invitation-service/internal/handler"
	"invitation-service/internal/service/mocks"

	"github.com/go-chi/chi/v5"
)

func newTestRouter(svc *mocks.InvitationService) http.Handler {
	h := handler.NewInvitationHandler(svc)
	r := chi.NewRouter()
	r.Post("/api/v1/invitations/{code}", h.UseInvitation)
	return r
}
