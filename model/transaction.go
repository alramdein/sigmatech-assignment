package model

type Transaction struct {
	ID                int     `json:"id"`
	CustomerID        int     `json:"customer_id"`
	ContractNumber    string  `json:"contract_number"`
	OTR               float64 `json:"otr"`
	AdminFee          float64 `json:"admin_fee"`
	InstallmentAmount float64 `json:"installment_amount"`
	InterestAmount    float64 `json:"interest_amount"`
	AssetName         string  `json:"asset_name"`
	Tenor             int     `json:"tenor"`
}
