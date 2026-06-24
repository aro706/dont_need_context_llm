package main

import (
	"fmt"
	"os"

	"micro-mcp/internal/agentconfig"
	"micro-mcp/internal/api"
	"micro-mcp/internal/database"
	"micro-mcp/internal/mcp"
	"micro-mcp/internal/parser"
	"micro-mcp/internal/watcher"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "install" {
		err := agentconfig.InstallGeminiCLIConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Installation failed: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	db, err := database.InitDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, "DB Error: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	parser.InitParser()

	fsWatcher, _ := fsnotify.NewWatcher()

	go watcher.StartWatchLoop(db, fsWatcher)
	go api.StartHTTPServer(db)

	mcp.StartStdioLoop(db, fsWatcher)
}
