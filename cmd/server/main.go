package main

import (
	"context"
	"goAuth/internal/storage"
	"goAuth/internal/handler"
	"goAuth/internal/repository"
	"goAuth/internal/router"
	"goAuth/internal/services"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)	

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, reading from system environment")
	}

	pool, err := storage.ConnectPostgres();
	if err != nil {
		log.Fatal("error from postgres db connection: ", err)
	}
	defer pool.Close()
	log.Println("✅ postgres connected successfully")

	redisClient, err := storage.ConnectRedis(context.Background())
	if err != nil {
		log.Fatal("failed to connect to redis: ", err)
	}
	defer redisClient.Close()
	log.Println("✅ redis connected successfully")

	PORT := os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("unable to find PORT")
	}

	userRepository := repository.NewPostgresRepository(pool)
	userService := services.NewUserService(userRepository)
	userHandler := handler.NewUserHandler(userService)
	app := router.SetUpRouter(userHandler, redisClient)

	srv := &http.Server {
		Addr: ":" + PORT,
		Handler: app,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 60 * time.Second,
	}

	go func() {
		log.Printf("goAuth server is running on PORT: %s", PORT)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed{
			log.Fatalf("goAuth server failed to start: %v\n", err)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	<-sig
	log.Println("shutdown signal received, draining connections...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed")
	}else {
		log.Println("goAuth server shut down cleanly")
	}
}