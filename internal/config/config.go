package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		Env      string   `env-required:"true" yaml:"env" env:"ENV"`
		HTTP     HTTP     `yaml:"http" env-prefix:"HTTP_"`
		Postgres Postgres `yaml:"postgres" env-prefix:"POSTGRES_"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"PORT"`
	}

	Postgres struct {
		Host     string `env-required:"true" yaml:"host" env:"HOST"`
		Port     string `env-required:"true" yaml:"port" env:"PORT"`
		User     string `env-required:"true" yaml:"user" env:"USER"`
		Password string `env-required:"true" yaml:"password" env:"PASSWORD"`
		DB       string `env-required:"true" yaml:"db" env:"DB"`
	}
)

func New(path string) (*Config, error) {
	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
