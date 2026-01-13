package decisionlog

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

// DecisionType represents types of decisions made by the gateway
type DecisionType string

const (
	DecisionAllow DecisionType = "ALLOW"
	DecisionBlock DecisionType = "BLOCK"
	DecisionRoute DecisionType = "ROUTE"
	DecisionChaos DecisionType = "CHAOS"
)

// DecisionLog represents a structured log for intelligent decisions
type DecisionLog struct {
	Timestamp   time.Time           `json:"timestamp"`
	Decision    DecisionType        `json:"decision"`
	Reason      string              `json:"reason"`
	Tenant      string              `json:"tenant,omitempty"`
	Route       string              `json:"route,omitempty"`
	Target      string              `json:"target,omitempty"`
	Method      string              `json:"method,omitempty"`
	RequestID   string              `json:"request_id,omitempty"`
	ExtraFields map[string]any      `json:"extra,omitempty"`
}

// LogDecision prints a decision log in structured JSON format
func LogDecision(r *http.Request, decision DecisionType, reason string, extra map[string]any) {
	dl := DecisionLog{
		Timestamp:   time.Now(),
		Decision:    decision,
		Reason:      reason,
		Method:      r.Method,
		Route:       r.URL.Path,
		RequestID:   r.Header.Get("X-Request-ID"),
		Tenant:      r.Header.Get("X-Tenant-ID"),
		ExtraFields: extra,
	}

	data, err := json.Marshal(dl)
	if err != nil {
		log.Printf("[DECISION LOG ERROR] failed to marshal: %v", err)
		return
	}

	log.Println(string(data))
}
