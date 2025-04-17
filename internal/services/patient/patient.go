package services

import (
	"context"

	domain "github.com/kotai-tech/server/internal/domain"
	ports "github.com/kotai-tech/server/internal/port"
)

type Patient struct{}

func NewService() ports.PatientService {
	return &Patient{}
}

func (p *Patient) GetPatients(ctx context.Context) ([]domain.Patient, error) {
	// This is a stub implementation. Replace with actual database call.
	return []domain.Patient{
		{Name: "John Doe", Email: "jd@gmail.com"},
	}, nil
}
