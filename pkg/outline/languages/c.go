package languages

import (
	"fmt"
	"strings"

	"github.com/tree-sitter/go-tree-sitter"
)

func processCNode(node *tree_sitter.Node, indentLevel int, content []byte, result *strings.Builder) {
	indent := strings.Repeat("\t", indentLevel)

	// Process based on node type
	switch node.Kind() {
	case "translation_unit":
		var i uint
		for i = 0; i < node.NamedChildCount(); i++ {
			child := node.NamedChild(i)
			processCNode(child, indentLevel, content, result)
		}

	case "preproc_def", "preproc_function_def":
		processCDefine(node, content, result, indent)

	case "preproc_include":
		processCInclude(node, content, result, indent)

	case "function_definition":
		processCFunction(node, content, result, indent)

	case "declaration":
		processCDeclaration(node, content, result, indent)

	case "struct_specifier", "union_specifier", "enum_specifier":
		processCStructUnionEnum(node, content, result, indent)

	case "type_definition":
		processCTypedef(node, content, result, indent)

	case "namespace_definition":
		processCNamespace(node, indentLevel, content, result)

	case "class_specifier":
		processCClass(node, content, result, indent)

	case "template_declaration":
		processCTemplateDeclaration(node, indentLevel, content, result)

	default:
		// Handle other node types by checking children
		var i uint
		for i = 0; i < node.NamedChildCount(); i++ {
			child := node.NamedChild(i)
			processCNode(child, indentLevel, content, result)
		}
	}
}

func processCDefine(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	defineText := getNodeText(node, content)
	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%s // line %d\n", indent, defineText, lineNum))
}

func processCInclude(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	includeText := getNodeText(node, content)
	result.WriteString(fmt.Sprintf("%s%s\n", indent, includeText))
}

func processCFunction(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	declaratorNode := node.ChildByFieldName("declarator")
	if declaratorNode == nil {
		return
	}

	// Extract function name from declarator
	functionName := extractFunctionName(declaratorNode, content)
	if functionName == "" {
		return
	}

	// Get full function signature
	signature := extractFunctionSignature(node, content)

	// Get documentation comment if present
	doc := findDocComment(node, content, "c")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
				result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.TrimSpace(line)))
			} else {
				result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
			}
		}
	}

	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%s { //... } // line %d\n\n", indent, signature, lineNum))
}

func processCDeclaration(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	declarationText := getNodeText(node, content)

	// Skip function declarations that are just prototypes
	if strings.Contains(declarationText, ";") && !strings.Contains(declarationText, "{") {
		lineNum := getNodeLineNumber(node)
		result.WriteString(fmt.Sprintf("%s%s // line %d\n", indent, strings.TrimSpace(declarationText), lineNum))
	}
}

func processCStructUnionEnum(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var structType string
	switch node.Kind() {
	case "struct_specifier":
		structType = "struct"
	case "union_specifier":
		structType = "union"
	case "enum_specifier":
		structType = "enum"
	}

	// Get struct/union/enum name
	nameNode := node.ChildByFieldName("name")
	name := ""
	if nameNode != nil {
		name = getNodeText(nameNode, content)
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "c")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
				result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.TrimSpace(line)))
			} else {
				result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
			}
		}
	}

	lineNum := getNodeLineNumber(node)
	if name != "" {
		result.WriteString(fmt.Sprintf("%s%s %s { // line %d\n", indent, structType, name, lineNum))
	} else {
		result.WriteString(fmt.Sprintf("%s%s { // line %d\n", indent, structType, lineNum))
	}

	// Process fields/members
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		processCStructBody(bodyNode, 1, content, result)
	}

	result.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

func processCStructBody(bodyNode *tree_sitter.Node, indentLevel int, content []byte, result *strings.Builder) {
	indent := strings.Repeat("\t", indentLevel)

	for i := uint(0); i < bodyNode.NamedChildCount(); i++ {
		child := bodyNode.NamedChild(i)
		if child.Kind() == "field_declaration" {
			fieldText := getNodeText(child, content)
			result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.TrimSpace(fieldText)))
		} else if child.Kind() == "enumerator" {
			enumText := getNodeText(child, content)
			result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.TrimSpace(enumText)))
		}
	}
}

func processCTypedef(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	typedefText := getNodeText(node, content)
	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%s // line %d\n\n", indent, strings.TrimSpace(typedefText), lineNum))
}

// C++ specific functions
func processCNamespace(node *tree_sitter.Node, indentLevel int, content []byte, result *strings.Builder) {
	indent := strings.Repeat("\t", indentLevel)

	nameNode := node.ChildByFieldName("name")
	name := ""
	if nameNode != nil {
		name = getNodeText(nameNode, content)
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "cpp")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
				result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.TrimSpace(line)))
			} else {
				result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
			}
		}
	}

	lineNum := getNodeLineNumber(node)
	if name != "" {
		result.WriteString(fmt.Sprintf("%snamespace %s { // line %d\n", indent, name, lineNum))
	} else {
		result.WriteString(fmt.Sprintf("%snamespace { // line %d\n", indent, lineNum))
	}

	// Process namespace body
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		for i := uint(0); i < bodyNode.NamedChildCount(); i++ {
			child := bodyNode.NamedChild(i)
			processCNode(child, indentLevel+1, content, result)
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

func processCClass(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	nameNode := node.ChildByFieldName("name")
	name := ""
	if nameNode != nil {
		name = getNodeText(nameNode, content)
	}

	// Get base classes if any
	baseClause := ""
	for i := uint(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(i)
		if child.Kind() == "base_class_clause" {
			baseClause = " : " + getNodeText(child, content)
			break
		}
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "cpp")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
				result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.TrimSpace(line)))
			} else {
				result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
			}
		}
	}

	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%sclass %s%s { // line %d\n", indent, name, baseClause, lineNum))

	// Process class body
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		processCClassBody(bodyNode, strings.Repeat("\t", 1), content, result)
	}

	result.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

