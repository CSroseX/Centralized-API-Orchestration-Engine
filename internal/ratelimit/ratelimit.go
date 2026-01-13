package ratelimit

import (
    "context"
    "net/http"
    "strconv"
    "time"

    "github.com/redis/go-redis/v9"
    "github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/tenant"
    "github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/decisionlog"
)

type RateLimiter struct {
    redis *redis.Client
    limit int
    refill time.Duration
}
// constructor to make rate limiting configure. 
func NewRateLimiter(redis *redis.Client, limit int, refill time.Duration) *RateLimiter {
    return &RateLimiter{
        redis: redis,
        limit: limit,
        refill: refill,
    }
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t, ok := tenant.FromContext(r.Context())
		if !ok {
			decisionlog.LogDecision(r, decisionlog.DecisionBlock, "Tenant not found for rate limiting", nil)
			http.Error(w, "Tenant not found", http.StatusUnauthorized)
			return
		}

		key := "ratelimit:" + t.ID
		ctx := context.Background()

		tokensStr, err := rl.redis.Get(ctx, key).Result()
		if err == redis.Nil {
			rl.redis.Set(ctx, key, rl.limit-1, rl.refill)
			decisionlog.LogDecision(r, decisionlog.DecisionAllow, "Rate limit OK (first request)", map[string]any{
				"tenant": t.ID,
				"limit":  rl.limit,
				"used":   1,
			})
			next.ServeHTTP(w, r)
			return
		}

		tokens, _ := strconv.Atoi(tokensStr)
		if tokens <= 0 {
			decisionlog.LogDecision(r, decisionlog.DecisionBlock, "Rate limit exceeded", map[string]any{
				"tenant": t.ID,
				"limit":  rl.limit,
				"used":   tokens,
			})
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		rl.redis.Decr(ctx, key)
		decisionlog.LogDecision(r, decisionlog.DecisionAllow, "Rate limit OK", map[string]any{
			"tenant": t.ID,
			"limit":  rl.limit,
			"used":   tokens,
		})

		next.ServeHTTP(w, r)
	})
}

