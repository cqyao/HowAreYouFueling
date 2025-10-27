package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"HowAreYouFueling/types"
	"HowAreYouFueling/util"
)

func GetAccessToken() (string, error) {
	authHeader := os.Getenv("AUTH_HEADER")

	url := "https://api.onegov.nsw.gov.au/oauth/client_credential/accesstoken?grant_type=client_credentials"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", authHeader)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, _ := io.ReadAll(res.Body)

	var access types.AccessResponse
	if err := json.Unmarshal(body, &access); err != nil {
		return "", err
	}

	issuedAtMs, _ := strconv.ParseInt(access.Issued, 10, 64)
	expiresInSec, _ := strconv.ParseInt(access.Expiry, 10, 64)

	expiry := time.UnixMilli(issuedAtMs).Add(time.Duration(expiresInSec) * time.Second)
	os.WriteFile("token_expiry.txt", []byte(expiry.Format("2006-01-02 15:04:05")), 0644)

	if err := util.UpdateEnvValue(".env", "ACCESS_TOKEN", access.AccessToken); err != nil {
		fmt.Println("Failed to update .env:", err)
	}

	return access.AccessToken, nil
}

func LoadToken() (string, error) {
	t := time.Now()
	data, err := os.ReadFile("token_expiry.txt")
	if err != nil {
		return "", err
	}
	if string(data) > t.Format("2006-01-02 15:04:05") {
		return os.Getenv("ACCESS_TOKEN"), nil
	}
	return GetAccessToken()
}
