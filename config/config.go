package config

type Config struct {
	OpenAPIUrl string `env:"OPEN_API_URL,required"`
	OpenAPIKey string `env:"OPEN_API_KEY,required"`
	Database   DatabaseConfig
}

type DatabaseConfig struct {
	User     string `env:"DB_CONFIG_USER,required"`
	Password string `env:"DB_CONFIG_PASSWORD,required"`
	Host     string `env:"DB_CONFIG_HOST,required"`
	Port     int    `env:"DB_CONFIG_PORT,required"`
	Database string `env:"DB_CONFIG_DB,required"`
}
