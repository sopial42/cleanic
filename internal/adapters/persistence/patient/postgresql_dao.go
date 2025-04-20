package persistence

import (
	"github.com/uptrace/bun"

	patient "github.com/kotai-tech/server/internal/domains/patient"
)

type patientDAO struct {
	bun.BaseModel `bun:"table:patient"`
	ID            int64  `bun:"id,autoincrement"`
	Firstname     string `bun:"firstname"`
	Lastname      string `bun:"lastname"`
	Email         string `bun:"email"`
}

func patientFromDAOToDomain(p patientDAO) patient.Patient {
	return patient.Patient{
		ID:        p.ID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     p.Email,
	}
}
func patientFromDomainToDAO(p patient.Patient) patientDAO {
	return patientDAO{
		ID:        p.ID,
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     p.Email,
	}
}

func patientFromDAOsToDomains(pDAOs []patientDAO) []patient.Patient {
	var patients []patient.Patient
	for _, pDAO := range pDAOs {
		patients = append(patients, patientFromDAOToDomain(pDAO))
	}
	return patients
}

// func patientFromDomainsToDAOs(p []patient.Patient) []patientDAO {
// 	var patientDAOs []patientDAO
// 	for _, patient := range p {
// 		patientDAOs = append(patientDAOs, patientFromDomainToDAO(patient))
// 	}
// 	return patientDAOs
// }
