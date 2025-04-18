package rest

import (
	ports "github.com/kotai-tech/server/internal/port/in"
	"github.com/labstack/echo/v4"
)

type PatientHandler struct {
	PatientService ports.PatientService
}

func SetHandler(e *echo.Echo, svc ports.PatientService) {
	p := &PatientHandler{
		PatientService: svc,
	}

	e.GET("/patients", p.getPatients)
}

func (h *PatientHandler) getPatients(context echo.Context) error {
	ctx := context.Request().Context()
	patients, err := h.PatientService.GetPatients(ctx)
	if err != nil {
		return context.JSON(500, err)
	}

	return context.JSON(200, patients)
}
