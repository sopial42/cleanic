package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	persistence "github.com/sopial42/cleanic/internal/adapters/persistence"
	patientPersistence "github.com/sopial42/cleanic/internal/adapters/persistence/patient"
	userPersistence "github.com/sopial42/cleanic/internal/adapters/persistence/user"
	patientHTTPHandler "github.com/sopial42/cleanic/internal/adapters/rest/patient"
	userHTTPHandler "github.com/sopial42/cleanic/internal/adapters/rest/user"
	"github.com/sopial42/cleanic/internal/config"
	patientSVC "github.com/sopial42/cleanic/internal/services/patient"
	userSVC "github.com/sopial42/cleanic/internal/services/user"
)

func main() {
	config := config.Load()
	pgClient := persistence.NewPGClient(config.DBConfig)

	patientPersistence := patientPersistence.NewPGClient(pgClient)
	patientService := patientSVC.NewPatientService(patientPersistence)

	userPersistence := userPersistence.NewPGClient(pgClient)
	userService := userSVC.NewUserService(userPersistence)

	engine := echo.New()
	engine.Use(middleware.Logger())

	patientHTTPHandler.SetHandler(engine, patientService)
	userHTTPHandler.SetHandler(engine, userService)

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
