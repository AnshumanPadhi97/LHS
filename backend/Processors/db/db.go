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

	createTables := `
	CREATE TABLE IF NOT EXISTS stacks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		code TEXT UNIQUE,
		name TEXT,
		label TEXT,
		description TEXT,
		version INTEGER,
		tags TEXT,
		volumes TEXT,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS stack_services (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		container_id TEXT,
		stack_code TEXT,
		name TEXT,
		code TEXT,
		image TEXT,
		ports TEXT,
		env TEXT,
		volumes TEXT,
		depends_on TEXT,
		tunnel INTEGER,
		build_context TEXT,
		build_dockerfile TEXT,
		FOREIGN KEY(stack_code) REFERENCES stacks(code)
	);`

	_, err = DB.Exec(createTables)
	if err != nil {
		return err
	}
	return nil
}
