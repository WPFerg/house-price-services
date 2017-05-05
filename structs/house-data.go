package structs

type HouseData struct {
	Id                string   `json:"id"`
	Date              string   `json:"date"`
	Postcode          string   `json:"postcode"`
	FlagA             string   `json:"flagA"`
	FlagB             string   `json:"flagB"`
	FlagC             string   `json:"flagC"`
	HouseNameOrNumber string   `json:"houseNameOrNumber"`
	AdditionalNumber  string   `json:"additionalNumber"`
	Address           []string `json:"address"`
	County            string   `json:"county"`
	Cost              int      `json:"cost"`
}
