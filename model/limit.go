package model

type Limit struct {
	ID         int     `json:"id"`
	CustomerID int     `json:"customer_id"`
	Tenor1     float64 `json:"tenor_1"`
	Tenor2     float64 `json:"tenor_2"`
	Tenor3     float64 `json:"tenor_3"`
	Tenor4     float64 `json:"tenor_4"`
}
