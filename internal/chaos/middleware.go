package chaos

import (
	"math/rand"
	"net/http"
	"time"
	"github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/decisionlog"
)

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg := Get()

		if !cfg.Enabled {
			decisionlog.LogDecision(r, decisionlog.DecisionChaos, "Chaos disabled", nil)
			next.ServeHTTP(w, r)
			return
		}

		if cfg.Route != "" && cfg.Route != r.URL.Path {
			decisionlog.LogDecision(r, decisionlog.DecisionChaos, "Chaos not active for this route", nil)
			next.ServeHTTP(w, r)
			return
		}

		if cfg.Delay > 0 {
			decisionlog.LogDecision(r, decisionlog.DecisionChaos, "Injecting delay", map[string]any{
				"delay_ms": cfg.Delay.Milliseconds(),
			})
			time.Sleep(cfg.Delay)
		}

		if cfg.ErrorRate > 0 && rand.Intn(100) < cfg.ErrorRate {
			decisionlog.LogDecision(r, decisionlog.DecisionChaos, "Injecting error", map[string]any{
				"error_code": http.StatusServiceUnavailable,
			})
			http.Error(w, "Service Unavailable (chaos)", http.StatusServiceUnavailable)
			return
		}

		if cfg.DropRate > 0 && rand.Intn(100) < cfg.DropRate {
			decisionlog.LogDecision(r, decisionlog.DecisionChaos, "Dropping request", nil)
			return
		}

		next.ServeHTTP(w, r)
	})
}

