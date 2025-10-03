package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/WhiteMaks/go_academy_portal_service/util"
	"github.com/golang-migrate/migrate/v4"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("failed to load config:", err)
	}

	err = migrateDatabase("file://.database/migrations", config.Database.ConnectionString)
	if err != nil {
		log.Fatal("failed to migrate the database:", err)
	}

	err = migrateFunctions(".database/functions", config.Database.ConnectionString)
	if err != nil {
		log.Fatal("failed to migrate functions:", err)
	}
}

func migrateDatabase(migrationsDir string, connectionString string) error {
	migrations, err := migrate.New(migrationsDir, connectionString)
	if err != nil {
		return err
	}

	err = migrations.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func migrateFunctions(functionsDir string, connectionString string) error {
	dbConnection, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}

	files, _ := filepath.Glob(filepath.Join(functionsDir, "*.sql"))
	for _, f := range files {
		sqlBytes, _ := os.ReadFile(f)
		_, err := dbConnection.Exec(string(sqlBytes))
		if err != nil {
			return fmt.Errorf("failed in %s: %w", f, err)
		}
	}
	return nil
}
