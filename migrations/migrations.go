package migrations

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	// migrate dependency.
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed warehouse
var migrations embed.FS

type options struct {
	steps int
}

func ApplyMigrations(dbURL string, opts ...Option) error {
	o := new(options)
	for _, opt := range opts {
		opt(o)
	}

	source, err := iofs.New(migrations, "warehouse")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, dbURL)
	if err != nil {
		return err
	}

	if o.steps != 0 {
		err = m.Steps(o.steps)
	} else {
		err = m.Up()
	}
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return errors.Join(m.Close())
}
