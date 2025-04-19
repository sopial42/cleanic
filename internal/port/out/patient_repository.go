package repository

import (
	"context"

	domain "github.com/kotai-tech/server/internal/domain"
)

type PatientRepository interface {
	ListPatients(ctx context.Context) ([]domain.Patient, error)
	GetPatientByID(ctx context.Context, id int64) (domain.Patient, error)
	InsertPatient(ctx context.Context, patient domain.Patient) (domain.Patient, error)
	UpdatePatient(ctx context.Context, patient domain.Patient) (domain.Patient, error)
	DeletePatient(ctx context.Context, id int64) error
}
