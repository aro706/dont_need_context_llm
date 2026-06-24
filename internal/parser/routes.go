package parser

import (
	"fmt"
	"micro-mcp/internal/models"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

func ProcessRoutes(node *sitter.Node, source []byte, filePath, projectName, currentCaller string, nodes *[]models.Node, calls *[]models.Call) {
	if node.Type() == "call_expression" {
		funcNode := node.ChildByFieldName("function")
		if funcNode != nil {
			funcContent := funcNode.Content(source)

			// Detect Express and API calls
			if funcNode.Type() == "member_expression" {
				objNode := funcNode.ChildByFieldName("object")
				propNode := funcNode.ChildByFieldName("property")
				if objNode != nil && propNode != nil {
					methodName := propNode.Content(source)

					if methodName == "get" || methodName == "post" || methodName == "put" || methodName == "delete" || methodName == "patch" {
						argsNode := node.ChildByFieldName("arguments")
						if argsNode != nil && argsNode.ChildCount() > 1 {
							firstArg := argsNode.Child(1)
							if firstArg != nil && (firstArg.Type() == "string" || firstArg.Type() == "string_fragment") {
								pathStr := strings.Trim(firstArg.Content(source), "\"'`")
								if (strings.HasPrefix(pathStr, "/") || len(pathStr) > 0) && !strings.Contains(pathStr, " ") {
									routeName := fmt.Sprintf("%s %s", strings.ToUpper(methodName), pathStr)
									*nodes = append(*nodes, models.Node{ProjectName: projectName, FilePath: filepath.Base(filePath), FunctionName: routeName, Type: "express_endpoint"})
								}
							}
						}
					}

					if objNode.Content(source) == "axios" || objNode.Content(source) == "fetch" || methodName == "get" || methodName == "post" {
						argsNode := node.ChildByFieldName("arguments")
						if argsNode != nil && argsNode.ChildCount() > 1 {
							firstArg := argsNode.Child(1)
							if firstArg != nil && (firstArg.Type() == "string" || firstArg.Type() == "string_fragment") {
								pathStr := strings.Trim(firstArg.Content(source), "\"'`")
								if strings.Contains(pathStr, "/") {
									apiCall := fmt.Sprintf("API_CALL %s %s", strings.ToUpper(methodName), pathStr)
									*nodes = append(*nodes, models.Node{ProjectName: projectName, FilePath: filepath.Base(filePath), FunctionName: apiCall, Type: "backend_request"})
								}
							}
						}
					}
				}
			}

			if funcContent == "fetch" {
				argsNode := node.ChildByFieldName("arguments")
				if argsNode != nil && argsNode.ChildCount() > 1 {
					firstArg := argsNode.Child(1)
					if firstArg != nil && (firstArg.Type() == "string" || firstArg.Type() == "string_fragment") {
						pathStr := strings.Trim(firstArg.Content(source), "\"'`")
						apiCall := fmt.Sprintf("API_CALL FETCH %s", pathStr)
						*nodes = append(*nodes, models.Node{ProjectName: projectName, FilePath: filepath.Base(filePath), FunctionName: apiCall, Type: "backend_request"})
					}
				}
			}

			// Add Call Dependency Edge
			if currentCaller != "" {
				calleeName := funcContent
				if funcNode.Type() == "member_expression" {
					propNode := funcNode.ChildByFieldName("property")
					if propNode != nil {
						calleeName = propNode.Content(source)
					}
				}
				*calls = append(*calls, models.Call{ProjectName: projectName, Caller: currentCaller, Callee: calleeName})
			}
		}
	}
}
