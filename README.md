# gosu-scores-tracker

A Go library that tracks osu! players' scores as they set them and builds a fresh ranking from only the scores set after tracking begins. It sits on top of [gosu-api](https://github.com/Minoxs/gosu-api), which owns osu! API access and the weighted ranking; this library adds the tracking layer: the global feed tail, per-user polling, and the accumulation into a live ranking. A tracked user's ranking is a weighted top-100, deduplicated per beatmap keeping the best, with total pp using osu!'s weighting of 100%, 95%, 90.25% and so on down the list.

## What it provides

- `tracker.Client` polls a single user's recent scores and folds new ones into a `gosu.Ranking`. `NewClient` builds one for a known user id; `Login` looks the user up by name first.
- `RealtimePoller` tails osu!'s global recent-scores feed, fanning every passing score out through an embedded `Broadcaster`. The feed's scores are lean: a beatmap id with no embedded beatmap.
- `FilterTracker` narrows a feed to a chosen set of users from a start time; `UserPoller` does the same by polling the per-user endpoint, which returns scores with their maps embedded.

## PP

Each score's pp is whatever the osu! API returns; nothing here computes it. The API omits pp for a score that earns none (unranked mods, an unranked map, or one osu! has not finished processing), and that decodes as zero. A zero-pp score contributes nothing to the ranking; any further policy on it is the caller's.

## Rate limiting

Every osu! request goes through a `gosu.Client`, and every client reserves a slot on a single process-wide pacer, so the whole program stays inside the osu! API terms of use regardless of how many clients run. The default is 60 requests per minute; change it with `gosu.SetRateLimit`. A client is built at a priority with `gosu.NewClient(prio)`: when several requests wait at once the pacer grants the higher priority first, so a burst of background traffic never delays a higher-priority request by more than one slot. The priority is the client's, chosen once when it is built; the request methods carry none of their own.

## Usage

Authorize with an OAuth client-credentials token through `gosu.NewClient(prio).GetGuestToken`, then hand a `tracker.AuthProvider` to a poller or client. For a single user, `Client.Update` fetches new scores, paging by offset while every score on a page is newer than the last seen, and folds them into the ranking; `Client.KeepUpdated` runs that on an interval until the user goes idle. For a population, build a `RealtimePoller` and subscribe, or a `UserPoller`/`FilterTracker` over a tracked set. `Client.Ranking`, `Client.GetTotalPP`, and `Client.Restore` read the ranking and rehydrate it from persisted scores.

## Building

The module is pure Go, with no cgo and no native dependency.

```bash
go build ./...
go test ./...
```
