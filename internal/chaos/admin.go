package chaos

import (
	"encoding/json"
	"net/http"
	"time"
)

type Request struct {
	Route     string `json:"route"`
	DelayMs  int    `json:"delay_ms"`
	ErrorPct int    `json:"error_rate"`
	DropPct  int    `json:"drop_rate"`
	Duration int    `json:"duration_sec"`
}

func EnableHandler(w http.ResponseWriter, r *http.Request) {
	var req Request
	json.NewDecoder(r.Body).Decode(&req)

	cfg := Config{
		Enabled:   true,
		Route:     req.Route,
		Delay:     time.Duration(req.DelayMs) * time.Millisecond,
		ErrorRate: req.ErrorPct,
		DropRate:  req.DropPct,
	}

	if req.Duration > 0 {
		cfg.ExpiresAt = time.Now().Add(time.Duration(req.Duration) * time.Second)
	}

	Set(cfg)

	w.Write([]byte("Chaos enabled"))
}

func DisableHandler(w http.ResponseWriter, r *http.Request) {
	Clear()
	w.Write([]byte("Chaos disabled"))
}
