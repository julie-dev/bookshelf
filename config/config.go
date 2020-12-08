package config

type Config struct {
	RESTAPI_URL string `env:"RESTAPI_URL,required"`
	RESTAPI_KEY string `env:"RESTAPI_KEY,required"`
}
