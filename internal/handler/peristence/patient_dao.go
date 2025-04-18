package persistence

import (
	domain "github.com/kotai-tech/server/internal/domain"
	"github.com/uptrace/bun"
)

type patientDAO struct {
	bun.BaseModel `bun:"table:patient"`
	ID            int64  `bun:"id,autoincrement"`
	Firstname     string `bun:"firstname"`
	Lastname      string `bun:"lastname"`
	Email         string `bun:"email"`
}

func patientFromDAOToDomain(p patientDAO) domain.Patient {
	return domain.Patient{
		ID:        p.ID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     p.Email,
	}
}
func patientFromDomainToDAO(p domain.Patient) patientDAO {
	return patientDAO{
		ID:        p.ID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     p.Email,
	}
}
func patientFromDomainsToDAOs(p []domain.Patient) []patientDAO {
	var patientDAOs []patientDAO
	for _, patient := range p {
		patientDAOs = append(patientDAOs, patientFromDomainToDAO(patient))
	}
	return patientDAOs
}

func patientFromDAOsToDomains(pDAOs []patientDAO) []domain.Patient {
	var patients []domain.Patient
	for _, pDAO := range pDAOs {
		patients = append(patients, patientFromDAOToDomain(pDAO))
	}
	return patients
}
