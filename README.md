# micro-mcp

A high-performance, Go-native codebase intelligence engine for AI agents.

micro-mcp indexes repositories in milliseconds, maps structural dependencies, and provides an interactive graph visualization of your codebase. Built entirely in Go, it delivers cross-platform reliability, fast startup times, and minimal resource usage.

---

## Overview

micro-mcp is a modern codebase intelligence platform powered by Tree-sitter and an in-memory SQLite database. It builds a structured knowledge graph of files, functions, API routes, and call relationships, enabling AI agents to understand large codebases efficiently.

### Key Capabilities

* Fast repository indexing with recursive directory traversal
* Intelligent exclusion of directories such as `.git`, `node_modules`, and build artifacts
* Native support for the Model Context Protocol (MCP)
* Interactive architecture visualization through a built-in web interface
* Multi-project workspace management with isolated knowledge graphs
* Real-time synchronization using filesystem watchers

---

## Architecture

```text
internal/
├── agentconfig/    # MCP and client configuration management
├── api/            # HTTP server and graph APIs
├── database/       # In-memory SQLite persistence layer
├── indexer/        # Repository traversal and indexing orchestration
├── mcp/            # JSON-RPC 2.0 MCP implementation
├── parser/         # Tree-sitter based code analysis
├── watcher/        # Real-time filesystem monitoring

cmd/server/         # Application entry point
frontend/           # Graph visualization interface
```

---

## Features

### Codebase Intelligence

* Automatic discovery of files, functions, classes, and routes
* Dependency analysis across modules and services
* Call graph generation for understanding execution flow
* Cross-file symbol relationships and reference tracking

### Interactive Visualization

* Graph-based architecture exploration
* Project-wide dependency navigation
* Workspace selection and project isolation
* Lightweight browser-based interface

### MCP Integration

The server exposes a set of MCP tools for AI agents:

#### `index_project`

Indexes a repository and builds its knowledge graph.

```json
{
  "path": "/absolute/path/to/project"
}
```

#### `search_symbols`

Searches for files, functions, routes, or symbols.

#### `get_file_dependencies`

Returns dependency relationships and call chains.

#### `read_files`

Provides context-aware file access for agent workflows.

---

## Getting Started

### Prerequisites

* Go 1.22 or later

### Build

```bash
go build -o micro-mcp.exe ./cmd/server
```

### Install MCP Configuration

```bash
./micro-mcp.exe install
```

### Run

```bash
go run ./cmd/server
```

---

## Usage

### Index a Project

Call the MCP tool:

```json
{
  "path": "/absolute/path/to/project"
}
```

### Open the Visualization

Navigate to:

```text
http://localhost:8080
```

Select a project from the workspace list to explore its architecture and dependency graph.

---

## How It Works

1. The indexer scans the repository.
2. Tree-sitter parses source files into ASTs.
3. Symbols, routes, and dependencies are extracted.
4. Data is stored in an in-memory SQLite knowledge graph.
5. MCP tools expose the graph to AI agents.
6. The web interface visualizes relationships interactively.
7. Filesystem watchers keep the graph synchronized automatically.

---

## Configuration

micro-mcp automatically:

* Detects file additions, modifications, and deletions
* Reindexes affected components
* Updates dependency relationships
* Maintains isolated workspaces for multiple projects

No external database or additional services are required.

---

## License

MIT License
