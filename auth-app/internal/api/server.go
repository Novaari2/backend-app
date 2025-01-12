package api

import (
	"auth-app/internal/config"
	"auth-app/internal/entity"
	"auth-app/internal/users"
	"auth-app/internal/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewServer() *server {
	cfg := config.Load()

	dsn := fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.DbName,
		cfg.Database.Password,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})

	if err != nil {
		log.Fatal().Err(err).Msg("Unable to open database connection")
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to get sql db")
	}

	sqlDB.SetMaxIdleConns(5)
	sqlDB.SetMaxOpenConns(10)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = db.AutoMigrate(&entity.User{})
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to migrate database")
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	passwordGenerator := utils.GeneratePassword
	userReposiory := users.NewRepository(db)
	userService := users.NewService(userReposiory, passwordGenerator)
	userHandler := users.NewHTTPHandler(userService)

	r.Use(cors.New(cors.Options{
		AllowedOrigins:     []string{"*"},
		AllowedMethods:     []string{"POST", "GET", "PUT", "DELETE", "HEAD", "OPTIONS"},
		AllowedHeaders:     []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		MaxAge:             60,
		AllowCredentials:   true,
		OptionsPassthrough: false,
		Debug:              false,
	}).Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", userHandler.RegisterUserHandler)
	})

	return &server{router: r}
}

type server struct {
	router chi.Router
}

func (s *server) Run(port int) {
	addr := fmt.Sprintf(":%d", port)

	httpSrv := &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	done := make(chan bool)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-quit
		log.Info().Msg("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		httpSrv.SetKeepAlivesEnabled(false)
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Fatal().Msgf("Could not gracefully shutdown the server: %v", err)
		}
		close(done)
	}()

	log.Info().Msgf("Listening and serving on port %d", port)

	if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal().Msgf("Could not listen on %s: %v", addr, err)
	}

	<-done
	log.Info().Msg("Server stopped")
}
