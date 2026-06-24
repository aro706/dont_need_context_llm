package database

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE Nodes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_name TEXT,
		file_path TEXT,
		function_name TEXT,
		type TEXT
	);
	CREATE TABLE Calls (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_name TEXT,
		caller TEXT,
		callee TEXT
	);`

	_, err = db.Exec(createTableQuery)
	return db, err
}

func ClearDB(db *sql.DB) {
	db.Exec("DELETE FROM Nodes")
	db.Exec("DELETE FROM Calls")
	db.Exec("DELETE FROM sqlite_sequence WHERE name='Nodes' OR name='Calls'")
}
