package config

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

type Config struct {
	ListenAddr string
	TTL        time.Duration
}

var (
	cfg  Config
	once sync.Once
)

const prefix = "KVEVRI"

func Get() Config {
	once.Do(func() {
		cfg = Config{ListenAddr: getenv("LISTEN_ADDR", ":8080")}

		ttlStr := getenv("TTL", "1h")
		ttl, err := time.ParseDuration(ttlStr)
		if err != nil {
			log.Printf("failed to parse TTL: %s", err)
		}

		cfg.TTL = ttl
	})

	return cfg
}

func getenv(key, fallback string) string {
	fullKey := fmt.Sprintf("%s_%s", prefix, key)

	if value, ok := os.LookupEnv(fullKey); ok {
		return value
	}

	return fallback
}
