package config

import "os"

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	JwtSecret  string
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

var GlobalConfig = Config{
	DBHost:     getEnvOrDefault("DB_HOST", "localhost"),
	DBUser:     getEnvOrDefault("DB_USER", "root"),
	DBPassword: getEnvOrDefault("DB_PASSWORD", "zth20041017"),
	DBName:     getEnvOrDefault("DB_NAME", "todo-app"),
	JwtSecret:  getEnvOrDefault("JWT_SECRET", "BiEryqig6Hg7UlkFJ3ODpb8lXGhuOU1TegOdbPxxGcytsHOBDg1KBWWYVBdPvEHe"),
}
