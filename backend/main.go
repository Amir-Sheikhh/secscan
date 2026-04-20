package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/amir-sheikh/secscan/backend/internal/api"
)

func main() {
	server := api.NewServer(api.Config{
		AllowedOrigins:      csvEnv("ALLOWED_ORIGINS", []string{"http://localhost:3000"}),
		ScanTimeout:         durationEnv("SECSCAN_SCAN_TIMEOUT", 45*time.Second),
		EnableActiveProbes:  strings.EqualFold(os.Getenv("SECSCAN_ENABLE_ACTIVE_PROBES"), "true"),
		EnableExternalIntel: true,
	})

	port := envOr("PORT", "8080")
	log.Printf("secscan backend listening on :%s", port)
	if err := http.ListenAndServe(":"+port, server.Router()); err != nil {
		log.Fatal(err)
	}
}

func envOr(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func csvEnv(key string, fallback []string) []string {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	items := strings.Split(raw, ",")
	values := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			values = append(values, item)
		}
	}
	if len(values) == 0 {
		return fallback
	}
	return values
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}

	value, err := time.ParseDuration(raw)
	if err != nil {
		return fallback
	}
	return value
}
