//go:build integration

package warehouse_test

import (
	"context"
	_ "embed"
	"log/slog"
	"net"
	"net/url"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/stdlib"
	service "github.com/s02190058/warehouse/internal/service/warehouse"
	storage "github.com/s02190058/warehouse/internal/storage/warehouse"
	"github.com/s02190058/warehouse/migrations"
	"github.com/s02190058/warehouse/pkg/db/postgres"
	"github.com/s02190058/warehouse/pkg/slogger"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	postgresUser     = "postgres"
	postgresPassword = "postgres"
	postgresDB       = "postgres"

	postgresContainerPort = "5432"

	migrationStepWarehouse = 1
)

const (
	setupSuiteTimeout    = 120 * time.Second
	tearDownSuiteTimeout = 10 * time.Second

	setupTestTimeout    = 5 * time.Second
	tearDownTestTimeout = 5 * time.Second
)

func createPostgresTestContainer(t *testing.T, ctx context.Context) (testcontainers.Container, error) {
	t.Helper()

	env := map[string]string{
		"POSTGRES_USER":     postgresUser,
		"POSTGRES_PASSWORD": postgresPassword,
		"POSTGRES_DB":       postgresDB,
	}

	dbURL := func(host string, port nat.Port) string {
		u := &url.URL{
			Scheme: "postgres",
			User:   url.UserPassword(postgresUser, postgresPassword),
			Host:   net.JoinHostPort(host, port.Port()),
			Path:   postgresDB,
		}

		return u.String()
	}

	req := testcontainers.ContainerRequest{
		Image:        "postgres:16",
		Env:          env,
		ExposedPorts: []string{postgresContainerPort + "/tcp"},
		WaitingFor: wait.ForSQL(postgresContainerPort, "pgx", dbURL).
			WithStartupTimeout(10 * time.Second).WithQuery("SELECT 10"),
	}

	postgresC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	return postgresC, nil
}

type Suite struct {
	suite.Suite

	logger *slog.Logger

	postgresC testcontainers.Container

	storage *storage.Storage
	db      *postgres.Database
}

func (s *Suite) SetupSuite() {
	t := s.T()

	logger, err := slogger.New("local")
	require.NoError(t, err)

	s.logger = logger

	ctx, cancel := context.WithTimeout(context.Background(), setupSuiteTimeout)
	defer cancel()

	postgresC, err := createPostgresTestContainer(t, ctx)
	require.NoError(t, err)

	port, err := postgresC.MappedPort(ctx, postgresContainerPort)
	require.NoError(t, err)

	s.postgresC = postgresC

	cfg := postgres.Config{
		Host:     "localhost",
		Port:     port.Port(),
		User:     postgresUser,
		Password: postgresPassword,
		DB:       postgresDB,
	}

	db, err := postgres.New(logger, cfg)
	require.NoError(t, err)

	s.db = db

	s.storage = storage.New(db)

	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(postgresUser, postgresPassword),
		Host:   net.JoinHostPort("localhost", port.Port()),
		Path:   postgresDB,
	}

	params := url.Values{}
	params.Set("sslmode", "disable")

	u.RawQuery = params.Encode()

	dbURL := u.String()

	err = migrations.ApplyMigrations(dbURL, migrations.WithSteps(migrationStepWarehouse))
	require.NoError(t, err)
}

func (s *Suite) TeardownSuite() {
	t := s.T()

	ctx, cancel := context.WithTimeout(context.Background(), tearDownSuiteTimeout)
	defer cancel()

	err := s.postgresC.Terminate(ctx)
	require.NoError(t, err)

	s.db.Close()
}

//go:embed testdata/setup_test.sql
var setupTestSQL string

func (s *Suite) SetupTest() {
	t := s.T()

	ctx, cancel := context.WithTimeout(context.Background(), setupTestTimeout)
	defer cancel()

	_, err := s.db.Query(ctx).Exec(ctx, setupTestSQL)
	require.NoError(t, err)
}

//go:embed testdata/tear_down_test.sql
var tearDownSQL string

func (s *Suite) TearDownTest() {
	t := s.T()

	ctx, cancel := context.WithTimeout(context.Background(), tearDownTestTimeout)
	defer cancel()

	_, err := s.db.Query(ctx).Exec(ctx, tearDownSQL)
	require.NoError(t, err)
}

func (s *Suite) TestGet() {
	t := s.T()

	testCases := []struct {
		name        string
		id          int
		isAvailable bool
		err         error
	}{
		{
			name:        "available",
			id:          1,
			isAvailable: true,
		},
		{
			name:        "unavailable",
			id:          2,
			isAvailable: false,
		},
		{
			name: "unknown warehouse",
			id:   4,
			err:  service.ErrWarehouseNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx := context.Background()

			wh, err := s.storage.Get(ctx, tc.id)
			require.Equal(t, tc.err, err)
			require.Equal(t, tc.isAvailable, wh.IsAvailable)
		})
	}
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(Suite))
}
