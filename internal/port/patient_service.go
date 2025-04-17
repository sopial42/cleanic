package repository

import (
	"context"

	model "github.com/kotai-tech/server/internal/domain"
)

type PatientService interface {
	GetPatients(ctx context.Context) ([]model.Patient, error)
}
