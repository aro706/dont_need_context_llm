package indexer

import (
	"database/sql"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"micro-mcp/internal/database"
	"micro-mcp/internal/parser"

	"github.com/fsnotify/fsnotify"
)

func WalkAndIndex(db *sql.DB, targetDir string, projectName string, watcher *fsnotify.Watcher) {
	filepath.Walk(targetDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil || strings.Contains(path, "node_modules") || strings.Contains(path, ".git") {
			return nil
		}

		if info.IsDir() && watcher != nil {
			watcher.Add(path)
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".js") {
			sourceCode, err := os.ReadFile(path)
			if err == nil {
				// 1. Get parsed data from AST layer
				nodes, calls := parser.ParseFile(sourceCode, path, projectName)

				// 2. Save parsed data using Database layer
				for _, n := range nodes {
					database.InsertNode(db, n)
				}
				for _, c := range calls {
					database.InsertCall(db, c)
				}
			}
		}
		return nil
	})
}
