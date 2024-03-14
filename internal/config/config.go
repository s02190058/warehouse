package config

import "github.com/ilyakaznacheev/cleanenv"

type (
	Config struct {
		Env  string `env-required:"true" yaml:"env" env:"ENV"`
		HTTP HTTP   `yaml:"http" env-prefix:"HTTP_"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"PORT"`
	}
)

func New(path string) (*Config, error) {
	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
