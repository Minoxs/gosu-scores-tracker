package osu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const OAUTH_URL = "https://osu.ppy.sh/oauth"
const BASE_URL = "https://osu.ppy.sh/api/v2"
const JSON = "application/json"

var CACHED_CLIENT = &http.Client{}

const (
	GET  = "GET"
	POST = "POST"
	AUTH = "Authorization"
)

func buildOAUTHUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", OAUTH_URL, endpoint)
}

func buildAPIUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", BASE_URL, endpoint)
}

func createRequestBody(i interface{}) io.Reader {
	data, _ := json.Marshal(i)
	return bytes.NewBuffer(data)
}

func GetGuestToken(clientID int, clientSecret string) *GuestToken {
	body := authGrant{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		GrantType:    "client_credentials",
		Scope:        "public",
	}
	r, err := http.Post(buildOAUTHUrl("token"), JSON, createRequestBody(body))

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if r.StatusCode == 200 {
		res := &GuestToken{}
		_ = json.NewDecoder(r.Body).Decode(res)
		return res
	} else {
		fmt.Printf("StatusCode=%d", r.StatusCode)
		return nil
	}
}

func GetUserID(token *GuestToken, username string) (int, error) {
	endpoint := fmt.Sprintf("users/%s/osu/?key=username", username)

	req, _ := http.NewRequest(GET, buildAPIUrl(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))

	res, err := CACHED_CLIENT.Do(req)
	if err != nil {
		log.Printf("Error while sending GetUserID request: %s\n", err)
		return 0, err
	}
	defer res.Body.Close()

	body := &User{}
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		log.Printf("Error while decoding body: %s\n", err)
		return 0, err
	}
	return body.ID, err
}

func GetRecentScores(token *GuestToken, userid int) []Score {
	endpoint := fmt.Sprintf("users/%d/scores/recent/?mode=osu&limit=25", userid)

	req, _ := http.NewRequest(GET, buildAPIUrl(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))

	res, err := CACHED_CLIENT.Do(req)
	if err != nil {
		log.Printf("Error while sending GetUserID request: %s\n", err)
		return nil
	}
	defer res.Body.Close()

	body := make([]Score, 0)
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		log.Printf("Error while decoding body: %s\n", err)
		return nil
	}
	return body
}
