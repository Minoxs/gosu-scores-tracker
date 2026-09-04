# gosu-scores-tracker

Score tracking plugin for [gosu-api](https://github.com/Minoxs/gosu-api).

## Install

```bash
go get github.com/minoxs/gosu-scores-tracker
```

## Usage

Both pollers run off your `gosu.App` and share its rate ceiling. Subscribe, run in a goroutine, close the channel to unsubscribe. Drain it or it blocks.

### Players

This mode of tracking is what I initially came up with, but it doesn't scale well for many users.
Since this endpoints gives you a full beatmap for free, it works great for tools that calculate pp and stuff like that.
Important to note that this method also gives you **unranked**, **loved** and **failed** scores.

```go
var poller = tracker.NewUserPoller(app, tracker.PollConfig{})
var scores = poller.Subscribe()
poller.Track(peppyID, time.Now()) 

go poller.Run(ctx)
for s := range scores {
	fmt.Println(s.UserID, s.Beatmapset.Title, s.PP)
}
```

**osu!standard only**

### All Players

This mode of tracking was yoinked almost one to one from [kaysting's osu-scores-cache](https://github.com/kaysting/osu-score-cache).
This one will give you scores for _all players_ for any mode you choose, so do expect quite a constant stream of scores.
I'm not totally sure, but I think this will only emit finished scores on ranked maps.

```go
var feed = tracker.NewRealtimePoller(app, tracker.RealtimeConfig{Priority: 1})
var scores = feed.Subscribe()

go feed.Run(ctx)
for s := range scores {
	fmt.Println(s.UserID, s.BeatmapID, s.PP)
}
```

Persist `feed.Cursor()`, pass it to `feed.Resume()` before `Run` to skip nothing across a restart.
The cursor ensures you're not missing scores, so you don't even have to poll that often to keep up with it.
If the poller starts falling behind, it will also speed up the interval to catch up, in case that happens.

## Rate limiting 

Same rules as [gosu-api](https://github.com/Minoxs/gosu-api) applies here. **Please be careful**.

## AI disclosure

Most of this code is pure human made slop, but the more recent versions have had a hefty usage of AI to quickly add new endpoints and rate limiting (and now the repo split).
To be honest, I hate it probably as much as you do, but a few years working as a software engineer and your hopes and dreams will absolutely be crushed as well.
Still, I had a lot of fun designing this thing and have reviewed things thoroughly, once I get things at a stable point I will make sure everything is cleaned up.
