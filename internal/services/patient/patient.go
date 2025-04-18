package services

import (
	"context"

	domain "github.com/kotai-tech/server/internal/domain"
	in "github.com/kotai-tech/server/internal/port/in"
	out "github.com/kotai-tech/server/internal/port/out"
)

type Patient struct {
	out.PatientRepository
}

func NewService(repository out.PatientRepository) in.PatientService {
	return &Patient{
		PatientRepository: repository,
	}
}

func (p *Patient) GetPatients(ctx context.Context) ([]domain.Patient, error) {
	patients, err := p.PatientRepository.ListPatients(ctx)
	if err != nil {
		return nil, err
	}

	return patients, nil
}
