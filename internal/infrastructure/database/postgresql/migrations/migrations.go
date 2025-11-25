package migrations

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/infrastructure/database/postgresql"
)

type MigrateMode string

const (
	MigrateUp   MigrateMode = "up"
	MigrateDrop MigrateMode = "drop"
)

const (
	localFolder = "./internal/infrastructure/database/postgresql/migrations/sql"
)

func Migrate(cfg *config.Config, mode MigrateMode) error {
	var conn *pgx.Conn
	var err error
	var folder string

	conn, err = postgresql.Conn(postgresql.Config{
		Username: cfg.PostgreSQL.Username,
		Password: cfg.PostgreSQL.Password,
		Host:     cfg.PostgreSQL.Host,
		Port:     cfg.PostgreSQL.Port,
		Database: cfg.PostgreSQL.Database,
		Pool: postgresql.ConfigPool{
			MaxConns:          cfg.PostgreSQL.Pool.MaxConns,
			MinConns:          cfg.PostgreSQL.Pool.MinConns,
			MaxConnLifeTime:   cfg.PostgreSQL.Pool.MaxConnLifeTime,
			MaxConnIdleTime:   cfg.PostgreSQL.Pool.MaxConnIdleTime,
			HealthCheckPeriod: cfg.PostgreSQL.Pool.HealthCheckPeriod,
		},
	})
	if err != nil {
		return err
	}

	folder = localFolder

	var re *regexp.Regexp

	switch mode {
	case MigrateUp:
		re = regexp.MustCompile(`^.*\.up\.sql$`)
	case MigrateDrop:
		re = regexp.MustCompile(`^.*\.down\.sql$`)
	default:
		log.Fatal(errors.New("[flag error] - choose: --run=create or --run=drop"))
	}

	migrationFiles, err := os.ReadDir(folder)
	if err != nil {
		return err
	}

	files := make([]string, 0, int(len(migrationFiles)/2))

	for _, f := range migrationFiles {
		if re.MatchString(f.Name()) {
			files = append(files, f.Name())
		}
	}

	var sortErrors []error
	sort.Slice(files, func(i, j int) bool {
		numi, err := strconv.Atoi(files[i][:6])
		if err != nil {
			sortErrors = append(sortErrors, err)
			return false
		}
		numj, err := strconv.Atoi(files[j][:6])
		if err != nil {
			sortErrors = append(sortErrors, err)
			return false
		}

		if mode == "drop" {
			return numi > numj
		}
		return numi < numj
	})

	if len(sortErrors) > 0 {
		return sortErrors[0]
	}

	for _, f := range files {
		data, err := os.ReadFile(fmt.Sprintf("%s/%s", folder, f))
		if err != nil {
			return err
		}
		queries := strings.Split(string(data), ";")

		for _, query := range queries {
			query = strings.TrimSpace(query)
			if query == "" {
				continue
			}

			_, err = conn.Exec(context.Background(), query)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
