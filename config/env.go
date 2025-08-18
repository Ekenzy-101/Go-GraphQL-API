package config

import "os"

type DBType string

var (
	DataBaseMongo    DBType = "mongo"
	DataBasePostgres DBType = "postgres"
)

func AccessTokenSecret() string {
	return os.Getenv("ACCESS_TOKEN_SECRET")
}

func CacheURL() string {
	return os.Getenv("CACHE_URL")
}

func DataBaseName() string {
	return os.Getenv("DATABASE_NAME")
}

func DataBaseType() DBType {
	return DBType(os.Getenv("DATABASE_TYPE"))
}

func DataBaseURL() string {
	return os.Getenv("DATABASE_URL")
}

func Port() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	return port
}

func IsProduction() bool {
	return os.Getenv("GIN_MODE") == "release"
}
