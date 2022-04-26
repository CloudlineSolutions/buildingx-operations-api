package buildingx

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type AuthRequest struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Audience     string `json:"audience"`
	GrantType    string `json:"grant_type"`
	URL          string `json:"url"`
}
type SBToken struct {
	AccessToken string `json:"access_token"`
	Expiration  int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

func GetToken() (string, error) {

	// verify that you have the environment variables needed
	clientID := os.Getenv("BUILDINGX_CLIENT_ID")
	if clientID == "" {
		return "", errors.New("missing client id")
	}
	clientSecret := os.Getenv("BUILDINGX_CLIENT_SECRET")
	if clientSecret == "" {
		return "", errors.New("missing client secret")
	}
	audience := os.Getenv("BUILDINGX_AUDIENCE")
	if audience == "" {
		return "", errors.New("missing audience")
	}
	authURL := os.Getenv("BUILDINGX_AUTH_URL")
	if authURL == "" {
		return "", errors.New("missing authorization URL")
	}

	authRequest := AuthRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Audience:     audience,
		GrantType:    "client_credentials",
		URL:          authURL,
	}
	authRequestBytes, _ := json.Marshal(authRequest)
	authRequestReader := bytes.NewReader(authRequestBytes)

	client := &http.Client{Timeout: time.Duration(20) * time.Second}
	req, err := http.NewRequest("POST", authRequest.URL, authRequestReader)
	if err != nil {
		return "", err

	}
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("unexpected error while invoking http client: %s", err.Error())
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err := decoder.Decode(&data)
		if err == nil {
			return "", errors.New("error message from BuildingX: " + data["detail"].(string))
		}
		return "", errors.New("got non-200 response from BuildingX API with no additional information")

	}
	tkn := SBToken{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &tkn); err != nil {
		return "", errors.New("Error parsing API response. String submitted: " + string(body))
	}
	return string(tkn.AccessToken), nil

}
