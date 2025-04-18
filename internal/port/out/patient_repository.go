package repository

import (
	"context"

	domain "github.com/kotai-tech/server/internal/domain"
)

type PatientRepository interface {
	ListPatients(ctx context.Context) ([]domain.Patient, error)
}
