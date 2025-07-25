package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	Host string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port int    `yaml:"port" env:"PORT" env-default:"8080"`
}

type PostgresConfig struct {
	Server   ServerConfig `yaml:"server_db"`
	User     string       `yaml:"user" env:"DB_USER" env-default:"user"`
	Password string       `yaml:"password" env:"DB_PASSWORD" env-default:"password"`
	Name     string       `yaml:"name" env:"DB_NAME" env-default:"mydatabase"`
}

type MongoConfig struct {
	Host     string `yaml:"host" env:"MONGO_HOST" env-default:"localhost"`
	Port     int    `yaml:"port" env:"MONGO_PORT" env-default:"27017"`
	User     string `yaml:"user" env:"MONGO_USER" env-default:"mongo_user"`
	Password string `yaml:"password" env:"MONGO_PASSWORD" env-default:"mongo_password"`
	Name     string `yaml:"name" env:"MONGO_DB" env-default:"reportsdb"`
}

type DatabaseConfig struct {
	Postgres PostgresConfig `yaml:"postgres"`
	Mongo    MongoConfig    `yaml:"mongo"`
}

type AuthConfig struct {
	JWTSecret string `yaml:"jwt_secret" env:"JWT_SECRET" env-default:"mysecretkey"`
}

type Config struct {
	Server   ServerConfig   `yaml:"server_config"`
	Database DatabaseConfig `yaml:"database"`
	Auth     AuthConfig     `yaml:"auth"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}
