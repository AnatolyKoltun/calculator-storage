package models

import "time"

type Calculation struct {
	ID        int       `json:"id"`
	Argument1 float64   `json:"argument1"`
	Argument2 float64   `json:"argument2"`
	Operator  string    `json:"operator"`
	Result    float64   `json:"result"`
	CreatedAt time.Time `json:"created_at"`
}

type FilterRequest struct {
	DateFrom string `form:"date_from" time_format:"2006-01-02"`
	DateTo   string `form:"date_to"   time_format:"2006-01-02"`
}
