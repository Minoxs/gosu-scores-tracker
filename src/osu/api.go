package osu

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const (
	OAUTH_URL = "https://osu.ppy.sh/oauth"
	API_V2    = "https://osu.ppy.sh/api/v2"

	GET  = "GET"
	POST = "POST"
	AUTH = "Authorization"
	JSON = "application/json"
)

var (
	CACHED_CLIENT = &http.Client{}
)

func buildOAUTHUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", OAUTH_URL, endpoint)
}

func APIv2URL(endpoint string) string {
	return fmt.Sprintf("%s/%s", API_V2, endpoint)
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

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
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

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
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

func GetBeatmapScores(token *GuestToken, userID int, beatmapID int) []Score {
	endpoint := fmt.Sprintf("beatmaps/%d/scores/users/%d/all", beatmapID, userID)

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))

	res, err := CACHED_CLIENT.Do(req)
	if err != nil {
		log.Printf("Error while getting scores: %v", err)
		return nil
	}
	defer res.Body.Close()

	s := struct{ Scores []Score }{make([]Score, 0)}
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		log.Printf("Error while decoding body: %s\n", err)
		return nil
	}
	return s.Scores
}
