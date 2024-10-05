package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/hhertout/twirp_auth/internal/router"
	"github.com/hhertout/twirp_auth/migrations"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	if os.Getenv("GO_ENV") == "development" {
		logger, _ = zap.NewDevelopment()
	}
	defer logger.Sync()

	if os.Getenv("MIGRATION_ENABLE") == "true" {
		logger.Sugar().Info("🚀 Running migrations")
		m := migrations.NewMigration("/", logger)
		if err := m.MigrateAll(); err != nil {
			logger.Sugar().Fatalf("Error running migrations", err)
			return
		}
		logger.Sugar().Info("🚀 Migrations ran successfully")
	}

	router := router.GetRouter(logger)

	if os.Getenv("GO_ENV") == "development" {
		logger.Sugar().Info("🔨🔨 Caution : The server will be running under development mode 🔨🔨")
	}

	port, _ := strconv.Atoi(os.Getenv("PORT"))
	logger.Sugar().Info("🚀 Server running on port", port)
	http.ListenAndServe(fmt.Sprintf("%v:%v", "0.0.0.0", port), router)
}
