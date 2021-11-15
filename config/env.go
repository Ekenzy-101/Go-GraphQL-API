package config

import "os"

func AccessTokenSecret() string {
	return os.Getenv("ACCESS_TOKEN_SECRET")
}

func CacheURL() string {
	return os.Getenv("CACHE_URL")
}

func DataBaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func DataBaseName() string {
	return os.Getenv("DATABASE_NAME")
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	return port
}

func IsProduction() bool {
	return os.Getenv("APP_ENV") == "prod"
}
