package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	userCLI "github.com/sopial42/cleanic/internal/adapters/clients/user"
	persistence "github.com/sopial42/cleanic/internal/adapters/persistence"
	authPersistence "github.com/sopial42/cleanic/internal/adapters/persistence/auth"
	patientPersistence "github.com/sopial42/cleanic/internal/adapters/persistence/patient"
	userPersistence "github.com/sopial42/cleanic/internal/adapters/persistence/user"
	authHTTPHandler "github.com/sopial42/cleanic/internal/adapters/rest/auth"
	authMiddleware "github.com/sopial42/cleanic/internal/adapters/rest/middleware"
	patientHTTPHandler "github.com/sopial42/cleanic/internal/adapters/rest/patient"
	userHTTPHandler "github.com/sopial42/cleanic/internal/adapters/rest/user"
	"github.com/sopial42/cleanic/internal/config"
	authSVC "github.com/sopial42/cleanic/internal/services/auth"
	patientSVC "github.com/sopial42/cleanic/internal/services/patient"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
)

func main() {
	config := config.Load()
	pgClient := persistence.NewPGClient(config.DB)

	refreshMiddleware := authMiddleware.NewAuthRefreshMiddleware(config.JWT.RefreshTokenConfig)
	accessMiddleware := authMiddleware.NewAuthAccessMiddleware(config.JWT.AccessTokenConfig)
	userPersistence := userPersistence.NewPGClient(pgClient)
	authPersistence := authPersistence.NewPGClient(pgClient)
	userService := userSVC.NewUserService(userPersistence)

	userClient := userCLI.NewInMemoryUserClient(userService)
	authService := authSVC.NewAuthService(userClient, config.JWT, authPersistence)

	patientPersistence := patientPersistence.NewPGClient(pgClient)
	patientService := patientSVC.NewPatientService(patientPersistence)

	engine := echo.New()
	engine.Use(middleware.Logger())
	engine.Use(session.Middleware(sessions.NewCookieStore(config.JWT.CookieStoreConfig.Secret)))

	patientHTTPHandler.SetHandler(engine, patientService, accessMiddleware)
	userHTTPHandler.SetHandler(engine, userService, accessMiddleware)
	authHTTPHandler.SetHandler(engine, config.JWT.CookieStoreConfig, authService, refreshMiddleware)

	go func() {
		if err := engine.Start(":8080"); err != nil {
			log.Printf("Shutting down the server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down the server gracefully")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := engine.Shutdown(ctx); err != nil {
		log.Printf("Unable to shutdown server gracefully: %v\n", err)
		return
	}
	log.Println("Server has shut down gracefully")
}
