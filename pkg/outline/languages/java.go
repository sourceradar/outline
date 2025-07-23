package languages

import (
	"fmt"
	"strings"

	"github.com/tree-sitter/go-tree-sitter"
)

func processJavaNode(node *tree_sitter.Node, indentLevel int, content []byte, result *strings.Builder) {
	indent := strings.Repeat("\t", indentLevel)

	switch node.Kind() {
	case "program":
		var i uint
		for i = 0; i < node.NamedChildCount(); i++ {
			child := node.NamedChild(i)
			processJavaNode(child, indentLevel, content, result)
		}

	case "package_declaration":
		processJavaPackage(node, content, result, indent)

	case "import_declaration":
		processJavaImport(node, content, result, indent)

	case "class_declaration":
		processJavaClass(node, content, result, indent, indentLevel)

	case "interface_declaration":
		processJavaInterface(node, content, result, indent, indentLevel)

	case "enum_declaration":
		processJavaEnum(node, content, result, indent, indentLevel)

	case "method_declaration":
		processJavaMethod(node, content, result, indent)

	case "constructor_declaration":
		processJavaConstructor(node, content, result, indent)

	case "field_declaration":
		processJavaField(node, content, result, indent)
	}
}

func processJavaPackage(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	packageText := getNodeText(node, content)
	result.WriteString(fmt.Sprintf("%s%s\n\n", indent, packageText))
}

func processJavaImport(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	importText := getNodeText(node, content)
	result.WriteString(fmt.Sprintf("%s%s\n", indent, importText))
}

func processJavaClass(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string, indentLevel int) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := getNodeText(nameNode, content)

	// Get modifiers
	modifiers := getJavaModifiers(node, content)
	modifierText := ""
	if len(modifiers) > 0 {
		modifierText = strings.Join(modifiers, " ") + " "
	}

	// Get superclass
	superclassNode := node.ChildByFieldName("superclass")
	superclassText := ""
	if superclassNode != nil {
		superclassText = " " + getNodeText(superclassNode, content)
	}

	// Get interfaces
	interfacesNode := node.ChildByFieldName("interfaces")
	interfacesText := ""
	if interfacesNode != nil {
		interfacesText = " " + getNodeText(interfacesNode, content)
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "java")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
		}
	}

	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%sclass %s%s%s { // line %d\n", indent, modifierText, name, superclassText, interfacesText, lineNum))

	// Process class body
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		for i := uint(0); i < bodyNode.NamedChildCount(); i++ {
			child := bodyNode.NamedChild(i)
			processJavaNode(child, indentLevel+1, content, result)
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

func processJavaInterface(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string, indentLevel int) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := getNodeText(nameNode, content)

	// Get modifiers
	modifiers := getJavaModifiers(node, content)
	modifierText := ""
	if len(modifiers) > 0 {
		modifierText = strings.Join(modifiers, " ") + " "
	}

	// Get extends clause
	extendsText := ""
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child.Kind() == "extends_interfaces" {
			extendsText = " " + getNodeText(child, content)
			break
		}
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "java")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
		}
	}

	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%sinterface %s%s { // line %d\n", indent, modifierText, name, extendsText, lineNum))

	// Process interface body
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		for i := uint(0); i < bodyNode.NamedChildCount(); i++ {
			child := bodyNode.NamedChild(i)
			processJavaNode(child, indentLevel+1, content, result)
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

func processJavaEnum(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string, indentLevel int) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := getNodeText(nameNode, content)

	// Get modifiers
	modifiers := getJavaModifiers(node, content)
	modifierText := ""
	if len(modifiers) > 0 {
		modifierText = strings.Join(modifiers, " ") + " "
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "java")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
		}
	}

	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%senum %s { // line %d\n", indent, modifierText, name, lineNum))

	// Process enum body - constants and methods
	bodyNode := node.ChildByFieldName("body")
	if bodyNode != nil {
		for i := uint(0); i < bodyNode.NamedChildCount(); i++ {
			child := bodyNode.NamedChild(i)
			if child.Kind() == "enum_constant" {
				constantName := getNodeText(child, content)
				result.WriteString(fmt.Sprintf("%s\t%s,\n", indent, constantName))
			} else if child.Kind() == "enum_body_declarations" {
				// Process methods and other declarations inside the enum
				for j := uint(0); j < child.NamedChildCount(); j++ {
					subchild := child.NamedChild(j)
					processJavaNode(subchild, indentLevel+1, content, result)
				}
			} else {
				processJavaNode(child, indentLevel+1, content, result)
			}
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n\n", indent))
}

func processJavaMethod(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := getNodeText(nameNode, content)

	// Get modifiers
	modifiers := getJavaModifiers(node, content)
	modifierText := ""
	if len(modifiers) > 0 {
		modifierText = strings.Join(modifiers, " ") + " "
	}

	// Get return type
	typeNode := node.ChildByFieldName("type")
	typeText := "void"
	if typeNode != nil {
		typeText = getNodeText(typeNode, content)
	}

	// Get parameters
	parametersNode := node.ChildByFieldName("parameters")
	parametersText := "()"
	if parametersNode != nil {
		parametersText = getNodeText(parametersNode, content)
	}

	// Get throws clause
	throwsText := ""
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child.Kind() == "throws" {
			throwsText = " " + getNodeText(child, content)
			break
		}
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "java")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
		}
	}

	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%s%s %s%s%s { //... } // line %d\n\n", indent, modifierText, typeText, name, parametersText, throwsText, lineNum))
}

