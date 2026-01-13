package main

import (
	"log"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/chaos"
	"github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/middleware"
	"github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/observability"
	"github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/proxy"
	"github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/ratelimit"
	"github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/tenant"
)

func main() {
	// ---- Tracing ----
	shutdown := observability.InitTracer("api-gateway")
	defer shutdown()

	// ---- Chaos auto-recovery watcher ----
	chaos.AutoRecover()

	// ---- Redis ----
	rdb := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	// ---- Rate Limiter ----
	rl := ratelimit.NewRateLimiter(rdb, 5, time.Minute)

	// ---- Backend proxies ----
	userHandler, _ := proxy.ProxyHandler("http://localhost:9001")
	orderHandler, _ := proxy.ProxyHandler("http://localhost:9002")

	// ---- Secure pipelines ----
	securedUserHandler :=
		chaos.Middleware(
			tenant.Middleware(
				rl.Middleware(userHandler),
			),
		)

	securedOrderHandler :=
		chaos.Middleware(
			tenant.Middleware(
				rl.Middleware(orderHandler),
			),
		)

	// ---- Router ----
	router := proxy.NewRouter()
	router.AddRoute("/users", securedUserHandler)
	router.AddRoute("/orders", securedOrderHandler)

	// ---- Global middleware stack ----
	finalHandler :=
		middleware.Logging(
			middleware.Metrics(
				middleware.Tracing(router),
			),
		)

	http.Handle("/", finalHandler)

	// ---- Admin Chaos Control ----
	http.HandleFunc("/admin/chaos/enable", chaos.EnableHandler)
	http.HandleFunc("/admin/chaos/disable", chaos.DisableHandler)

	log.Println("API Gateway running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
