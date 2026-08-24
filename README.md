# Osu-Phantom

A Go library that tracks a user's recent osu! scores and builds a fresh ranking from only the scores they set after tracking begins. Each tracked user is a `phantom.Client` that polls the osu! API for new plays and keeps a weighted top-100: scores are deduplicated per beatmap keeping the best, and total pp uses osu!'s weighting of 100%, 95%, 90.25% and so on down the list.

## PP

Each score's pp is whatever the osu! API returns; the library computes none. The API omits pp for a score that earns none (unranked mods, an unranked map, or one osu! has not finished processing), and that decodes as zero. A zero-pp score contributes nothing to the ranking; any further policy on it is the caller's.

## Rate limiting

Every osu! API request goes through an `osu.Client`, and every client reserves a slot on a single process-wide pacer, so the whole program stays inside the osu! API terms of use regardless of how many clients run. The default is 60 requests per minute; change it with `osu.SetRateLimit`. A client is built at a priority with `osu.NewClient(prio)`: when several requests wait at once the pacer grants the higher priority first, so a burst of background traffic never delays a higher-priority request by more than one slot. The priority is the client's, chosen once when it is built; the request methods carry none of their own. The caller assigns the levels their meaning.

## Usage

Authorize with an OAuth client-credentials token through `osu.NewClient(prio).GetGuestToken`, then hand a `phantom.AuthProvider` to a client. `phantom.NewClient` builds a client for a known user id; `phantom.Login` looks the user up by name first. `Client.Update` fetches new scores, paging by offset while every score on a page is newer than the last seen, and folds them into the ranking. `Client.KeepUpdated` runs that on an interval until the user goes idle. `Client.Ranking`, `Client.GetTotalPP`, and `Client.Restore` read the ranking and rehydrate it from persisted scores.

## Building

The module is pure Go, with no cgo and no native dependency.

```bash
go build ./...
go test ./...
```
