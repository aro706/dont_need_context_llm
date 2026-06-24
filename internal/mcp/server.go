package mcp

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"micro-mcp/internal/models"

	"github.com/fsnotify/fsnotify"
)

func StartStdioLoop(db *sql.DB, watcher *fsnotify.Watcher) {
	reader := bufio.NewReader(os.Stdin)
	for {
		input, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		var req models.Request
		if err := json.Unmarshal([]byte(input), &req); err != nil {
			continue
		}

		// --- THE FIX ---
		// JSON-RPC specifies that requests without an ID are "notifications".
		// A server MUST NOT reply to notifications, so we silently ignore them.
		if req.ID == nil {
			continue
		}
		// ---------------

		var resp models.Response
		resp.JSONRPC = "2.0"
		resp.ID = req.ID

		switch req.Method {
		case "initialize":
			resp.Result = map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"serverInfo":      map[string]string{"name": "micro-mcp-go", "version": "1.0.0"},
				"capabilities":    map[string]interface{}{"tools": map[string]bool{"listChanged": false}},
			}
		case "tools/list":
			resp.Result = GetToolList()
		case "tools/call":
			HandleToolCall(db, watcher, req, &resp)
		default:
			resp.Error = map[string]string{"code": "-32601", "message": "Method not found"}
		}

		output, _ := json.Marshal(resp)
		fmt.Println(string(output))
	}
}

// ... keep your GetToolList() function below exactly as it is
// Kept isolated to keep server.go clean
func GetToolList() map[string]interface{} {
	return map[string]interface{}{
		"tools": []map[string]interface{}{
			{
				"name":        "index_project",
				"description": "Indexes a specific local project directory into the system database.",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"path": map[string]interface{}{"type": "string", "description": "The absolute path of the project folder."},
					},
					"required": []string{"path"},
				},
			},
			{
				"name":        "search_symbols",
				"description": "Searches for structural symbols, functions, or database references across projects.",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"query":   map[string]interface{}{"type": "string", "description": "Term like 'vector', 'embedding'."},
						"project": map[string]interface{}{"type": "string", "description": "Optional project filter."},
					},
					"required": []string{"query"},
				},
			},
			{
				"name":        "get_file_dependencies",
				"description": "Returns all caller and callee relationships for functions.",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"project": map[string]interface{}{"type": "string", "description": "Name of the project folder."},
					},
					"required": []string{"project"},
				},
			},
			{
				"name":        "read_files",
				"description": "Reads the exact content of specific files.",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"paths": map[string]interface{}{
							"type":  "array",
							"items": map[string]interface{}{"type": "string"},
						},
					},
					"required": []string{"paths"},
				},
			},
		},
	}
}
