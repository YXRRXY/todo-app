package config

type Config struct {
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	JwtSecret  string
}

var GlobalConfig = Config{
	DBHost:     "localhost",
	DBUser:     "root",
	DBPassword: "123456",
	DBName:     "todo-app",
	JwtSecret:  "BiEryqig6Hg7UlkFJ3ODpb8lXGhuOU1TegOdbPxxGcytsHOBDg1KBWWYVBdPvEHe",
}