func processCClassBody(bodyNode *tree_sitter.Node, indent string, content []byte, result *strings.Builder) {
	currentVisibility := "private" // Default for class

	for i := uint(0); i < bodyNode.NamedChildCount(); i++ {
		child := bodyNode.NamedChild(i)

		switch child.Kind() {
		case "access_specifier":
			visibility := getNodeText(child, content)
			currentVisibility = strings.TrimSuffix(visibility, ":")
			result.WriteString(fmt.Sprintf("%s%s:\n", indent, currentVisibility))

		case "function_definition":
			signature := extractFunctionSignature(child, content)
			lineNum := getNodeLineNumber(child)
			result.WriteString(fmt.Sprintf("%s\t%s { //... } // line %d\n", indent, signature, lineNum))

		case "declaration":
			declText := getNodeText(child, content)
			if strings.Contains(declText, "(") && strings.Contains(declText, ")") {
				// Method declaration
				lineNum := getNodeLineNumber(child)
				result.WriteString(fmt.Sprintf("%s\t%s // line %d\n", indent, strings.TrimSpace(declText), lineNum))
			} else {
				// Field declaration
				lineNum := getNodeLineNumber(child)
				result.WriteString(fmt.Sprintf("%s\t%s // line %d\n", indent, strings.TrimSpace(declText), lineNum))
			}

		case "constructor_declaration", "destructor_declaration":
			signature := extractFunctionSignature(child, content)
			lineNum := getNodeLineNumber(child)
			result.WriteString(fmt.Sprintf("%s\t%s { //... } // line %d\n", indent, signature, lineNum))
		}
	}
}

func extractFunctionName(declaratorNode *tree_sitter.Node, content []byte) string {
	// Handle different declarator types
	switch declaratorNode.Kind() {
	case "function_declarator":
		nameNode := declaratorNode.ChildByFieldName("declarator")
		if nameNode != nil {
			return extractFunctionName(nameNode, content)
		}
	case "identifier":
		return getNodeText(declaratorNode, content)
	case "pointer_declarator":
		nameNode := declaratorNode.ChildByFieldName("declarator")
		if nameNode != nil {
			return extractFunctionName(nameNode, content)
		}
	}
	return ""
}

func extractFunctionSignature(node *tree_sitter.Node, content []byte) string {
	// Try to build a clean function signature
	var parts []string

	// Get return type if present
	typeNode := node.ChildByFieldName("type")
	if typeNode != nil {
		returnType := getNodeText(typeNode, content)
		if !strings.Contains(returnType, "(") { // Avoid including function pointer types
			parts = append(parts, returnType)
		}
	}

	// Get declarator (contains function name and parameters)
	declaratorNode := node.ChildByFieldName("declarator")
	if declaratorNode != nil {
		declaratorText := getNodeText(declaratorNode, content)
		parts = append(parts, declaratorText)
	}

	signature := strings.Join(parts, " ")

	// Clean up the signature
	signature = strings.ReplaceAll(signature, "\n", " ")
	signature = strings.ReplaceAll(signature, "\t", " ")

	// Normalize multiple spaces to single space
	for strings.Contains(signature, "  ") {
		signature = strings.ReplaceAll(signature, "  ", " ")
	}

	return strings.TrimSpace(signature)
}

// ExtractCOutline extracts C outline directly from the code
func ExtractCOutline(root *tree_sitter.Node, content []byte) string {
	var result = new(strings.Builder)

	// Function to process a node and its children
	processCNode(root, 0, content, result)

	return result.String()
}

func processCTemplateDeclaration(node *tree_sitter.Node, indentLevel int, content []byte, result *strings.Builder) {
	indent := strings.Repeat("\t", indentLevel)

	// Get the template declaration text
	templateText := getNodeText(node, content)
	lines := strings.Split(templateText, "\n")

	// Get documentation comment if present
	doc := findDocComment(node, content, "cpp")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			if strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") || strings.HasPrefix(line, "*") {
				result.WriteString(fmt.Sprintf("%s%s\n", indent, strings.TrimSpace(line)))
			} else {
				result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
			}
		}
	}

	lineNum := getNodeLineNumber(node)
	// Write first line (template declaration)
	if len(lines) > 0 {
		result.WriteString(fmt.Sprintf("%s%s // line %d\n", indent, strings.TrimSpace(lines[0]), lineNum))
	}

	// Process the templated declaration (class, function, etc.)
	for i := uint(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(i)
		if child.Kind() != "template_parameter_list" {
			processCNode(child, indentLevel, content, result)
		}
	}
}

// ExtractCppOutline extracts C++ outline directly from the code
func ExtractCppOutline(root *tree_sitter.Node, content []byte) string {
	var result = new(strings.Builder)

	// Function to process a node and its children (same as C, but handles C++ constructs)
	processCNode(root, 0, content, result)

	return result.String()
}
