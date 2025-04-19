package rest

import (
	"fmt"
	"net/http"

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


	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/patients", p.getPatients)
		apiV1.GET("/patient/:id", p.getPatient)
		apiV1.POST("/patient", p.createPatient)
		apiV1.PATCH("/patient", p.updatePatient)
		apiV1.DELETE("/patient/:id", p.deletePatient)
	}
}

func (h *PatientHandler) getPatients(context echo.Context) error {
	ctx := context.Request().Context()
	patients, err := h.PatientService.GetPatients(ctx)
	if err != nil {
		log.Error("Error get patients: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, patients)
}

func (h *PatientHandler) getPatient(context echo.Context) error {
	ctx := context.Request().Context()

	id := context.Param("id")
	if id == "" {
		return context.JSON(http.StatusBadRequest, "id is required")
	}

	var patientID int64
	_, err := fmt.Sscanf(id, "%d", &patientID)
	if err != nil {
		log.Error("Error converting id to int64: ", err)
		return context.JSON(http.StatusBadRequest, "invalid id format")
	}

	patient, err := h.PatientService.GetPatientByID(ctx, patientID)
	if err != nil {
		log.Error("Error get patient: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) createPatient(context echo.Context) error {
	ctx := context.Request().Context()

	patient := new(domain.Patient)
	if err := context.Bind(patient); err != nil {
		log.Error("Error bind patient: ", err)
		return context.JSON(http.StatusBadRequest, err)
	}

	patientCreated, err := h.PatientService.CreatePatient(ctx, *patient)
	if err != nil {
		log.Error("Error creating patient: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusCreated, patientCreated)
}

func (h *PatientHandler) updatePatient(context echo.Context) error {
	ctx := context.Request().Context()

	patient := new(domain.Patient)
	if err := context.Bind(patient); err != nil {
		log.Error("Error bind patient: ", err)
		return context.JSON(http.StatusBadRequest, err)
	}

	patientUpdated, err := h.PatientService.UpdatePatient(ctx, *patient)
	if err != nil {
		log.Error("Error updating patient: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, patientUpdated)
}

func (h *PatientHandler) deletePatient(context echo.Context) error {
	ctx := context.Request().Context()

	id := context.Param("id")
	if id == "" {
		return context.JSON(http.StatusBadRequest, "id is required")
	}

	// Convert id to int64
	var patientID int64
	_, err := fmt.Sscanf(id, "%d", &patientID)
	if err != nil {
		log.Error("Error converting id to int64: ", err)
		return context.JSON(http.StatusBadRequest, "invalid id format")
	}

	err = h.PatientService.DeletePatient(ctx, patientID)
	if err != nil {
		log.Error("Error deleting patient: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.NoContent(http.StatusNoContent)
}
