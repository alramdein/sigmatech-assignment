package model

type Customer struct {
	ID          int     `json:"id"`
	NIK         string  `json:"nik"`
	Password    string  `json:"password,omitempty"`
	FullName    string  `json:"full_name"`
	LegalName   string  `json:"legal_name"`
	BirthPlace  string  `json:"birth_place"`
	BirthDate   string  `json:"birth_date"`
	Salary      float64 `json:"salary"`
	KTPPhoto    []byte  `json:"ktp_photo"`
	SelfiePhoto []byte  `json:"selfie_photo"`
}
