package ELKauth

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	tokenURL     = "https://lk.orb.ru/oauth/token"
	clientID     = "27"
	clientSecret = "SuDGnNtf0MBN9SE6YXyvUID5VluDYuJC4qwtN7pK"
	redirectURI  = "http://hackathon-2.orb.ru/login-elk-redirect"
	grantType    = "authorization_code"
	userDataURL  = "https://lk.orb.ru/api/get_user"
	scope        = "email+session+auth_method"
)

func ExchangeCodeForToken(code string) (*TokenResponse, error) {
	reqBody := TokenRequest{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURI:  redirectURI,
		Code:         code,
		GrantType:    grantType,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", tokenURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to exchange token, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tokenResponse TokenResponse
	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return nil, err
	}

	return &tokenResponse, nil
}

func GetUserData(accessToken string) (*UserData, error) {
	params := url.Values{}
	params.Set("scope", scope)

	requestURL := fmt.Sprintf("%s?%s", userDataURL, params.Encode())

	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to fetch user data, status: %d, response: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var userData UserData
	if err := json.Unmarshal(body, &userData); err != nil {
		return nil, err
	}

	return &userData, nil
}
