package agentconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func InstallGeminiCLIConfig() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not find home directory: %v", err)
	}

	geminiDir := filepath.Join(homeDir, ".gemini")
	if err := os.MkdirAll(geminiDir, 0755); err != nil {
		return fmt.Errorf("could not create .gemini directory: %v", err)
	}

	settingsPath := filepath.Join(geminiDir, "settings.json")
	var settings map[string]interface{}

	fileData, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			settings = make(map[string]interface{})
		} else {
			return fmt.Errorf("could not read settings.json: %v", err)
		}
	} else {
		if err := json.Unmarshal(fileData, &settings); err != nil {
			settings = make(map[string]interface{})
		}
	}

	if _, exists := settings["mcpServers"]; !exists {
		settings["mcpServers"] = make(map[string]interface{})
	}

	mcpServers, ok := settings["mcpServers"].(map[string]interface{})
	if !ok {
		mcpServers = make(map[string]interface{})
		settings["mcpServers"] = mcpServers
	}

	executablePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("could not get executable path: %v", err)
	}

	mcpServers["micro-mcp"] = map[string]interface{}{
		"command": executablePath,
		"args":    []string{},
		"env":     map[string]string{},
	}

	updatedData, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("could not marshal settings: %v", err)
	}

	if err := os.WriteFile(settingsPath, updatedData, 0644); err != nil {
		return fmt.Errorf("could not write settings.json: %v", err)
	}

	instructionsPath := filepath.Join(geminiDir, "GEMINI.md")
	instructions := `# micro-mcp Integration
When exploring codebases, always use the 'search_symbols', 'get_file_dependencies', and 'read_files' tools provided by the 'micro-mcp' MCP server to gain deep structural context before attempting any architectural changes.
`
	if err := os.WriteFile(instructionsPath, []byte(instructions), 0644); err != nil {
		return fmt.Errorf("could not write GEMINI.md: %v", err)
	}

	fmt.Println("Successfully configured Gemini CLI for micro-mcp.")
	return nil
}