func processJavaConstructor(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	nameNode := node.ChildByFieldName("name")
	if nameNode == nil {
		return
	}

	name := getNodeText(nameNode, content)

	// Get modifiers
	modifiers := getJavaModifiers(node, content)
	modifierText := ""
	if len(modifiers) > 0 {
		modifierText = strings.Join(modifiers, " ") + " "
	}

	// Get parameters
	parametersNode := node.ChildByFieldName("parameters")
	parametersText := "()"
	if parametersNode != nil {
		parametersText = getNodeText(parametersNode, content)
	}

	// Get throws clause
	throwsText := ""
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child.Kind() == "throws" {
			throwsText = " " + getNodeText(child, content)
			break
		}
	}

	// Get documentation comment if present
	doc := findDocComment(node, content, "java")
	if doc != "" {
		docLines := strings.Split(doc, "\n")
		for _, line := range docLines {
			result.WriteString(fmt.Sprintf("%s// %s\n", indent, strings.TrimSpace(line)))
		}
	}

	lineNum := getNodeLineNumber(node)
	result.WriteString(fmt.Sprintf("%s%s%s%s%s { //... } // line %d\n\n", indent, modifierText, name, parametersText, throwsText, lineNum))
}

func processJavaField(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	typeNode := node.ChildByFieldName("type")
	if typeNode == nil {
		return
	}

	typeText := getNodeText(typeNode, content)

	// Get modifiers
	modifiers := getJavaModifiers(node, content)
	modifierText := ""
	if len(modifiers) > 0 {
		modifierText = strings.Join(modifiers, " ") + " "
	}

	// Get all variable declarators
	for i := uint(0); i < node.NamedChildCount(); i++ {
		child := node.NamedChild(i)
		if child.Kind() == "variable_declarator" {
			nameNode := child.ChildByFieldName("name")
			if nameNode != nil {
				name := getNodeText(nameNode, content)
				
				// Get initializer if present
				valueNode := child.ChildByFieldName("value")
				valueText := ""
				if valueNode != nil {
					valueText = " = " + getNodeText(valueNode, content)
				}

				lineNum := getNodeLineNumber(node)
				result.WriteString(fmt.Sprintf("%s%s%s %s%s; // line %d\n", indent, modifierText, typeText, name, valueText, lineNum))
			}
		}
	}
}

func getJavaModifiers(node *tree_sitter.Node, content []byte) []string {
	var modifiers []string
	
	for i := uint(0); i < node.ChildCount(); i++ {
		child := node.Child(i)
		if child.Kind() == "modifiers" {
			for j := uint(0); j < child.ChildCount(); j++ {
				modifier := child.Child(j)
				modifierText := getNodeText(modifier, content)
				if modifierText != "" && modifierText != " " {
					modifiers = append(modifiers, modifierText)
				}
			}
			break
		}
	}
	
	return modifiers
}

// ExtractJavaOutline extracts Java outline directly from the code
func ExtractJavaOutline(root *tree_sitter.Node, content []byte) string {
	var result = new(strings.Builder)

	processJavaNode(root, 0, content, result)

	return result.String()
}