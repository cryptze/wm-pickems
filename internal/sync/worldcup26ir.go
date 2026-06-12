package sync

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

const wc26URL = "https://worldcup26.ir/get/games"

type wc26Game struct {
	HomeTeamName string `json:"home_team_name_en"`
	AwayTeamName string `json:"away_team_name_en"`
	HomeScore    string `json:"home_score"`
	AwayScore    string `json:"away_score"`
	Finished     string `json:"finished"`
	TimeElapsed  string `json:"time_elapsed"`
}

// wc26Sync pulls worldcup26.ir's live JSON and applies results.
// Matches by normalised team-name pair — same approach as API-Football sync.
// Idempotent: a record is only saved when something actually changed.
func wc26Sync(ctx context.Context, app core.App) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, wc26URL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "wm-tips/1.0")
	resp, err := (&http.Client{Timeout: 20 * time.Second}).Do(req)
	if err != nil {
		return fmt.Errorf("worldcup26ir fetch: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("worldcup26ir: status %d", resp.StatusCode)
	}
	var doc struct {
		Games []wc26Game `json:"games"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&doc); err != nil {
		return err
	}

	matches, err := app.FindRecordsByFilter("matches", "id != ''", "", 0, 0)
	if err != nil {
		return err
	}
	teams, _ := app.FindRecordsByFilter("teams", "id != ''", "", 0, 0)
	teamName := map[string]string{}
	for _, t := range teams {
		teamName[t.Id] = canonName(t.GetString("name"))
	}
	byPair := map[string]*core.Record{}
	for _, m := range matches {
		h, a := teamName[m.GetString("homeTeam")], teamName[m.GetString("awayTeam")]
		if h != "" && a != "" {
			byPair[h+"|"+a] = m
		}
	}

	updated := 0
	for _, g := range doc.Games {
		status := "scheduled"
		switch {
		case g.Finished == "TRUE":
			status = "finished"
		case g.TimeElapsed == "1H" || g.TimeElapsed == "HT" || g.TimeElapsed == "2H" ||
			g.TimeElapsed == "ET" || g.TimeElapsed == "P":
			status = "live"
		}

		if status == "scheduled" {
			continue
		}

		ftH, errH := strconv.Atoi(g.HomeScore)
		ftA, errA := strconv.Atoi(g.AwayScore)
		if errH != nil || errA != nil {
			continue
		}

		key := canonName(g.HomeTeamName) + "|" + canonName(g.AwayTeamName)
		rec, ok := byPair[key]
		if !ok {
			continue
		}

		// Skip if nothing changed.
		if rec.GetString("status") == status &&
			rec.GetInt("ftHome") == ftH && rec.GetInt("ftAway") == ftA {
			continue
		}

		applyResult(rec, status, &ftH, &ftA, nil, nil, nil, nil)
		if app.Save(rec) == nil {
			updated++
		}
	}
	if err := ResolveBracket(app); err != nil {
		return err
	}
	return nil
}
