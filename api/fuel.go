package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"HowAreYouFueling/types"
)

func GetFuelPrices(token, apiKey string, payload types.Payload) (*types.Response, error) {
	url := "https://api.onegov.nsw.gov.au/FuelPriceCheck/v1/fuel/prices/location"
	jsonBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", url, strings.NewReader(string(jsonBytes)))

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("authorization", "Bearer "+token)
	req.Header.Add("apikey", apiKey)
	req.Header.Add("transactionid", fmt.Sprintf("%d", time.Now().UnixNano()))
	req.Header.Add("requesttimestamp", time.Now().Format("02/01/2006 03:04:05 PM"))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)
	var resp types.Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	fmt.Println(&resp)
	return &resp, nil
}
