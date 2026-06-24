package reader

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ReadFilesContext(paths []string) string {
	var fileContents []string
	for _, path := range paths {
		content, err := os.ReadFile(strings.Trim(path, `"' `))
		if err == nil {
			fileContents = append(fileContents, fmt.Sprintf("--- FILE: %s ---\n%s\n", filepath.Base(path), string(content)))
		}
	}
	return strings.Join(fileContents, "\n")
}
