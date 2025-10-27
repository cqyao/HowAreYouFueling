package types

type Response struct {
	Stations []struct {
		Brand   string `json:"brand"`
		Code    string `json:"code"`
		Name    string `json:"name"`
		Address string `json:"address"`
	} `json:"stations"`
	Prices []struct {
		StationCode string  `json:"stationcode"`
		FuelType    string  `json:"fueltype"`
		Price       float64 `json:"price"`
		LastUpdated string  `json:"lastupdate"`
	} `json:"prices"`
}

type AccessResponse struct {
	AccessToken string `json:"access_token"`
	Expiry      string `json:"expires_in"`
	Issued      string `json:"issued_at"`
}

type Payload struct {
	FuelType      string   `json:"fueltype"`
	Brand         []string `json:"brand"`
	NamedLocation string   `json:"namedlocation"`
	Reference     struct {
		Latitude  string `json:"latitude"`
		Longitude string `json:"longitude"`
	} `json:"referencepoint"`
	SortBy        string `json:"sortby"`
	SortAscending string `json:"sortascending"`
}
