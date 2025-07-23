package languages

import (
	"fmt"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// ExtractPythonOutline extracts Python outline directly from the code
func ExtractPythonOutline(root *sitter.Node, content []byte) string {
	var result strings.Builder

	// Function to process a node and its children
	var processNode func(node *sitter.Node, indentLevel int)
	processNode = func(node *sitter.Node, indentLevel int) {
		indent := strings.Repeat(" ", indentLevel*4)

		// Process based on node type
		switch node.Kind() {
		case "module":
			// Process all children
			for i := 0; i < int(node.NamedChildCount()); i++ {
				child := node.NamedChild(uint(i))
				processNode(child, indentLevel)
			}

		case "import_statement", "import_from_statement":
			// Handle import statements (both 'import' and 'from ... import')
			importText := getNodeText(node, content)
			result.WriteString(fmt.Sprintf("%s\n", importText))

		case "function_definition":
			// For Python functions
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// In Python, names starting with _ are considered private
				isPublic := !strings.HasPrefix(name, "_")
				if !isPublic {
					return
				}

				// Get the parameter list
				paramNode := node.ChildByFieldName("parameters")
				paramText := ""
				if paramNode != nil {
					paramText = getNodeText(paramNode, content)
				}

				// Get return type annotation if any
				returnNode := node.ChildByFieldName("return_type")
				returnText := ""
				if returnNode != nil {
					returnText = " -> " + getNodeText(returnNode, content)
				}

				// Find documentation (Python uses docstrings inside the function body)
				doc := ""
				bodyNode := node.ChildByFieldName("body")
				if bodyNode != nil && bodyNode.NamedChildCount() > 0 {
					firstChild := bodyNode.NamedChild(0)
					if firstChild.Kind() == "expression_statement" && firstChild.NamedChildCount() > 0 {
						exprChild := firstChild.NamedChild(0)
						if exprChild.Kind() == "string" {
							doc = getNodeText(exprChild, content)
							// Clean up docstring
							doc = strings.Trim(doc, "\"'")
						}
					}
				}

				// Write function definition
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%sdef %s%s%s: # line %d", indent, name, paramText, returnText, lineNum))
				if doc != "" {
					result.WriteString(fmt.Sprintf(" \"\"\"%s\"\"\"", doc))
				}
				result.WriteString("\n")
				result.WriteString(fmt.Sprintf("%s    ...\n\n", indent))
			}

		case "class_definition":
			// For Python classes
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// In Python, names starting with _ are considered private
				isPublic := !strings.HasPrefix(name, "_")
				if !isPublic {
					return
				}

				// Get base classes if any
				baseNode := node.ChildByFieldName("base_clause")
				baseText := ""
				if baseNode != nil {
					baseText = getNodeText(baseNode, content)
				}

				// Write class definition
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%sclass %s%s: # line %d\n", indent, name, baseText, lineNum))

				// Find documentation (Python uses docstrings inside the class body)
				doc := ""
				bodyNode := node.ChildByFieldName("body")
				if bodyNode != nil && bodyNode.NamedChildCount() > 0 {
					firstChild := bodyNode.NamedChild(0)
					if firstChild.Kind() == "expression_statement" && firstChild.NamedChildCount() > 0 {
						exprChild := firstChild.NamedChild(0)
						if exprChild.Kind() == "string" {
							doc = getNodeText(exprChild, content)
							doc = strings.Trim(doc, "\"'")
							result.WriteString(fmt.Sprintf("%s    \"\"\"%s\"\"\"\n", indent, doc))
						}
					}
				}

				// Process class body for methods
				hasMethods := false
				if bodyNode != nil {
					for i := 0; i < int(bodyNode.NamedChildCount()); i++ {
						child := bodyNode.NamedChild(uint(i))
						if child.Kind() == "function_definition" {
							methodNameNode := child.ChildByFieldName("name")
							if methodNameNode != nil {
								methodName := getNodeText(methodNameNode, content)

								// Skip private methods
								if strings.HasPrefix(methodName, "_") {
									continue
								}

								hasMethods = true

								// Get parameters
								paramNode := child.ChildByFieldName("parameters")
								paramText := ""
								if paramNode != nil {
									paramText = getNodeText(paramNode, content)
								}

								// Get return type
								returnNode := child.ChildByFieldName("return_type")
								returnText := ""
								if returnNode != nil {
									returnText = " -> " + getNodeText(returnNode, content)
								}

								// Get docstring
								methodDoc := ""
								methodBodyNode := child.ChildByFieldName("body")
								if methodBodyNode != nil && methodBodyNode.NamedChildCount() > 0 {
									firstChild := methodBodyNode.NamedChild(0)
									if firstChild.Kind() == "expression_statement" && firstChild.NamedChildCount() > 0 {
										exprChild := firstChild.NamedChild(0)
										if exprChild.Kind() == "string" {
											methodDoc = getNodeText(exprChild, content)
											methodDoc = strings.Trim(methodDoc, "\"'")
										}
									}
								}

								// Write method definition
								methodLineNum := getNodeLineNumber(child)
								result.WriteString(fmt.Sprintf("%s    def %s%s%s: # line %d", indent, methodName, paramText, returnText, methodLineNum))
								if methodDoc != "" {
									result.WriteString(fmt.Sprintf(" \"\"\"%s\"\"\"", methodDoc))
								}
								result.WriteString("\n")
								result.WriteString(fmt.Sprintf("%s        ...\n\n", indent))
							}
						}
					}
				}

				// If no methods were found, add 'pass'
				if !hasMethods {
					result.WriteString(fmt.Sprintf("%s    pass\n\n", indent))
				} else {
					result.WriteString("\n")
				}
			}
		}
	}

	processNode(root, 0)
	return result.String()
}
