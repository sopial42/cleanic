package rest

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"

	"github.com/sopial42/cleanic/internal/adapters/rest/middleware"
	patient "github.com/sopial42/cleanic/internal/domains/patient"
	"github.com/sopial42/cleanic/internal/domains/user"
	patientSVC "github.com/sopial42/cleanic/internal/services/patient"
)

type PatientHandler struct {
	patientSVC.Service
}

func SetHandler(e *echo.Echo, service patientSVC.Service, accessMiddleware middleware.AuthAccessMiddleware) {
	p := &PatientHandler{
		service,
	}

	requireDoctor := accessMiddleware.RequireRoles(user.Roles{user.RoleDoctor})
	apiV1 := e.Group("/api/v1")
	{
		apiV1.GET("/patients", p.getPatients, requireDoctor)
		apiV1.GET("/patient/:id", p.getPatient, requireDoctor)
		apiV1.POST("/patient", p.createPatient, requireDoctor)
		apiV1.PATCH("/patient", p.updatePatient, requireDoctor)
		apiV1.DELETE("/patient/:id", p.deletePatient, requireDoctor)
	}
}

func (h *PatientHandler) getPatients(context echo.Context) error {
	ctx := context.Request().Context()
	patients, err := h.GetPatients(ctx)
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

	patient, err := h.GetPatientByID(ctx, patientID)
	if err != nil {
		log.Error("Error get patient: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusOK, patient)
}

func (h *PatientHandler) createPatient(context echo.Context) error {
	ctx := context.Request().Context()
	newPatient := new(patient.Patient)
	if err := context.Bind(newPatient); err != nil {
		log.Error("Error bind patient: ", err)
		return context.JSON(http.StatusBadRequest, err)
	}

	patientCreated, err := h.CreatePatient(ctx, *newPatient)
	if err != nil {
		log.Error("Error creating patient: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.JSON(http.StatusCreated, patientCreated)
}

func (h *PatientHandler) updatePatient(context echo.Context) error {
	ctx := context.Request().Context()
	patient := new(patient.Patient)
	if err := context.Bind(patient); err != nil {
		log.Error("Error bind patient: ", err)
		return context.JSON(http.StatusBadRequest, err)
	}

	patientUpdated, err := h.UpdatePatient(ctx, *patient)
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

	err = h.DeletePatient(ctx, patientID)
	if err != nil {
		log.Error("Error deleting patient: ", err)
		return context.JSON(http.StatusInternalServerError, err)
	}

	return context.NoContent(http.StatusNoContent)
}
