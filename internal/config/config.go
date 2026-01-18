package config

import "os"

type Config struct {
	Addr   string
	DBPath string
}

func Default() Config {
	return Config{
		Addr:   envOr("ASD_ADDR", "127.0.0.1:8080"),
		DBPath: envOr("ASD_DB_PATH", "asd.db"),
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
