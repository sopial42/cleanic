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

	persistence "github.com/kotai-tech/server/internal/adapters/persistence/patient"
	patientHTTPHandler "github.com/kotai-tech/server/internal/adapters/rest/patient"
	"github.com/kotai-tech/server/internal/config"
	patientSVC "github.com/kotai-tech/server/internal/services/patient"
)

func main() {
	config := config.Load()

	patientPersistence := persistence.NewPGClient(config.DBConfig)
	patientService := patientSVC.NewPatientService(patientPersistence)

	engine := echo.New()
	engine.Use(middleware.Logger())
	patientHTTPHandler.SetHandler(engine, patientService)

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
