package sqlite

import (
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"

	_ "modernc.org/sqlite"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

// Database represents a SQLite database
type Database struct {
	DB *sql.DB
}

// NewDatabase creates a SQLite connection and then runs any applicable migrations or seeders.
func NewDatabase(conStr string) (*Database, error) {
	con, err := sql.Open("sqlite", conStr)
	if err != nil {
		return nil, err
	}

	db := &Database{
		DB: con,
	}

	// run the migrations
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("migrate: %w", err)
	}

	return db, nil
}

func (d *Database) migrate() error {
	// you have to enable WAL outside of a transaction, annoyingly
	if _, err := d.DB.Exec("PRAGMA journal_mode = wal;"); err != nil {
		return fmt.Errorf("unable to enable wal: %w", err)
	}

	// you have to enable foreign key checks outside of a transaction.
	if _, err := d.DB.Exec("PRAGMA foreign_keys = ON;"); err != nil {
		return fmt.Errorf("unable to enable foreign keys: %w", err)
	}

	// create the migrations table if it doesn't yet exist.
	if _, err := d.DB.Exec("CREATE TABLE IF NOT EXISTS migrations (name TEXT PRIMARY KEY);"); err != nil {
		return fmt.Errorf("create migration table: %w", err)
	}

	// retrieve a list of migration files to execute.
	fileNames, err := fs.Glob(migrationFS, "migrations/*.sql")
	if err != nil {
		return fmt.Errorf("globbing migration files: %w", err)
	}
	sort.Strings(fileNames)

	// begin a transaction, rolling it back by default.
	tx, err := d.DB.Begin()
	if err != nil {
		return fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Rollback()

	// then execute them all.
	for _, fileName := range fileNames {
		if err = d.migrateFile(fileName, tx); err != nil {
			return err
		}
	}

	// commit the transaction if we got this far.
	tx.Commit()

	return nil
}

func (d *Database) migrateFile(fileName string, tx *sql.Tx) error {
	// check if the migration has been ran before and, if it has, return early.
	var c int
	if err := tx.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", fileName).Scan(&c); err != nil {
		return err
	} else if c != 0 {
		return nil
	}

	// read the file and execute it against the database.
	if buf, err := fs.ReadFile(migrationFS, fileName); err != nil {
		return err
	} else if _, err := tx.Exec(string(buf)); err != nil {
		return err
	}

	// insert the migration record into the table.
	if _, err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", fileName); err != nil {
		return err
	}

	return nil
}
