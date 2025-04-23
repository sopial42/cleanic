package patient

type Patient struct {
	ID        ID     `json:"id"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Email     Email  `json:"email"`
}

type ID int64

type Email string
