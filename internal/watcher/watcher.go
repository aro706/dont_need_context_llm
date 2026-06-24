package watcher

import (
	"database/sql"
	"strings"

	"micro-mcp/internal/database"
	"micro-mcp/internal/indexer"

	"github.com/fsnotify/fsnotify"
)

func StartWatchLoop(db *sql.DB, watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if strings.HasSuffix(event.Name, ".js") && (event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create) {
				// Wipe DB and Resync all projects
				database.ClearDB(db)

				projects := indexer.GetAllProjects()
				for path, name := range projects {
					indexer.WalkAndIndex(db, path, name, watcher)
				}
			}
		case <-watcher.Errors:
		}
	}
}
