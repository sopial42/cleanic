package main

import (
	"github.com/labstack/echo/v4"

	patientSVC "github.com/kotai-tech/server/internal/services/patient"
	"github.com/kotai-tech/server/internal/adapter/rest"
)


func main() {
	e := echo.New()

	patientService := patientSVC.NewService()
	rest.SetHandler(e, patientService)

	e.Logger.Fatal(e.Start(":8080"))
}
