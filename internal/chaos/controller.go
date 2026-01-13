package chaos

import (
	"sync"
	"time"
)

var (
	mu     sync.RWMutex
	config Config
)

func Set(cfg Config) {
	mu.Lock()
	defer mu.Unlock()
	config = cfg
}

func Get() Config {
	mu.RLock()
	defer mu.RUnlock()
	return config
}

func Clear() {
	mu.Lock()
	defer mu.Unlock()
	config = Config{}
}

func AutoRecover() {
	go func() {
		for {
			time.Sleep(1 * time.Second)
			mu.Lock()
			if config.Enabled && !config.ExpiresAt.IsZero() &&
				time.Now().After(config.ExpiresAt) {
				config = Config{}
			}
			mu.Unlock()
		}
	}()
}
