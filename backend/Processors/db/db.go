package db

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./LHS.db")
	if err != nil {
		return err
	}
	_, err = DB.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		return err
	}
	createTables := `
	CREATE TABLE IF NOT EXISTS stacks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS stack_services (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		stack_id INTEGER,
		container_id TEXT,
		name TEXT,
		image TEXT,
		build_path TEXT,
		build_dockerfile TEXT,
		ports TEXT,
		env TEXT,
		volumes TEXT,
		FOREIGN KEY(stack_id) REFERENCES stacks(id) ON DELETE CASCADE
	);`

	_, err = DB.Exec(createTables)
	if err != nil {
		return err
	}

	return nil
}
