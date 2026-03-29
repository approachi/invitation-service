package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"invitation-service/internal/handler"

	"github.com/go-chi/chi/v5"
)

func newRouter(h *handler.InvitationHandler) http.Handler {
	r := chi.NewRouter()
	r.Post("/api/v1/invitations/{code}", h.UseInvitation)
	return r
}

func startServer(srv *http.Server) error {
	go func() {
		log.Printf("starting server on %s", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Println("server stopped")
	return nil
}
