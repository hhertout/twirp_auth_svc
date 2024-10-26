package migrations

// Caution: each migration query must be separated by a double dash (--) in the migration file.

import (
	"database/sql"
	"errors"
	"io"
	"os"
	"slices"
	"sort"
	"strings"
	"time"

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
//
// @Caution: each migration query must be separated by a double dash (--) in the migration file.
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

	if err := m.createMigrationTable(); err != nil {
		return errors.New("failed to create migration table")
	}

	migrationFiles, err := m.GetMigrationFiles(m.basePath)
	if err != nil {
		return errors.New("failed to retrieve migration files")
	}

	migrationAlreadyExecuted, err := m.getExecutedMigrations()
	if err != nil {
		return errors.New("failed to retrieve executed migrations")
	}

	if len(migrationFiles) == 0 {
		m.logger.Sugar().Info("No migration file found! To add one, run 'make migration-generate'.")
	} else {
		for _, f := range migrationFiles {
			if slices.Contains(migrationAlreadyExecuted, f) {
				m.logger.Sugar().Infof("⏭️ Migration file already executed: %s\n", f)
				continue
			} else {
				err := m.migrateFromFile(f)
				if err != nil {
					return err
				}
				// Save the executed migration
				_, err = m.dbPool.Exec("INSERT INTO go_migrations.migration (filename, migrated_at) VALUES ($1, $2)", f, time.Now())
				if err != nil {
					return err
				}
				m.logger.Sugar().Infof("✅ Migrated file: %v", f)
			}
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

// Create migration table
func (m *Migration) createMigrationTable() error {
	//create migration schema if not exists
	_, err := m.dbPool.Exec(`
        CREATE SCHEMA IF NOT EXISTS go_migrations;
    `)
	if err != nil {
		return err
	}

	_, err = m.dbPool.Exec(`
        CREATE TABLE IF NOT EXISTS go_migrations.migration (
            id SERIAL PRIMARY KEY,
            filename VARCHAR(255) NOT NULL,
            migrated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        );
    `)
	if err != nil {
		return err
	}

	return nil
}

// Retieve migration already executed
func (m *Migration) getExecutedMigrations() ([]string, error) {
	var res []string

	rows, err := m.dbPool.Query("SELECT filename FROM go_migrations.migration")
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			return nil, err
		}
		res = append(res, filename)
	}

	return res, nil
}
