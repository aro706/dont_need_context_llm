package mcp

import (
	"database/sql"
	"encoding/json"
	"path/filepath"
	"strings"

	"micro-mcp/internal/database"
	"micro-mcp/internal/indexer"
	"micro-mcp/internal/models"
	"micro-mcp/internal/reader"

	"github.com/fsnotify/fsnotify"
)

// HandleToolCall routes the request to the correct internal system
func HandleToolCall(db *sql.DB, watcher *fsnotify.Watcher, req models.Request, resp *models.Response) {
	var callReq struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}
	json.Unmarshal(req.Params, &callReq)

	switch callReq.Name {

	case "index_project":
		var params models.IndexParams
		json.Unmarshal(callReq.Arguments, &params)
		cleanPath := strings.Trim(params.Path, `"'\ `)
		projectName := filepath.Base(cleanPath)

		indexer.AddProject(cleanPath, projectName)
		indexer.WalkAndIndex(db, cleanPath, projectName, watcher)

		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": "Successfully indexed project: " + projectName}},
		}

	case "search_symbols":
		var params models.SearchParams
		json.Unmarshal(callReq.Arguments, &params)
		resultText := database.SearchSymbols(db, params.Query, params.Project)

		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": resultText}},
		}

	case "get_file_dependencies":
		var params struct {
			Project string `json:"project"`
		}
		json.Unmarshal(callReq.Arguments, &params)
		resultText := database.GetDependencies(db, params.Project)

		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": resultText}},
		}

	case "read_files":
		var params struct {
			Paths []string `json:"paths"`
		}
		json.Unmarshal(callReq.Arguments, &params)
		resultText := reader.ReadFilesContext(params.Paths)

		resp.Result = map[string]interface{}{
			"content": []map[string]interface{}{{"type": "text", "text": resultText}},
		}
	}
}
