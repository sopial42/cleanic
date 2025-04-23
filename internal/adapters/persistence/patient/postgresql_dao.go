package persistence

import (
	"github.com/uptrace/bun"

	"github.com/sopial42/cleanic/internal/domains/patient"
)

type patientDAO struct {
	bun.BaseModel `bun:"table:patient"`
	ID            int64     `bun:"id,pk,autoincrement"`
	Firstname     string    `bun:"firstname"`
	Lastname      string    `bun:"lastname"`
	Email         string    `bun:"email"`
}

func patientFromDAOToDomain(p patientDAO) patient.Patient {
	return patient.Patient{
		ID:        patient.ID(p.ID),
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     patient.Email(p.Email),
	}
}
func patientFromDomainToDAO(p patient.Patient) patientDAO {
	return patientDAO{
		ID:        int64(p.ID),
		Firstname: p.Firstname,
		Lastname:  p.Lastname,
		Email:     string(p.Email),
	}
}

func patientFromDAOsToDomains(pDAOs []patientDAO) []patient.Patient {
	patients := make([]patient.Patient, len(pDAOs))
	for i, pDAO := range pDAOs {
		patients[i] = patientFromDAOToDomain(pDAO)
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
