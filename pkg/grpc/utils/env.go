package utils

import (
	"flag"
	"os"
	"strings"
)

// GetEnvOrDefault get env or default value
func GetEnvOrDefault(key string, def string) string {
	if v := flag.Lookup(strings.ToLower(key)); v != nil {
		return v.Value.String()
	}
	if v := os.Getenv(key); v != "" {
		return v
	}

	return def
}
