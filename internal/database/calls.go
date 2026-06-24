package database

import (
	"database/sql"
	"fmt"
	"micro-mcp/internal/models"
	"strings"
)

func InsertCall(db *sql.DB, c models.Call) {
	db.Exec("INSERT INTO Calls (project_name, caller, callee) VALUES (?, ?, ?)",
		c.ProjectName, c.Caller, c.Callee)
}

func GetDependencies(db *sql.DB, project string) string {
	rows, err := db.Query("SELECT caller, callee FROM Calls WHERE project_name = ?", project)
	if err != nil {
		return "Database error occurred."
	}
	defer rows.Close()

	var pairs []string
	for rows.Next() {
		var caller, callee string
		rows.Scan(&caller, &callee)
		pairs = append(pairs, fmt.Sprintf("Caller '%s' -> Calls '%s'", caller, callee))
	}

	output := strings.Join(pairs, "\n")
	if output == "" {
		return "No dependency patterns found."
	}
	return output
}
