package languages

import (
	"fmt"
	sitter "github.com/tree-sitter/go-tree-sitter"
	"strings"
)

// ExtractJSOutline extracts JavaScript outline directly from the code
func ExtractJSOutline(root *sitter.Node, content []byte) string {
	var result strings.Builder

	// Function to process a node and its children
	var processNode func(node *sitter.Node, indentLevel int)
	processNode = func(node *sitter.Node, indentLevel int) {
		indent := strings.Repeat(" ", indentLevel*2)

		// Process based on node type
		switch node.Kind() {
		case "program":
			// Process all children
			for i := 0; i < int(node.NamedChildCount()); i++ {
				child := node.NamedChild(uint(i))
				processNode(child, indentLevel)
			}

		case "import_statement":
			// Handle import statements
			importText := getNodeText(node, content)
			result.WriteString(fmt.Sprintf("%s\n", importText))

		case "function_declaration", "generator_function_declaration":
			// For JavaScript function declarations
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// Get parameters
				paramNode := node.ChildByFieldName("parameters")
				paramText := ""
				if paramNode != nil {
					paramText = getNodeText(paramNode, content)
				}

				// Get documentation comment (JSDoc) if present
				doc := findDocComment(node, content, "javascript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				// Write function declaration
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%sfunction %s%s { // line %d\n", indent, name, paramText, lineNum))
				result.WriteString(fmt.Sprintf("%s  // ...\n", indent))
				result.WriteString(fmt.Sprintf("%s}\n\n", indent))
			}

		case "method_definition":
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// Skip private methods (those starting with #)
				if strings.HasPrefix(name, "#") {
					return
				}

				// Get parameters
				paramNode := node.ChildByFieldName("parameters")
				paramText := ""
				if paramNode != nil {
					paramText = getNodeText(paramNode, content)
				}

				// Check if it's a static method
				isStatic := false
				for j := 0; j < int(node.ChildCount()); j++ {
					if node.Child(uint(j)).Kind() == "static" {
						isStatic = true
						break
					}
				}

				prefix := ""
				if isStatic {
					prefix = "static "
				}

				// Get documentation comment if present
				doc := findDocComment(node, content, "javascript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				// Write method definition
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%s%s%s%s { // line %d\n", indent, prefix, name, paramText, lineNum))
				result.WriteString(fmt.Sprintf("%s  // ...\n", indent))
				result.WriteString(fmt.Sprintf("%s}\n\n", indent))
			}

		case "class_declaration":
			// For JavaScript classes
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// Get extends clause if any
				var extendsText string
				for i := 0; i < int(node.ChildCount()); i++ {
					child := node.Child(uint(i))
					if child.Kind() == "class_heritage" {
						extendsText = " " + getNodeText(child, content)
						break
					}
				}

				// Get documentation comment if present
				doc := findDocComment(node, content, "javascript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				// Write class declaration
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%sclass %s%s { // line %d\n", indent, name, extendsText, lineNum))

				// Process class body
				bodyNode := node.ChildByFieldName("body")
				if bodyNode != nil {
					for i := 0; i < int(bodyNode.NamedChildCount()); i++ {
						child := bodyNode.NamedChild(uint(i))
						processNode(child, indentLevel+1)
					}
				}

				result.WriteString(fmt.Sprintf("%s}\n\n", indent))
			}

		case "export_statement":
			// Handle different types of export statements
			if node.NamedChildCount() > 0 {
				firstChild := node.NamedChild(0)

				// Check if it's a default export
				isDefault := false
				for i := 0; i < int(node.ChildCount()); i++ {
					if node.Child(uint(i)).Kind() == "default" {
						isDefault = true
						break
					}
				}

				switch firstChild.Kind() {
				case "function_declaration", "generator_function_declaration":
					nameNode := firstChild.ChildByFieldName("name")
					if nameNode != nil {
						name := getNodeText(nameNode, content)

						// Get parameters
						paramNode := firstChild.ChildByFieldName("parameters")
						paramText := ""
						if paramNode != nil {
							paramText = getNodeText(paramNode, content)
						}

						// Get documentation comment if present
						doc := findDocComment(node, content, "javascript")
						if doc != "" {
							docLines := strings.Split(doc, "\n")
							for _, line := range docLines {
								result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
							}
						}

						// Write export function declaration
						lineNum := getNodeLineNumber(firstChild)
						if isDefault {
							result.WriteString(fmt.Sprintf("%sexport default function %s%s { // line %d\n", indent, name, paramText, lineNum))
						} else {
							result.WriteString(fmt.Sprintf("%sexport function %s%s { // line %d\n", indent, name, paramText, lineNum))
						}
						result.WriteString(fmt.Sprintf("%s  // ...\n", indent))
						result.WriteString(fmt.Sprintf("%s}\n\n", indent))
					}

				case "class_declaration":
					nameNode := firstChild.ChildByFieldName("name")
					if nameNode != nil {
						name := getNodeText(nameNode, content)

						// Get extends clause if any
						var extendsText string
						for i := 0; i < int(firstChild.ChildCount()); i++ {
							child := firstChild.Child(uint(i))
							if child.Kind() == "class_heritage" {
								extendsText = " " + getNodeText(child, content)
								break
							}
						}

						// Get documentation comment if present
						doc := findDocComment(node, content, "javascript")
						if doc != "" {
							docLines := strings.Split(doc, "\n")
							for _, line := range docLines {
								result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
							}
						}

						// Write export class declaration
						lineNum := getNodeLineNumber(firstChild)
						if isDefault {
							result.WriteString(fmt.Sprintf("%sexport default class %s%s { // line %d\n", indent, name, extendsText, lineNum))
						} else {
							result.WriteString(fmt.Sprintf("%sexport class %s%s { // line %d\n", indent, name, extendsText, lineNum))
						}

						// Process class body
						bodyNode := firstChild.ChildByFieldName("body")
						if bodyNode != nil {
							for i := 0; i < int(bodyNode.NamedChildCount()); i++ {
								child := bodyNode.NamedChild(uint(i))
								processNode(child, indentLevel+1)
							}
						}

						result.WriteString(fmt.Sprintf("%s}\n\n", indent))
					}

				case "lexical_declaration", "variable_declaration":
					// Handle export const/let/var declarations
					if firstChild.NamedChildCount() > 0 {
						for i := 0; i < int(firstChild.NamedChildCount()); i++ {
							declarator := firstChild.NamedChild(uint(i))
							if declarator.Kind() == "variable_declarator" && declarator.NamedChildCount() >= 2 {
								nameNode := declarator.NamedChild(0)
								valueNode := declarator.NamedChild(1)

								if valueNode.Kind() == "arrow_function" || valueNode.Kind() == "function" {
									name := getNodeText(nameNode, content)

									// Get declaration type
									declType := "var"
									if firstChild.Kind() == "lexical_declaration" {
										if firstChild.Child(0).Kind() == "let" {
											declType = "let"
										} else {
											declType = "const"
										}
									}

									// Get parameters
									paramNode := valueNode.ChildByFieldName("parameters")
									paramText := ""
									if paramNode != nil {
										paramText = getNodeText(paramNode, content)
									}

									// Get documentation comment if present
									doc := findDocComment(node, content, "javascript")
									if doc != "" {
										docLines := strings.Split(doc, "\n")
										for _, line := range docLines {
											result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
										}
									}

									// Write export function
									lineNum := getNodeLineNumber(firstChild)
									if valueNode.Kind() == "arrow_function" {
										result.WriteString(fmt.Sprintf("%sexport %s %s = %s => { // line %d\n", indent, declType, name, paramText, lineNum))
									} else {
										result.WriteString(fmt.Sprintf("%sexport %s %s = function%s { // line %d\n", indent, declType, name, paramText, lineNum))
									}
									result.WriteString(fmt.Sprintf("%s  // ...\n", indent))
									result.WriteString(fmt.Sprintf("%s}\n\n", indent))
								} else {
									// Handle other exported variable declarations
									name := getNodeText(nameNode, content)
									declType := "var"
									if firstChild.Kind() == "lexical_declaration" {
										if firstChild.Child(0).Kind() == "let" {
											declType = "let"
										} else {
											declType = "const"
										}
									}
									lineNum := getNodeLineNumber(firstChild)
									result.WriteString(fmt.Sprintf("%sexport %s %s; // line %d\n\n", indent, declType, name, lineNum))
								}
							}
						}
					}

				case "export_clause":
					// Handle export { ... } statements
					exportText := getNodeText(node, content)
					lineNum := getNodeLineNumber(node)
					result.WriteString(fmt.Sprintf("%s%s // line %d\n\n", indent, exportText, lineNum))

				default:
					// Handle other export patterns like export * from '...'
					exportText := getNodeText(node, content)
					lineNum := getNodeLineNumber(node)
					result.WriteString(fmt.Sprintf("%s%s // line %d\n\n", indent, exportText, lineNum))
				}
			} else {
				// Fallback for other export patterns
				exportText := getNodeText(node, content)
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%s%s // line %d\n\n", indent, exportText, lineNum))
			}

		case "lexical_declaration", "variable_declaration":
			// For variable declarations that might contain arrow functions or require calls
			if node.NamedChildCount() > 0 {
				for i := 0; i < int(node.NamedChildCount()); i++ {
					declarator := node.NamedChild(uint(i))
					if declarator.Kind() == "variable_declarator" && declarator.NamedChildCount() >= 2 {
						nameNode := declarator.NamedChild(0)
						valueNode := declarator.NamedChild(1)

						if valueNode.Kind() == "arrow_function" || valueNode.Kind() == "function" {
							name := getNodeText(nameNode, content)

							// Get declaration type
							declType := "var"
							if node.Kind() == "lexical_declaration" {
								if node.Child(0).Kind() == "let" {
									declType = "let"
								} else {
									declType = "const"
								}
							}

							// Get parameters
							paramNode := valueNode.ChildByFieldName("parameters")
							paramText := ""
							if paramNode != nil {
								paramText = getNodeText(paramNode, content)
							}

							// Get documentation comment if present
							doc := findDocComment(node, content, "javascript")
							if doc != "" {
								docLines := strings.Split(doc, "\n")
								for _, line := range docLines {
									result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
								}
							}

							// Write function
							lineNum := getNodeLineNumber(node)
							if valueNode.Kind() == "arrow_function" {
								result.WriteString(fmt.Sprintf("%s%s %s = %s => { // line %d\n", indent, declType, name, paramText, lineNum))
							} else {
								result.WriteString(fmt.Sprintf("%s%s %s = function%s { // line %d\n", indent, declType, name, paramText, lineNum))
							}
							result.WriteString(fmt.Sprintf("%s  // ...\n", indent))
							result.WriteString(fmt.Sprintf("%s}\n\n", indent))
						} else if valueNode.Kind() == "call_expression" {
							// Check if this is a require() call
							functionNode := valueNode.ChildByFieldName("function")
							if functionNode != nil && getNodeText(functionNode, content) == "require" {
								// This is a require statement, include it in the outline
								requireText := getNodeText(node, content)
								result.WriteString(fmt.Sprintf("%s\n", requireText))
							}
						}
					}
				}
			}
		}
	}

	processNode(root, 0)
	return result.String()
}
