package migrations

// Caution: each migration query must be separated by a double dash (--) in the migration file.

import (
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/hhertout/twirp_auth/pkg/database"
	"go.uber.org/zap"
)

// Migration represents a database migration with a connection pool and a base path for migration files.
// Base path is the path to the directory containing the migration files from /migrations.
// DbPool is the connection pool to the database.
//
// @Caution: each migration query must be separated by a double dash (--) in the migration file.
type Migration struct {
	dbPool   *sql.DB
	basePath string
	logger   *zap.Logger
}

// NewMigration creates a new instance of Migration with the given base path.
// Returns a pointer to the newly created Migration.
// Base path is the path to the directory containing the migration files from /migrations.
func NewMigration(basePath string, logger *zap.Logger) *Migration {
	return &Migration{
		nil,
		basePath,
		logger,
	}
}

// Migrate executes a migration from a specified file.
// Connects to the database, runs the migration, and then closes the database connection.
// Returns an error if any occurs during the process.
func (m *Migration) Migrate(filename string) error {
	db, err := database.Connect()
	if err != nil {
		return errors.New("failed to connect to db")
	}
	m.dbPool = db.DbPool

	if err := m.migrateFromFile(filename); err != nil {
		return err
	}

	if err = m.dbPool.Close(); err != nil {
		return errors.New("failed to close db connection after executing migrations")
	}

	return nil
}

// MigrateAll executes all migration files found in the base path.
// Connects to the database, runs all migrations, and then closes the database connection.
// Returns an error if any occurs during the process.
func (m *Migration) MigrateAll() error {
	db, err := database.Connect()
	if err != nil {
		return errors.New("failed to connect to db")
	}
	m.dbPool = db.DbPool

	migrationFiles, err := m.GetMigrationFiles(m.basePath)
	if err != nil {
		return errors.New("failed to retrieve migration files")
	}
	if len(migrationFiles) == 0 {
		log.Println("No migration file found! To add one, run 'make migration-generate'.")
	} else {
		for _, f := range migrationFiles {
			err := m.migrateFromFile(f)
			if err != nil {
				return err
			}
			fmt.Println("✅ Migrated file: ", f)
		}
	}

	if err = m.dbPool.Close(); err != nil {
		return errors.New("failed to close db connection after migration")
	}

	return nil
}

// migrateFromFile executes the SQL queries from a specified migration file.
// Returns an error if any occurs during the process.
func (m *Migration) migrateFromFile(filename string) error {
	workingDir, _ := os.Getwd()
	fileOpen, err := os.Open(workingDir + m.basePath + "/migrations/" + filename)
	if err != nil {
		return err
	}
	defer fileOpen.Close()

	content, err := io.ReadAll(fileOpen)
	if err != nil {
		return err
	}

	queries := string(content)
	queriesSplit := strings.Split(queries, "--")

	for _, query := range queriesSplit {
		if strings.TrimSpace(query) == "" {
			continue
		}

		_, err = m.dbPool.Exec(query + ";")
		if err != nil {
			if err.Error() == "trigger set_viewed_param already exists" {
				continue
			}
			return err
		}
	}
	return nil
}

// GetMigrationFiles retrieves all SQL migration files from the base path.
// Returns a slice of file names and an error if any occurs during the process.
func (m *Migration) GetMigrationFiles(basePath string) ([]string, error) {
	var res []string
	baseDir := "migrations"
	workingDir, _ := os.Getwd()

	dir, err := os.ReadDir(workingDir + basePath + baseDir)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dir {
		if !dirEntry.IsDir() && strings.HasSuffix(dirEntry.Name(), ".sql") {
			res = append(res, dirEntry.Name())
		}
	}

	sort.Strings(res)

	return res, nil
}
