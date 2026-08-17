# Osu-Phantom

A Go library that tracks a user's recent osu! scores and builds a fresh ranking from only the scores they set after tracking begins. Each tracked user is a `phantom.Client` that polls the osu! API for new plays and keeps a weighted top-100: scores are deduplicated per beatmap keeping the best, and total pp uses osu!'s weighting of 100%, 95%, 90.25% and so on down the list.

## PP calculation

Each score's pp is taken from the osu! API when present, otherwise computed locally by [crosu-pp](https://github.com/Minoxs/crosu-pp), a shared-library wrapper around [rosu-pp](https://github.com/MaxOhn/rosu-pp). The local path is cgo, so builds need cgo enabled and the `crosu_pp` shared library discoverable at runtime (on `PATH` on Windows).

## Rate limiting

Every request the package makes to the osu! API passes through a single process-wide pacer, so the whole program stays inside the osu! API terms of use regardless of how many clients run. The default is 60 requests per minute; change it with `osu.SetRateLimit`.

## Usage

Authorize with an OAuth client-credentials token through `osu.GetGuestToken`, then hand a `phantom.AuthProvider` to a client. `phantom.NewClient` builds a client for a known user id; `phantom.Login` looks the user up by name first. `Client.Update` fetches new scores, paging by offset while every score on a page is newer than the last seen, and folds them into the ranking. `Client.KeepUpdated` runs that on an interval until the user goes idle. `Client.Ranking`, `Client.GetTotalPP`, and `Client.Restore` read the ranking and rehydrate it from persisted scores.

## Building

```bash
CGO_ENABLED=1 go build ./...
go test ./...
```

Tests that exercise pp calculation need the `crosu_pp` shared library on `PATH`. See [crosu-pp](https://github.com/Minoxs/crosu-pp) for building it.
