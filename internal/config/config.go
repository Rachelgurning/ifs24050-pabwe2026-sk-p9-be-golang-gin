package config

import (
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv         string
	AppPort        string
	AppKey         string
	AdminUsername  string
	AdminPassword  string
	DBConnection   string
	DBHost         string
	DBPort         string
	DBDatabase     string
	DBUsername     string
	DBPassword     string
	DBSSLMode      string
	DatabaseURL    string
	SearchEndpoint string
	OpenAIEndpoint string
	OpenAIKey      string
	OpenAIModel    string
	SeedSample     bool
}

func Get(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func Load() *Config {
	_ = godotenv.Load()
	port := "8081"
	for _, k := range []string{"APP_PORT", "PORT"} {
		if v := Get(k, ""); v != "" {
			if _, err := strconv.Atoi(v); err == nil {
				port = v
			}
		}
	}
	return &Config{
		AppEnv:         Get("APP_ENV", "development"),
		AppPort:        port,
		AppKey:         Get("APP_KEY", "delcom_blog_secret_key_change_me"),
		AdminUsername:  Get("ADMIN_USERNAME", "admin"),
		AdminPassword:  Get("ADMIN_PASSWORD", "admin"),
		DBConnection:   Get("DB_CONNECTION", "pgsql"),
		DBHost:         Get("DB_HOST", "127.0.0.1"),
		DBPort:         Get("DB_PORT", "5432"),
		DBDatabase:     Get("DB_DATABASE", "pabwe_p9_blog"),
		DBUsername:     Get("DB_USERNAME", "postgres"),
		DBPassword:     Get("DB_PASSWORD", "postgres"),
		DBSSLMode:      Get("DB_SSLMODE", "disable"),
		DatabaseURL:    Get("DATABASE_URL", ""),
		SearchEndpoint: Get("SEARCH_ENGINE_ENDPOINT", "https://delcom.org/api/v1/search"),
		OpenAIEndpoint: Get("OPENAI_ENDPOINT", "https://delcom.org/api/v1"),
		OpenAIKey:      Get("OPENAI_KEY", ""),
		OpenAIModel:    Get("OPENAI_MODEL", ""),
		SeedSample:     Get("SEED_SAMPLE_POSTS", "true") != "false",
	}
}
