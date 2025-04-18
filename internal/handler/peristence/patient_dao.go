package persistence

import (
	domain "github.com/kotai-tech/server/internal/domain"
	"github.com/uptrace/bun"
)

type patientDAO struct {
	bun.BaseModel `bun:"table:patient"`
	ID            int64  `bun:"id"`
	Firstname     string `bun:"firstname"`
	Lastname      string `bun:"lastname"`
	Email         string `bun:"email"`
}

func (p *patientDAO) ToDomain() domain.Patient {
	return domain.Patient{
		ID:        p.ID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     p.Email,
	}
}
func FromDomain(p domain.Patient) patientDAO {
	return patientDAO{
		ID:        p.ID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     p.Email,
	}
}
func FromDomainList(p []domain.Patient) []patientDAO {
	var patientDAOs []patientDAO
	for _, patient := range p {
		patientDAOs = append(patientDAOs, FromDomain(patient))
	}
	return patientDAOs
}
func ToDomainList(p []patientDAO) []domain.Patient {
	var patients []domain.Patient
	for _, patient := range p {
		patients = append(patients, patient.ToDomain())
	}
	return patients
}
