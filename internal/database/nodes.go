package database

import (
	"database/sql"
	"fmt"
	"micro-mcp/internal/models"
	"strings"
)

func InsertNode(db *sql.DB, n models.Node) {
	db.Exec("INSERT INTO Nodes (project_name, file_path, function_name, type) VALUES (?, ?, ?, ?)",
		n.ProjectName, n.FilePath, n.FunctionName, n.Type)
}

func SearchSymbols(db *sql.DB, query string, projectFilter string) string {
	var rows *sql.Rows
	var err error

	if projectFilter != "" {
		rows, err = db.Query("SELECT file_path, function_name, type, project_name FROM Nodes WHERE (function_name LIKE ? OR file_path LIKE ?) AND project_name = ?", "%"+query+"%", "%"+query+"%", projectFilter)
	} else {
		rows, err = db.Query("SELECT file_path, function_name, type, project_name FROM Nodes WHERE function_name LIKE ? OR file_path LIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if err != nil {
		return "Database error occurred."
	}
	defer rows.Close()

	var output []string
	for rows.Next() {
		var path, name, nodeType, proj string
		rows.Scan(&path, &name, &nodeType, &proj)
		output = append(output, fmt.Sprintf("[%s] (%s) '%s' found in file: %s", proj, nodeType, name, path))
	}
	if len(output) == 0 {
		return "No matching elements found."
	}
	return strings.Join(output, "\n")
}
