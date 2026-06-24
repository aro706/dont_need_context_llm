package parser

import (
	"fmt"
	"micro-mcp/internal/models"
	"path/filepath"

	sitter "github.com/smacker/go-tree-sitter"
)

func ProcessFunctions(node *sitter.Node, source []byte, filePath, projectName, currentCaller string, nodes *[]models.Node) string {

	if node.Type() == "import_statement" {
		*nodes = append(*nodes, models.Node{ProjectName: projectName, FilePath: filepath.Base(filePath), FunctionName: node.Content(source), Type: "import"})
	}

	if node.Type() == "variable_declarator" {
		nameNode := node.ChildByFieldName("name")
		initNode := node.ChildByFieldName("value")
		if nameNode != nil {
			*nodes = append(*nodes, models.Node{ProjectName: projectName, FilePath: filepath.Base(filePath), FunctionName: nameNode.Content(source), Type: "variable"})
			if initNode != nil && initNode.Type() == "call_expression" {
				callee := initNode.ChildByFieldName("function")
				if callee != nil && callee.Content(source) == "require" {
					*nodes = append(*nodes, models.Node{ProjectName: projectName, FilePath: filepath.Base(filePath), FunctionName: fmt.Sprintf("require: %s", initNode.Content(source)), Type: "import"})
				}
			}
		}
	}

	if node.Type() == "function_declaration" || node.Type() == "arrow_function" {
		var funcName string
		if node.Type() == "function_declaration" {
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				funcName = nameNode.Content(source)
			}
		} else if node.Type() == "arrow_function" {
			parent := node.Parent()
			if parent != nil && parent.Type() == "variable_declarator" {
				nameNode := parent.ChildByFieldName("name")
				if nameNode != nil {
					funcName = nameNode.Content(source)
				}
			}
		}

		if funcName != "" {
			*nodes = append(*nodes, models.Node{ProjectName: projectName, FilePath: filepath.Base(filePath), FunctionName: funcName, Type: "function"})
			return funcName
		}
	}

	return currentCaller
}
