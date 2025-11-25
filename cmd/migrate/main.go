package main

import (
	"flag"

	"github.com/neokofg/callap-backend/internal/application/config"
	"github.com/neokofg/callap-backend/internal/infrastructure/database/postgresql/migrations"
	zapPkg "github.com/neokofg/callap-backend/pkg/zap"
	"go.uber.org/zap"
)

func main() {
	run := flag.String("run", "up", "")

	flag.Parse()

	logger := zapPkg.InitZap()
	cfg := config.InitConfig(logger)

	if err := migrations.Migrate(cfg, migrations.MigrateMode(*run)); err != nil {
		logger.Fatal("Unable to migrate database", zap.Error(err))
	} else {
		logger.Info("Successfully migrated database", zap.String("run", *run))
	}
}
