package osu

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/minoxs/osu-phantom/pkg/osu/optimization"
	"github.com/minoxs/osu-phantom/pkg/osu/player"
	"io"
	"log/slog"
	"net/http"
)

type Credentials struct {
	ClientID     int
	ClientSecret string
}

// APIVersion is the osu! API v2 version this package requests, sent as the
// x-api-version header on every v2 call. osu-web gates response-shape changes on
// it and compares it as an integer, so pinning the newest published version keeps
// responses on the current solo_score shape. Update this when osu! ships a newer
// version whose shape this package has been adjusted to decode.
const APIVersion = "20241024"

// apiVersionHeader is the header name osu-web reads the version from.
const apiVersionHeader = "x-api-version"

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
	defer r.Body.Close()

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
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		return 0, err
	}
	defer res.Body.Close()

	body := &User{}
	err = json.NewDecoder(res.Body).Decode(body)
	if err != nil {
		return 0, err
	}
	return body.ID, err
}

func decodeProfile(res *http.Response) (*player.Profile, error) {
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return nil, ErrUserNotFound
	}
	if res.StatusCode != 200 {
		return nil, errors.New("status_code=" + res.Status)
	}

	profile := &player.Profile{}
	if err := json.NewDecoder(res.Body).Decode(profile); err != nil {
		return nil, err
	}
	return profile, nil
}

// GetUser fetches a full osu!standard profile by user id.
// Returns ErrUserNotFound when no user carries that id.
func GetUser(token *GuestToken, id int64) (*player.Profile, error) {
	req, _ := http.NewRequest(GET, APIv2URL(fmt.Sprintf("users/%d/osu", id)), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		return nil, err
	}
	return decodeProfile(res)
}

// GetUserByName fetches a full osu!standard profile by username.
// Returns ErrUserNotFound when no user carries that name.
func GetUserByName(token *GuestToken, username string) (*player.Profile, error) {
	req, _ := http.NewRequest(GET, APIv2URL(fmt.Sprintf("users/%s/osu?key=username", username)), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		return nil, err
	}
	return decodeProfile(res)
}

// GetRecentScores fetches one page of a user's recent osu!standard scores,
// newest first. limit is capped at 100 by the osu! API; offset pages past the
// newest results.
func GetRecentScores(token *GuestToken, userid, limit, offset int) player.Scores {
	endpoint := fmt.Sprintf("users/%d/scores/recent/?mode=osu&limit=%d&offset=%d", userid, limit, offset)

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		slog.Error("Error while sending request", "Error", err)
		return nil
	}
	defer res.Body.Close()

	scores := make(player.Scores, 0)
	err = json.NewDecoder(res.Body).Decode(&scores)
	if err != nil {
		slog.Error("Error while decoding response", "Error", err)
		return nil
	}

	return scores
}

func GetBeatmapScores(token *GuestToken, userID int, beatmapID int) player.Scores {
	endpoint := fmt.Sprintf("beatmaps/%d/scores/users/%d/all", beatmapID, userID)

	req, _ := http.NewRequest(GET, APIv2URL(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		slog.Error("Error while sending request", "Error", err)
		return nil
	}
	defer res.Body.Close()

	s := struct {
		Scores player.Scores `json:"scores"`
	}{}
	err = json.NewDecoder(res.Body).Decode(&s)
	if err != nil {
		slog.Error("Error while decoding response", "Error", err)
		return nil
	}

	return s.Scores
}

func DownloadBeatmap(id int64) (buf []byte, err error) {
	if beatmap, found := optimization.GetBeatmap(id); found {
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
	slog.Info("Beatmap downloaded", "ID", id, "Size", len(buf))
	if err == nil {
		optimization.PutBeatmap(id, buf)
	}

	return
}
