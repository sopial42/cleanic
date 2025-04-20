package patient

import (
	"context"

	patient "github.com/kotai-tech/server/internal/domains/patient"
)

type patientService struct {
	persistence Persistence
}

func NewPatientService(persistence Persistence) Service {
	return &patientService{
		persistence: persistence,
	}
}

func (p *patientService) GetPatients(ctx context.Context) ([]patient.Patient, error) {
	patients, err := p.persistence.ListPatients(ctx)
	if err != nil {
		return nil, err
	}

	return patients, nil
}

func (p *patientService) GetPatientByID(ctx context.Context, id int64) (patient.Patient, error) {
	currentPatient, err := p.persistence.GetPatientByID(ctx, id)
	if err != nil {
		return patient.Patient{}, err
	}

	return currentPatient, nil
}

func (p *patientService) CreatePatient(ctx context.Context, inputPatient patient.Patient) (patient.Patient, error) {
	patientCreated, err := p.persistence.InsertPatient(ctx, inputPatient)
	if err != nil {
		return patient.Patient{}, err
	}

	return patientCreated, nil
}

func (p *patientService) UpdatePatient(ctx context.Context, inputPatient patient.Patient) (patient.Patient, error) {
	patientUpdated, err := p.persistence.UpdatePatient(ctx, inputPatient)
	if err != nil {
		return patient.Patient{}, err
	}

	return patientUpdated, nil
}

func (p *patientService) DeletePatient(ctx context.Context, id int64) error {
	err := p.persistence.DeletePatient(ctx, id)
	if err != nil {
		return err
	}

	return nil
}
