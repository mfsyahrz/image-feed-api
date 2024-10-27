package config

import (
	"encoding/json"
	"fmt"

	"github.com/joeshaw/envdecode"
	"github.com/joho/godotenv"
)

type Config struct {
	Service     Service
	Postgres    Postgres
	FileStorage FileStorage
}

type Service struct {
	Name string `env:"SERVICE_NAME"`
	Port Port
}

type FileStorage struct {
	BaseURL  string `env:"FILE_STORAGE_BASE_URL"`
	BasePath string `env:"FILE_STORAGE_BASE_PATH"`
}

type Port struct {
	REST string `env:"SERVICE_PORT_REST"`
}

type Postgres struct {
	User            string `env:"POSTGRES_USER,required"`
	Password        string `env:"POSTGRES_PASSWORD,required"`
	Name            string `env:"POSTGRES_NAME,required"`
	Port            string `env:"POSTGRES_PORT,default=5432"`
	Host            string `env:"POSTGRES_HOST,default=localhost"`
	MaxOpenConns    int    `env:"POSTGRES_MAX_OPEN_CONNS,default=5"`
	MaxConnLifetime int    `env:"POSTGRES_MAX_CONN,default=10"`
	MaxIdleLifetime int    `env:"POSTGRES_MAX_IDLE,default=5"`
}

func (c *Config) String() string {
	js, _ := json.MarshalIndent(c, "", " ")
	return string(js)
}

func New(envFile string) (*Config, error) {

	_ = godotenv.Load(envFile)

	config := Config{}
	if err := envdecode.Decode(&config); err != nil {
		return nil, err
	}

	fmt.Println(config.String())
	return &config, nil
}
