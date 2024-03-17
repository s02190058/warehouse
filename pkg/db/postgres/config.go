package postgres

import (
	"net"
	"net/url"
)

const (
	scheme = "postgres"
)

type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DB       string
}

func (c Config) URL() string {
	u := &url.URL{
		Scheme: scheme,
		User:   url.UserPassword(c.User, c.Password),
		Host:   net.JoinHostPort(c.Host, c.Port),
		Path:   c.DB,
	}

	return u.String()
}
