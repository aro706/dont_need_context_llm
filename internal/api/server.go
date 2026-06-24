package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

func StartHTTPServer(db *sql.DB) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "frontend/index.html")
	})

	http.HandleFunc("/api/projects", func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT DISTINCT project_name FROM Nodes")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var projects []string
		for rows.Next() {
			var p string
			rows.Scan(&p)
			projects = append(projects, p)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(projects)
	})

	http.HandleFunc("/api/graph", func(w http.ResponseWriter, r *http.Request) {
		project := r.URL.Query().Get("project")

		nodesRows, err := db.Query("SELECT function_name, type FROM Nodes WHERE project_name = ?", project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer nodesRows.Close()

		var nodes []map[string]interface{}
		for nodesRows.Next() {
			var name, nType string
			nodesRows.Scan(&name, &nType)
			nodes = append(nodes, map[string]interface{}{
				"data": map[string]string{
					"id":    name,
					"label": name,
					"type":  nType,
				},
			})
		}

		edgesRows, err := db.Query("SELECT caller, callee FROM Calls WHERE project_name = ?", project)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer edgesRows.Close()

		var edges []map[string]interface{}
		for edgesRows.Next() {
			var caller, callee string
			edgesRows.Scan(&caller, &callee)
			edges = append(edges, map[string]interface{}{
				"data": map[string]string{
					"source": caller,
					"target": callee,
				},
			})
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"nodes": nodes,
			"edges": edges,
		})
	})

	_ = http.ListenAndServe(":8080", nil)
}
