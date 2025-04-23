package patient

import (
	"context"

	patient "github.com/sopial42/cleanic/internal/domains/patient"
)

type Service interface {
	CreatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	GetPatients(ctx context.Context) ([]patient.Patient, error)
	GetPatientByID(ctx context.Context, id int64) (patient.Patient, error)
	UpdatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	DeletePatient(ctx context.Context, id int64) error
}

type Persistence interface {
	ListPatients(ctx context.Context) ([]patient.Patient, error)
	GetPatientByID(ctx context.Context, id int64) (patient.Patient, error)
	InsertPatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	UpdatePatient(ctx context.Context, patient patient.Patient) (patient.Patient, error)
	DeletePatient(ctx context.Context, id int64) error
}
