package rest

import (
	domain "github.com/kotai-tech/server/internal/domain"
	ports "github.com/kotai-tech/server/internal/port/in"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type PatientHandler struct {
	PatientService ports.PatientService
}

func SetHandler(e *echo.Echo, svc ports.PatientService) {
	p := &PatientHandler{
		PatientService: svc,
	}

	e.GET("/patients", p.getPatients)
	e.POST("/patient", p.createPatient)
}

func (h *PatientHandler) getPatients(context echo.Context) error {
	ctx := context.Request().Context()
	patients, err := h.PatientService.GetPatients(ctx)
	if err != nil {
		log.Error("Error get patient: ", err)
		return context.JSON(500, err)
	}

	return context.JSON(200, patients)
}

func (h *PatientHandler) createPatient(context echo.Context) error {
	ctx := context.Request().Context()

	patient := new(domain.Patient)
	if err := context.Bind(patient); err != nil {
		log.Error("Error bind patient: ", err)
		return context.JSON(400, err)
	}

	patientCreated, err := h.PatientService.CreatePatient(ctx, *patient)
	if err != nil {
		log.Error("Error creating patient: ", err)
		return context.JSON(500, err)
	}

	return context.JSON(201, patientCreated)
}
