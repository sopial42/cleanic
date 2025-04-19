package services

import (
	"context"

	domain "github.com/kotai-tech/server/internal/domain"
	in "github.com/kotai-tech/server/internal/port/in"
	out "github.com/kotai-tech/server/internal/port/out"
)

type Patient struct {
	repository out.PatientRepository
}

func NewService(repository out.PatientRepository) in.PatientService {
	return &Patient{
		repository: repository,
	}
}

func (p *Patient) GetPatients(ctx context.Context) ([]domain.Patient, error) {
	patients, err := p.repository.ListPatients(ctx)
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (p *Patient) GetPatientByID(ctx context.Context, id int64) (domain.Patient, error) {
	patient, err := p.repository.GetPatientByID(ctx, id)
	if err != nil {
		return domain.Patient{}, err
	}

	return patient, nil
}

func (p *Patient) CreatePatient(ctx context.Context, patient domain.Patient) (domain.Patient, error) {
	patientCreated, err := p.repository.InsertPatient(ctx, patient)
	if err != nil {
		return domain.Patient{}, err
	}

	return patientCreated, nil
}

func (p *Patient) UpdatePatient(ctx context.Context, patient domain.Patient) (domain.Patient, error) {
	patientUpdated, err := p.repository.UpdatePatient(ctx, patient)
	if err != nil {
		return domain.Patient{}, err
	}

	return patientUpdated, nil
}

func (p *Patient) DeletePatient(ctx context.Context, id int64) error {
	err := p.repository.DeletePatient(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
