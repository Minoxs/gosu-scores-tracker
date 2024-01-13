package osu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
	"io"
	"log"
	"net/http"
)

type Credentials struct {
	ClientID     int
	ClientSecret string
}

func buildOAUTHUrl(endpoint string) string {
	return fmt.Sprintf("%s/%s", OAuthURL, endpoint)
}

func APIv2URL(endpoint string) string {
	return fmt.Sprintf("%s/%s", ApiV2, endpoint)
}

func createRequestBody(i interface{}) io.Reader {
	data, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(data)
}

func GetGuestToken(c Credentials) (*GuestToken, error) {
	body := AuthGrant{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		GrantType:    "client_credentials",
		Scope:        "public",
	}
	r, err := apiClient.Post(buildOAUTHUrl("token"), JSON, createRequestBody(body))

	if err != nil {
		return nil, err
	}

	if r.StatusCode == 200 {
		res := &GuestToken{}
		return res, json.NewDecoder(r.Body).Decode(res)
	} else {
		return nil, errors.New("status_code=" + r.Status)
	}
}

func GetUserID(token *GuestToken, username string) (int, error) {
	endpoint := fmt.Sprintf("users/%s/osu/?key=username", username)

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))

	res, err := apiClient.Do(req)
	if err != nil {
		log.Printf("Error while sending request: %s\n", err)
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

func GetRecentScores(token *GuestToken, userid int) player.Scores {
	endpoint := fmt.Sprintf("users/%d/scores/recent/?mode=osu&limit=10", userid)
	// TODO CONFIGURE LIMIT
	// TODO ALLOW FOR SCORE PAGINATION

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))

	res, err := apiClient.Do(req)
	if err != nil {
		log.Printf("Error while sending request: %s\n", err)
		return nil
	}
	defer res.Body.Close()

	scores := make(player.Scores, 0)
	err = json.NewDecoder(res.Body).Decode(&scores)
	if err != nil {
		log.Printf("Error while decoding body: %s\n", err)
		return nil
	}

	return scores
}

func GetBeatmapScores(token *GuestToken, userID int, beatmapID int) player.Scores {
	endpoint := fmt.Sprintf("beatmaps/%d/scores/users/%d/all", beatmapID, userID)

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))

	res, err := apiClient.Do(req)
	if err != nil {
		log.Printf("Error while sending request: %s\n", err)
		return nil
	}
	defer res.Body.Close()

	s := struct{ Scores player.Scores }{make(player.Scores, 0)}
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		log.Printf("Error while decoding body: %s\n", err)
		return nil
	}

	return s.Scores
}

// TODO support mode
func DownloadBeatmap(id int) (buf []byte, err error) {
	if beatmap, found := cache.Get(id); found {
		return beatmap, nil
	}

	var url = BaseURL + "/osu/" + fmt.Sprintf("%d", id)
	var res *http.Response

	res, err = apiClient.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		err = errors.New("status_code=" + res.Status)
		return
	}

	buf, err = io.ReadAll(res.Body)
	log.Printf("BeatmapSize=%d bytes", len(buf))

	if err == nil {
		cache.Set(id, buf)
		log.Printf("Set into cache : CurrentSize=%d\n", cache.CurrentSize())
	}

	return
}
