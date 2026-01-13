package chaos

import "time"

type Config struct {
	Enabled      bool
	Route        string        // empty = all routes
	Delay        time.Duration // artificial delay
	ErrorRate    int           // % chance to return 503
	DropRate     int           // % chance to drop request
	ExpiresAt    time.Time     // auto recovery time
}
