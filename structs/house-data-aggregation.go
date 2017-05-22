package structs

type HouseDataAggregation struct {
	ID      string  `json:"id"`
	Min     int     `json:"min"`
	Max     int     `json:"max"`
	Average float32 `json:"average"`
	Total   int     `json:"total"`
	Count   int     `json:"count"`
}
