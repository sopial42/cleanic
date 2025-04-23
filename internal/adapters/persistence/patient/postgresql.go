package persistence

import (
	"context"
	"fmt"

	"github.com/uptrace/bun"

	patient "github.com/sopial42/cleanic/internal/domains/patient"
	patientSVC "github.com/sopial42/cleanic/internal/services/patient"
)

type pgPersistence struct {
	clientDB *bun.DB
}

func NewPGClient(client *bun.DB) patientSVC.Persistence {
	return &pgPersistence{clientDB: client}
}

func (p *pgPersistence) InsertPatient(ctx context.Context, newPatient patient.Patient) (patient.Patient, error) {
	patientDAO := patientFromDomainToDAO(newPatient)

	_, err := p.clientDB.NewInsert().
		Model(&patientDAO).
		Returning("*").
		Exec(ctx)
	if err != nil {
		return patient.Patient{}, fmt.Errorf("err: %w", err)
	}

	// ID == 0 means that the insert failed
	if patientDAO.ID == 0 {
		return patient.Patient{}, fmt.Errorf("unable to create a new patient: %v", patientDAO)
	}

	return patientFromDAOToDomain(patientDAO), nil
}

func (p *pgPersistence) ListPatients(ctx context.Context) ([]patient.Patient, error) {
	var patientDAOs []patientDAO

	request := p.clientDB.NewSelect().Model(&patientDAOs)
	err := request.Scan(ctx)
	if err != nil {
		return nil, fmt.Errorf("err: %w", err)
	}

	return patientFromDAOsToDomains(patientDAOs), nil
}

func (p *pgPersistence) GetPatientByID(ctx context.Context, id int64) (patient.Patient, error) {
	var patientDAO patientDAO

	err := p.clientDB.NewSelect().
		Model(&patientDAO).
		Where("id = ?", id).
		Scan(ctx)
	if err != nil {
		return patient.Patient{}, fmt.Errorf("err: %w", err)
	}

	if patientDAO.ID == 0 {
		return patient.Patient{}, fmt.Errorf("patient not found with id: %d", id)
	}

	return patientFromDAOToDomain(patientDAO), nil
}

func (p *pgPersistence) UpdatePatient(ctx context.Context, updatedPatient patient.Patient) (patient.Patient, error) {
	patientDAO := patientFromDomainToDAO(updatedPatient)
	if patientDAO.ID == 0 {
		return patient.Patient{}, fmt.Errorf("unable to update any patient as ID is 0: %+v", patientDAO)
	}

	_, err := p.clientDB.NewUpdate().
		Model(&patientDAO).
		Where("id = ?", updatedPatient.ID).
		OmitZero().
		Returning("*").
		Exec(ctx)
	if err != nil {
		return patient.Patient{}, fmt.Errorf("unable to request patient update: %w", err)
	}

	return patientFromDAOToDomain(patientDAO), nil
}

func (p *pgPersistence) DeletePatient(ctx context.Context, id int64) error {
	_, err := p.clientDB.NewDelete().
		Model((*patientDAO)(nil)).
		Where("id = ?", id).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("unable to delete patient id: %d, err: %w", id, err)
	}

	return nil
}
