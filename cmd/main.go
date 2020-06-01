package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/DanielTitkov/netology-slowly/internal/app"
	"github.com/DanielTitkov/netology-slowly/internal/configs"
	"github.com/DanielTitkov/netology-slowly/internal/middleware"
	"github.com/go-chi/chi"
	chimiddleware "github.com/go-chi/chi/middleware"
)

func main() {
	log.Println("service starting")
	cfg, err := configs.ReadConfigs("./configs/default.yaml")
	if err != nil {
		log.Fatal("failed to load config", err)
	}

	handler := app.NewApp(cfg)
	router := chi.NewRouter()
	router.Use(chimiddleware.Logger)
	router.With(middleware.NewTimeout(cfg)).Post("/api/slow", handler.SlowHandler)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Printf("server is listening at port %s", cfg.Port)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	log.Printf("shutting down, waiting for %d seconds", cfg.ShutdownTimeout)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Duration(cfg.ShutdownTimeout)*time.Second,
	)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("service is down, exiting")
}
