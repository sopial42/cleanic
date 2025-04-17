package rest

import (
	ports "github.com/kotai-tech/server/internal/port"
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

func (h *PatientHandler) getPatients(c echo.Context) error {
	ctx := c.Request().Context()
	patients, err := h.PatientService.GetPatients(ctx)
	if err != nil {
		return c.JSON(500, err)
	}
	return c.JSON(200, patients)
}
