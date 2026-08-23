package osu

import (
	"bytes"
	"context"
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

func createRequestBody(i any) io.Reader {
	data, err := json.Marshal(i)
	if err != nil {
		panic(err)
	}
	return bytes.NewBuffer(data)
}

func GetGuestToken(ctx context.Context, c Credentials) (*GuestToken, error) {
	body := AuthGrant{
		ClientID:     c.ClientID,
		ClientSecret: c.ClientSecret,
		GrantType:    "client_credentials",
		Scope:        "public",
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, buildOAUTHUrl("token"), createRequestBody(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", JSON)
	r, err := apiClient.Do(req)

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

// GetUser fetches a full osu!standard profile by user id, at the priority carried
// on ctx. Returns ErrUserNotFound when no user carries that id.
func GetUser(ctx context.Context, token *GuestToken, id int64) (*player.Profile, error) {
	req, _ := http.NewRequestWithContext(ctx, GET, APIv2URL(fmt.Sprintf("users/%d/osu", id)), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		return nil, err
	}
	return decodeProfile(res)
}

// GetUserByName fetches a full osu!standard profile by username, at the priority
// carried on ctx. Returns ErrUserNotFound when no user carries that name.
func GetUserByName(ctx context.Context, token *GuestToken, username string) (*player.Profile, error) {
	req, _ := http.NewRequestWithContext(ctx, GET, APIv2URL(fmt.Sprintf("users/%s/osu?key=username", username)), nil)
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

// GetScores fetches one page of osu!'s global recent-scores feed for a ruleset,
// the passing scores every player has set, ascending by id. cursor is the
// cursor_string from a previous page, or empty for the newest page; the returned
// cursor_string fetches the scores newer than this page. Scores carry no embedded
// beatmap, only a beatmap_id.
func GetScores(ctx context.Context, token *GuestToken, ruleset, cursor string) (player.Scores, string, error) {
	endpoint := "scores?ruleset=" + ruleset
	if cursor != "" {
		endpoint += "&cursor_string=" + cursor
	}

	req, _ := http.NewRequestWithContext(ctx, GET, APIv2URL(endpoint), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, "", errors.New("status_code=" + res.Status)
	}

	page := struct {
		Scores       player.Scores `json:"scores"`
		CursorString string        `json:"cursor_string"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&page); err != nil {
		return nil, "", err
	}
	return page.Scores, page.CursorString, nil
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

// GetBeatmap fetches a single beatmap's metadata by id. The osu! API nests the
// owning beatmapset in the response, so both are returned: the map for its status
// and difficulty, the set for its title, artist, and cover art. Unlike the beatmap
// embedded in a score, this response carries a real max_combo.
func GetBeatmap(ctx context.Context, token *GuestToken, id int64) (player.Beatmap, player.BeatmapSet, error) {
	req, _ := http.NewRequestWithContext(ctx, GET, APIv2URL(fmt.Sprintf("beatmaps/%d", id)), nil)
	req.Header.Add(AUTH, createHeader(token.TokenType, token.AccessToken))
	req.Header.Add(apiVersionHeader, APIVersion)

	res, err := apiClient.Do(req)
	if err != nil {
		return player.Beatmap{}, player.BeatmapSet{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return player.Beatmap{}, player.BeatmapSet{}, errors.New("status_code=" + res.Status)
	}

	body := struct {
		player.Beatmap
		BeatmapSet player.BeatmapSet `json:"beatmapset"`
	}{}
	if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
		return player.Beatmap{}, player.BeatmapSet{}, err
	}
	return body.Beatmap, body.BeatmapSet, nil
}

func DownloadBeatmap(ctx context.Context, id int64) (buf []byte, err error) {
	if beatmap, found := optimization.GetBeatmap(id); found {
		return beatmap, nil
	}

	var url = BaseURL + "/osu/" + fmt.Sprintf("%d", id)
	var res *http.Response

	req, err := http.NewRequestWithContext(ctx, GET, url, nil)
	if err != nil {
		return nil, err
	}
	res, err = apiClient.Do(req)
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
