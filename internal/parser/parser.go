package parser

import (
	"context"
	"micro-mcp/internal/models"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

var astParser *sitter.Parser

func InitParser() {
	astParser = sitter.NewParser()
	astParser.SetLanguage(javascript.GetLanguage())
}

// ParseFile takes raw code and returns identified Nodes and Calls
func ParseFile(source []byte, filePath string, projectName string) ([]models.Node, []models.Call) {
	tree, _ := astParser.ParseCtx(context.Background(), nil, source)

	var nodes []models.Node
	var calls []models.Call

	walkAST(tree.RootNode(), source, filePath, projectName, "", &nodes, &calls)

	return nodes, calls
}

func walkAST(node *sitter.Node, source []byte, filePath, projectName, currentCaller string, nodes *[]models.Node, calls *[]models.Call) {

	// 1. Process standard functions/variables
	newCaller := ProcessFunctions(node, source, filePath, projectName, currentCaller, nodes)

	// 2. Process specific routing and API calls
	ProcessRoutes(node, source, filePath, projectName, currentCaller, nodes, calls)

	for i := 0; i < int(node.ChildCount()); i++ {
		walkAST(node.Child(i), source, filePath, projectName, newCaller, nodes, calls)
	}
}
