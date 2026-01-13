package proxy

import (
    "net/http"
    "strings"
    "github.com/CSroseX/Multi-tenant-Distributed-API-Gateway/internal/decisionlog"
)

type Route struct {
    Prefix  string
    Handler http.Handler
}

type Router struct {
    routes []Route
}

func NewRouter() *Router {
    return &Router{}
}

func (r *Router) AddRoute(prefix string, handler http.Handler) {
    r.routes = append(r.routes, Route{
        Prefix:  prefix,
        Handler: handler,
    })
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if strings.HasPrefix(req.URL.Path, route.Prefix) {
			decisionlog.LogDecision(req, decisionlog.DecisionRoute, "Routing to backend", map[string]any{
				"target": route.Prefix,
			})
			route.Handler.ServeHTTP(w, req)
			return
		}
	}

	decisionlog.LogDecision(req, decisionlog.DecisionBlock, "Route not found", nil)
	http.NotFound(w, req)
}
