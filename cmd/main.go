package main

import (
	"github.com/labstack/echo/v4"

	"github.com/kotai-tech/server/internal/config"
	persistence "github.com/kotai-tech/server/internal/handler/peristence"
	"github.com/kotai-tech/server/internal/handler/rest"
	patientSVC "github.com/kotai-tech/server/internal/services/patient"
)

func main() {
	config := config.Load()
	patientRepository := persistence.NewPGClient(config.DBConfig)
	patientService := patientSVC.NewService(patientRepository	)

	e := echo.New()
	rest.SetHandler(e, patientService)

	e.Logger.Fatal(e.Start(":8080"))
}
