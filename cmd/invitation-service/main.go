package main

import (
	"log"
	"net/http"

	"invitation-service/internal/handler"
	"invitation-service/internal/repository"
	"invitation-service/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	pool, err := initDB(cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()

	repo := repository.NewPostgresRepository(pool)
	svc := service.NewInvitationService(repo)
	h := handler.NewInvitationHandler(svc)

	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: newRouter(h),
	}

	return startServer(srv)
}
