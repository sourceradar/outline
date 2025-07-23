package languages

import (
	"fmt"
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// ExtractTSOutline extracts TypeScript outline directly from the code
func ExtractTSOutline(root *sitter.Node, content []byte) string {
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
			// For TypeScript functions
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// Get parameters
				paramNode := node.ChildByFieldName("parameters")
				paramText := ""
				if paramNode != nil {
					paramText = getNodeText(paramNode, content)
				}

				// Get return type if any
				returnNode := node.ChildByFieldName("return_type")
				returnText := ""
				if returnNode != nil {
					returnText = getNodeText(returnNode, content)
				}

				// Get documentation comment if present
				doc := findDocComment(node, content, "typescript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				// Write function declaration
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%sfunction %s%s%s { // line %d\n", indent, name, paramText, returnText, lineNum))
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

				// Get return type if any
				returnNode := node.ChildByFieldName("return_type")
				returnText := ""
				if returnNode != nil {
					returnText = getNodeText(returnNode, content)
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
				doc := findDocComment(node, content, "typescript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				// Write method definition
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%s%s%s%s%s { // line %d\n", indent, prefix, name, paramText, returnText, lineNum))
				result.WriteString(fmt.Sprintf("%s  // ...\n", indent))
				result.WriteString(fmt.Sprintf("%s}\n\n", indent))
			}

		case "class_declaration":
			// For TypeScript classes
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// Get heritage clause (extends/implements)
				var heritageText string
				for i := 0; i < int(node.ChildCount()); i++ {
					child := node.Child(uint(i))
					if child.Kind() == "class_heritage" {
						heritageText = " " + getNodeText(child, content)
						break
					}
				}

				// Get documentation comment if present
				doc := findDocComment(node, content, "typescript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				// Write class declaration
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%sclass %s%s { // line %d\n", indent, name, heritageText, lineNum))

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

		case "interface_declaration":
			// For TypeScript interfaces
			nameNode := node.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)

				// Get extends clause if any
				extendsNode := node.ChildByFieldName("extends_clause")
				extendsText := ""
				if extendsNode != nil {
					extendsText = " " + getNodeText(extendsNode, content)
				}

				// Get documentation comment if present
				doc := findDocComment(node, content, "typescript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				// Write interface declaration
				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%sinterface %s%s { // line %d\n", indent, name, extendsText, lineNum))

				// Process interface body for property and method signatures
				bodyNode := node.ChildByFieldName("body")
				if bodyNode != nil {
					for i := 0; i < int(bodyNode.NamedChildCount()); i++ {
						child := bodyNode.NamedChild(uint(i))

						if child.Kind() == "property_signature" {
							nameNode := child.ChildByFieldName("name")
							typeNode := child.ChildByFieldName("type")

							if nameNode != nil && typeNode != nil {
								propName := getNodeText(nameNode, content)
								propType := getNodeText(typeNode, content)

								// Check for optional marker
								optional := ""
								for j := 0; j < int(child.ChildCount()); j++ {
									if child.Child(uint(j)).Kind() == "?" {
										optional = "?"
										break
									}
								}

								// Get doc comment
								propDoc := findDocComment(child, content, "typescript")
								if propDoc != "" {
									result.WriteString(fmt.Sprintf("%s  // %s\n", indent, propDoc))
								}

								result.WriteString(fmt.Sprintf("%s  %s%s: %s;\n", indent, propName, optional, propType))
							}
						} else if child.Kind() == "method_signature" {
							nameNode := child.ChildByFieldName("name")
							paramNode := child.ChildByFieldName("parameters")
							returnNode := child.ChildByFieldName("return_type")

							if nameNode != nil {
								methodName := getNodeText(nameNode, content)

								paramText := ""
								if paramNode != nil {
									paramText = getNodeText(paramNode, content)
								}

								returnText := ""
								if returnNode != nil {
									returnText = ": " + getNodeText(returnNode, content)
								}

								// Get doc comment
								methodDoc := findDocComment(child, content, "typescript")
								if methodDoc != "" {
									result.WriteString(fmt.Sprintf("%s  // %s\n", indent, methodDoc))
								}

								result.WriteString(fmt.Sprintf("%s  %s%s%s;\n", indent, methodName, paramText, returnText))
							}
						}
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
						
						// Get return type if any
						returnNode := firstChild.ChildByFieldName("return_type")
						returnText := ""
						if returnNode != nil {
							returnText = getNodeText(returnNode, content)
						}
						
						// Get documentation comment if present
						doc := findDocComment(node, content, "typescript")
						if doc != "" {
							docLines := strings.Split(doc, "\n")
							for _, line := range docLines {
								result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
							}
						}
						
						// Write export function declaration
						lineNum := getNodeLineNumber(firstChild)
						if isDefault {
							result.WriteString(fmt.Sprintf("%sexport default function %s%s%s { // line %d\n", indent, name, paramText, returnText, lineNum))
						} else {
							result.WriteString(fmt.Sprintf("%sexport function %s%s%s { // line %d\n", indent, name, paramText, returnText, lineNum))
						}
						result.WriteString(fmt.Sprintf("%s  // ...\n", indent))
						result.WriteString(fmt.Sprintf("%s}\n\n", indent))
					}
				
				case "class_declaration":
					nameNode := firstChild.ChildByFieldName("name")
					if nameNode != nil {
						name := getNodeText(nameNode, content)
						
						// Get heritage clause (extends/implements)
						var heritageText string
						for i := 0; i < int(firstChild.ChildCount()); i++ {
							child := firstChild.Child(uint(i))
							if child.Kind() == "class_heritage" {
								heritageText = " " + getNodeText(child, content)
								break
							}
						}
						
						// Get documentation comment if present
						doc := findDocComment(node, content, "typescript")
						if doc != "" {
							docLines := strings.Split(doc, "\n")
							for _, line := range docLines {
								result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
							}
						}
						
						// Write export class declaration
						lineNum := getNodeLineNumber(firstChild)
						if isDefault {
							result.WriteString(fmt.Sprintf("%sexport default class %s%s { // line %d\n", indent, name, heritageText, lineNum))
						} else {
							result.WriteString(fmt.Sprintf("%sexport class %s%s { // line %d\n", indent, name, heritageText, lineNum))
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
				
				case "interface_declaration":
					nameNode := firstChild.ChildByFieldName("name")
					if nameNode != nil {
						name := getNodeText(nameNode, content)
						
						// Get extends clause if any
						extendsNode := firstChild.ChildByFieldName("extends_clause")
						extendsText := ""
						if extendsNode != nil {
							extendsText = " " + getNodeText(extendsNode, content)
						}
						
						// Get documentation comment if present
						doc := findDocComment(node, content, "typescript")
						if doc != "" {
							docLines := strings.Split(doc, "\n")
							for _, line := range docLines {
								result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
							}
						}
						
						// Write export interface declaration
						lineNum := getNodeLineNumber(firstChild)
						result.WriteString(fmt.Sprintf("%sexport interface %s%s { // line %d\n", indent, name, extendsText, lineNum))
						
						// Process interface body for property and method signatures
						bodyNode := firstChild.ChildByFieldName("body")
						if bodyNode != nil {
							for i := 0; i < int(bodyNode.NamedChildCount()); i++ {
								child := bodyNode.NamedChild(uint(i))
								
								if child.Kind() == "property_signature" {
									nameNode := child.ChildByFieldName("name")
									typeNode := child.ChildByFieldName("type")
									
									if nameNode != nil && typeNode != nil {
										propName := getNodeText(nameNode, content)
										propType := getNodeText(typeNode, content)
										
										// Check for optional marker
										optional := ""
										for j := 0; j < int(child.ChildCount()); j++ {
											if child.Child(uint(j)).Kind() == "?" {
												optional = "?"
												break
											}
										}
										
										// Get doc comment
										propDoc := findDocComment(child, content, "typescript")
										if propDoc != "" {
											result.WriteString(fmt.Sprintf("%s  // %s\n", indent, propDoc))
										}
										
										result.WriteString(fmt.Sprintf("%s  %s%s: %s;\n", indent, propName, optional, propType))
									}
								} else if child.Kind() == "method_signature" {
									nameNode := child.ChildByFieldName("name")
									paramNode := child.ChildByFieldName("parameters")
									returnNode := child.ChildByFieldName("return_type")
									
									if nameNode != nil {
										methodName := getNodeText(nameNode, content)
										
										paramText := ""
										if paramNode != nil {
											paramText = getNodeText(paramNode, content)
										}
										
										returnText := ""
										if returnNode != nil {
											returnText = ": " + getNodeText(returnNode, content)
										}
										
										// Get doc comment
										methodDoc := findDocComment(child, content, "typescript")
										if methodDoc != "" {
											result.WriteString(fmt.Sprintf("%s  // %s\n", indent, methodDoc))
										}
										
										result.WriteString(fmt.Sprintf("%s  %s%s%s;\n", indent, methodName, paramText, returnText))
									}
								}
							}
						}
						
						result.WriteString(fmt.Sprintf("%s}\n\n", indent))
					}
				
				case "type_alias_declaration":
					nameNode := firstChild.ChildByFieldName("name")
					typeNode := firstChild.ChildByFieldName("value")
					
					if nameNode != nil && typeNode != nil {
						name := getNodeText(nameNode, content)
						typeValue := getNodeText(typeNode, content)
						
						// Get documentation comment if present
						doc := findDocComment(node, content, "typescript")
						if doc != "" {
							docLines := strings.Split(doc, "\n")
							for _, line := range docLines {
								result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
							}
						}
						
						lineNum := getNodeLineNumber(firstChild)
						result.WriteString(fmt.Sprintf("%sexport type %s = %s; // line %d\n\n", indent, name, typeValue, lineNum))
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
									
									// Get return type if any
									returnNode := valueNode.ChildByFieldName("return_type")
									returnText := ""
									if returnNode != nil {
										returnText = getNodeText(returnNode, content)
									}
									
									// Get documentation comment if present
									doc := findDocComment(node, content, "typescript")
									if doc != "" {
										docLines := strings.Split(doc, "\n")
										for _, line := range docLines {
											result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
										}
									}
									
									// Write export function
									lineNum := getNodeLineNumber(firstChild)
									if valueNode.Kind() == "arrow_function" {
										result.WriteString(fmt.Sprintf("%sexport %s %s = %s%s => { // line %d\n", indent, declType, name, paramText, returnText, lineNum))
									} else {
										result.WriteString(fmt.Sprintf("%sexport %s %s = function%s%s { // line %d\n", indent, declType, name, paramText, returnText, lineNum))
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

							// Get return type if any
							returnNode := valueNode.ChildByFieldName("return_type")
							returnText := ""
							if returnNode != nil {
								returnText = getNodeText(returnNode, content)
							}

							// Get documentation comment if present
							doc := findDocComment(node, content, "typescript")
							if doc != "" {
								docLines := strings.Split(doc, "\n")
								for _, line := range docLines {
									result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
								}
							}

							// Write function
							lineNum := getNodeLineNumber(node)
							if valueNode.Kind() == "arrow_function" {
								result.WriteString(fmt.Sprintf("%s%s %s = %s%s => { // line %d\n", indent, declType, name, paramText, returnText, lineNum))
							} else {
								result.WriteString(fmt.Sprintf("%s%s %s = function%s%s { // line %d\n", indent, declType, name, paramText, returnText, lineNum))
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

		case "type_alias_declaration":
			// For TypeScript type aliases
			nameNode := node.ChildByFieldName("name")
			typeNode := node.ChildByFieldName("value")

			if nameNode != nil && typeNode != nil {
				name := getNodeText(nameNode, content)
				typeValue := getNodeText(typeNode, content)

				// Get documentation comment if present
				doc := findDocComment(node, content, "typescript")
				if doc != "" {
					docLines := strings.Split(doc, "\n")
					for _, line := range docLines {
						result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
					}
				}

				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%stype %s = %s; // line %d\n\n", indent, name, typeValue, lineNum))
			}
		}
	}

	processNode(root, 0)
	return result.String()
}
