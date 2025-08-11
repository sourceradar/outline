package languages

import (
	"fmt"
	"strings"

	"github.com/tree-sitter/go-tree-sitter"
)

func processSwiftNode(node *tree_sitter.Node, indentLevel int, content []byte, result *strings.Builder) {
	if node == nil {
		return
	}

	indent := strings.Repeat("  ", indentLevel)
	nodeType := node.Kind()

	switch nodeType {
	case "import_declaration":
		processSwiftImport(node, content, result, indent)
	case "class_declaration":
		processSwiftClass(node, content, result, indent)
	case "struct_declaration":
		processSwiftStruct(node, content, result, indent)
	case "protocol_declaration":
		processSwiftProtocol(node, content, result, indent)
	case "enum_declaration":
		processSwiftEnum(node, content, result, indent)
	case "function_declaration":
		processSwiftFunction(node, content, result, indent)
	case "init_declaration":
		processSwiftInit(node, content, result, indent)
	case "deinit_declaration":
		processSwiftDeinit(node, content, result, indent)
	case "variable_declaration", "property_declaration":
		processSwiftProperty(node, content, result, indent)
	case "subscript_declaration":
		processSwiftSubscript(node, content, result, indent)
	case "extension_declaration":
		processSwiftExtension(node, content, result, indent)
	case "typealias_declaration":
		processSwiftTypealias(node, content, result, indent)
	}

	// Only process top-level nodes, not all children recursively
	// This prevents duplicate processing of nodes already handled in specific processors
}

func processSwiftImport(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	text := getNodeText(node, content)
	result.WriteString(fmt.Sprintf("%s%s\n", indent, text))
}

func processSwiftClass(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var inheritance []string
	var modifiers []string

	// Check if this is actually a struct, enum, or extension by looking at the source text
	nodeText := getNodeText(node, content)
	isStruct := strings.Contains(nodeText, "struct ")
	isEnum := strings.Contains(nodeText, "enum ")
	isExtension := strings.Contains(nodeText, "extension ")

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "type_identifier":
			if name == "" {
				name = getNodeText(child, content)
			}
		case "user_type":
			// For extensions, the extended type is a user_type child
			if isExtension && name == "" {
				for j := 0; j < int(child.NamedChildCount()); j++ {
					typeChild := child.NamedChild(uint(j))
					if typeChild.Kind() == "type_identifier" {
						name = getNodeText(typeChild, content)
						break
					}
				}
			}
		case "inheritance_specifier":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				inheritChild := child.NamedChild(uint(j))
				if inheritChild.Kind() == "user_type" {
					for k := 0; k < int(inheritChild.NamedChildCount()); k++ {
						typeChild := inheritChild.NamedChild(uint(k))
						if typeChild.Kind() == "type_identifier" {
							inheritance = append(inheritance, getNodeText(typeChild, content))
						}
					}
				}
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	declType := "class"
	if isStruct {
		declType = "struct"
	} else if isEnum {
		declType = "enum"
	} else if isExtension {
		declType = "extension"
	}

	classDecl := declType + " " + name
	if len(modifiers) > 0 {
		classDecl = strings.Join(modifiers, " ") + " " + classDecl
	}
	if len(inheritance) > 0 {
		classDecl += ": " + strings.Join(inheritance, ", ")
	}

	result.WriteString(fmt.Sprintf("%s%s {\n", indent, classDecl))

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()
		if childType == "class_body" || childType == "struct_body" {
			processSwiftClassBody(child, content, result, indent+"  ")
		} else if childType == "enum_class_body" {
			processSwiftEnumClassBody(child, content, result, indent+"  ")
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n", indent))
}

func processSwiftStruct(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var protocols []string
	var modifiers []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "type_identifier":
			if name == "" {
				name = getNodeText(child, content)
			}
		case "inheritance_specifier":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				inheritChild := child.NamedChild(uint(j))
				if inheritChild.Kind() == "type_identifier" {
					protocols = append(protocols, getNodeText(inheritChild, content))
				}
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	structDecl := "struct " + name
	if len(modifiers) > 0 {
		structDecl = strings.Join(modifiers, " ") + " " + structDecl
	}
	if len(protocols) > 0 {
		structDecl += ": " + strings.Join(protocols, ", ")
	}

	result.WriteString(fmt.Sprintf("%s%s {\n", indent, structDecl))

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		if child.Kind() == "struct_body" {
			processSwiftStructBody(child, content, result, indent+"  ")
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n", indent))
}

func processSwiftProtocol(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var inheritance []string
	var modifiers []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "type_identifier":
			if name == "" {
				name = getNodeText(child, content)
			}
		case "inheritance_specifier":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				inheritChild := child.NamedChild(uint(j))
				if inheritChild.Kind() == "type_identifier" {
					inheritance = append(inheritance, getNodeText(inheritChild, content))
				}
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	protocolDecl := "protocol " + name
	if len(modifiers) > 0 {
		protocolDecl = strings.Join(modifiers, " ") + " " + protocolDecl
	}
	if len(inheritance) > 0 {
		protocolDecl += ": " + strings.Join(inheritance, ", ")
	}

	result.WriteString(fmt.Sprintf("%s%s {\n", indent, protocolDecl))

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		if child.Kind() == "protocol_body" {
			processSwiftProtocolBody(child, content, result, indent+"  ")
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n", indent))
}

func processSwiftEnum(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var rawType string
	var modifiers []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "type_identifier":
			if name == "" {
				name = getNodeText(child, content)
			} else if rawType == "" {
				rawType = getNodeText(child, content)
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	enumDecl := "enum " + name
	if len(modifiers) > 0 {
		enumDecl = strings.Join(modifiers, " ") + " " + enumDecl
	}
	if rawType != "" {
		enumDecl += ": " + rawType
	}

	result.WriteString(fmt.Sprintf("%s%s {\n", indent, enumDecl))

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		if child.Kind() == "enum_body" {
			processSwiftEnumBody(child, content, result, indent+"  ")
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n", indent))
}

func processSwiftFunction(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var params []string
	var returnType string
	var modifiers []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "simple_identifier":
			if name == "" {
				name = getNodeText(child, content)
			}
		case "function_parameter_list":
			params = extractSwiftParameters(child, content)
		case "parameter":
			// Function parameters can be direct children
			param := extractSwiftParameter(child, content)
			if param != "" {
				params = append(params, param)
			}
		case "function_type":
			returnType = getNodeText(child, content)
		case "user_type", "type_identifier":
			// Return type can be a direct user_type or type_identifier
			if returnType == "" {
				returnType = getNodeText(child, content)
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	funcDecl := "func " + name + "(" + strings.Join(params, ", ") + ")"
	if len(modifiers) > 0 {
		funcDecl = strings.Join(modifiers, " ") + " " + funcDecl
	}
	if returnType != "" {
		funcDecl += " -> " + returnType
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, funcDecl))
}

func processSwiftInit(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var params []string
	var modifiers []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "function_parameter_list":
			params = extractSwiftParameters(child, content)
		case "parameter":
			// For init methods, parameters are direct children
			param := extractSwiftParameter(child, content)
			if param != "" {
				params = append(params, param)
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	initDecl := "init(" + strings.Join(params, ", ") + ")"
	if len(modifiers) > 0 {
		initDecl = strings.Join(modifiers, " ") + " " + initDecl
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, initDecl))
}

func processSwiftDeinit(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	result.WriteString(fmt.Sprintf("%sdeinit\n", indent))
}

func processSwiftProperty(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var propType string
	var modifiers []string
	var isComputed bool

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "pattern":
			if child.Kind() == "pattern" {
				for j := 0; j < int(child.NamedChildCount()); j++ {
					patternChild := child.NamedChild(uint(j))
					if patternChild.Kind() == "simple_identifier" {
						name = getNodeText(patternChild, content)
					}
				}
			}
		case "type_annotation":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				typeChild := child.NamedChild(uint(j))
				if typeChild.Kind() == "user_type" {
					for k := 0; k < int(typeChild.NamedChildCount()); k++ {
						userTypeChild := typeChild.NamedChild(uint(k))
						if userTypeChild.Kind() == "type_identifier" {
							propType = getNodeText(userTypeChild, content)
						}
					}
				}
			}
		case "computed_property":
			isComputed = true
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	propDecl := name
	if len(modifiers) > 0 {
		propDecl = strings.Join(modifiers, " ") + " " + propDecl
	}
	if propType != "" {
		propDecl += ": " + propType
	}
	if isComputed {
		propDecl += " { get set }"
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, propDecl))
}

func processSwiftSubscript(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var params []string
	var returnType string
	var modifiers []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "function_parameter_list":
			params = extractSwiftParameters(child, content)
		case "parameter":
			// Subscript parameters are direct children
			param := extractSwiftParameter(child, content)
			if param != "" {
				params = append(params, param)
			}
		case "type_annotation":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				typeChild := child.NamedChild(uint(j))
				if typeChild.Kind() == "type_identifier" {
					returnType = getNodeText(typeChild, content)
				}
			}
		case "user_type", "type_identifier":
			// Return type can be a direct user_type
			if returnType == "" {
				returnType = getNodeText(child, content)
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	subscriptDecl := "subscript(" + strings.Join(params, ", ") + ")"
	if len(modifiers) > 0 {
		subscriptDecl = strings.Join(modifiers, " ") + " " + subscriptDecl
	}
	if returnType != "" {
		subscriptDecl += " -> " + returnType
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, subscriptDecl))
}

func processSwiftExtension(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var protocols []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "type_identifier":
			if name == "" {
				name = getNodeText(child, content)
			}
		case "inheritance_specifier":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				inheritChild := child.NamedChild(uint(j))
				if inheritChild.Kind() == "type_identifier" {
					protocols = append(protocols, getNodeText(inheritChild, content))
				}
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	extensionDecl := "extension " + name
	if len(protocols) > 0 {
		extensionDecl += ": " + strings.Join(protocols, ", ")
	}

	result.WriteString(fmt.Sprintf("%s%s {\n", indent, extensionDecl))

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		if child.Kind() == "extension_body" {
			processSwiftExtensionBody(child, content, result, indent+"  ")
		}
	}

	result.WriteString(fmt.Sprintf("%s}\n", indent))
}

func processSwiftTypealias(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var aliasType string
	var modifiers []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "type_identifier":
			if name == "" {
				name = getNodeText(child, content)
			} else if aliasType == "" {
				aliasType = getNodeText(child, content)
			}
		case "function_type", "user_type", "tuple_type":
			if aliasType == "" {
				aliasType = getNodeText(child, content)
			}
		case "modifiers":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				modChild := child.NamedChild(uint(j))
				modifiers = append(modifiers, getNodeText(modChild, content))
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	typealiasDecl := "typealias " + name
	if len(modifiers) > 0 {
		typealiasDecl = strings.Join(modifiers, " ") + " " + typealiasDecl
	}
	if aliasType != "" {
		typealiasDecl += " = " + aliasType
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, typealiasDecl))
}

func processSwiftClassBody(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		processSwiftNode(child, len(indent)/2, content, result)
	}
}

func processSwiftStructBody(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		processSwiftNode(child, len(indent)/2, content, result)
	}
}

func processSwiftProtocolBody(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "protocol_function_declaration":
			processSwiftProtocolFunction(child, content, result, indent)
		case "protocol_property_declaration":
			processSwiftProtocolProperty(child, content, result, indent)
		default:
			processSwiftNode(child, len(indent)/2, content, result)
		}
	}
}

func processSwiftEnumClassBody(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var enumCases []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		if childType == "enum_entry" {
			caseName := ""
			for j := 0; j < int(child.NamedChildCount()); j++ {
				entryChild := child.NamedChild(uint(j))
				if entryChild.Kind() == "simple_identifier" {
					caseName = getNodeText(entryChild, content)
					break
				}
			}
			if caseName != "" {
				enumCases = append(enumCases, caseName)
			}
		} else {
			// Handle other enum members like functions
			processSwiftNode(child, len(indent)/2, content, result)
		}
	}

	if len(enumCases) > 0 {
		result.WriteString(fmt.Sprintf("%scase %s\n", indent, strings.Join(enumCases, ", ")))
	}
}

func processSwiftEnumBody(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		if childType == "enum_case_declaration" {
			processSwiftEnumCase(child, content, result, indent)
		} else {
			processSwiftNode(child, len(indent)/2, content, result)
		}
	}
}

func processSwiftExtensionBody(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		processSwiftNode(child, len(indent)/2, content, result)
	}
}

func processSwiftEnumCase(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var cases []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		if child.Kind() == "enum_case" {
			caseName := ""
			for j := 0; j < int(child.NamedChildCount()); j++ {
				caseChild := child.NamedChild(uint(j))
				if caseChild.Kind() == "simple_identifier" {
					caseName = getNodeText(caseChild, content)
					break
				}
			}
			if caseName != "" {
				cases = append(cases, caseName)
			}
		}
	}

	if len(cases) > 0 {
		result.WriteString(fmt.Sprintf("%scase %s\n", indent, strings.Join(cases, ", ")))
	}
}

func extractSwiftParameters(node *tree_sitter.Node, content []byte) []string {
	var params []string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		if child.Kind() == "function_parameter" {
			var paramName string
			var paramType string

			for j := 0; j < int(child.NamedChildCount()); j++ {
				paramChild := child.NamedChild(uint(j))
				childType := paramChild.Kind()

				if childType == "simple_identifier" && paramName == "" {
					paramName = getNodeText(paramChild, content)
				} else if childType == "type_annotation" {
					for k := 0; k < int(paramChild.NamedChildCount()); k++ {
						typeChild := paramChild.NamedChild(uint(k))
						if typeChild.Kind() == "type_identifier" {
							paramType = getNodeText(typeChild, content)
						}
					}
				}
			}

			param := paramName
			if paramType != "" {
				param += ": " + paramType
			}
			params = append(params, param)
		}
	}

	return params
}

func processSwiftProtocolFunction(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var params []string
	var returnType string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "simple_identifier":
			if name == "" {
				name = getNodeText(child, content)
			}
		case "parameter":
			param := extractSwiftParameter(child, content)
			if param != "" {
				params = append(params, param)
			}
		case "user_type", "type_identifier":
			if returnType == "" {
				returnType = getNodeText(child, content)
			}
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	funcDecl := "func " + name + "(" + strings.Join(params, ", ") + ")"
	if returnType != "" {
		funcDecl += " -> " + returnType
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, funcDecl))
}

func processSwiftProtocolProperty(node *tree_sitter.Node, content []byte, result *strings.Builder, indent string) {
	var name string
	var propType string
	var requirements string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		switch childType {
		case "pattern":
			// Extract property name from pattern
			for j := 0; j < int(child.NamedChildCount()); j++ {
				patternChild := child.NamedChild(uint(j))
				if patternChild.Kind() == "simple_identifier" {
					name = getNodeText(patternChild, content)
				}
			}
		case "type_annotation":
			for j := 0; j < int(child.NamedChildCount()); j++ {
				typeChild := child.NamedChild(uint(j))
				if typeChild.Kind() == "user_type" {
					for k := 0; k < int(typeChild.NamedChildCount()); k++ {
						userTypeChild := typeChild.NamedChild(uint(k))
						if userTypeChild.Kind() == "type_identifier" {
							propType = getNodeText(userTypeChild, content)
						}
					}
				}
			}
		case "protocol_property_requirements":
			requirements = getNodeText(child, content)
		}
	}

	comment := findDocComment(node, content, "swift")
	if comment != "" {
		result.WriteString(fmt.Sprintf("%s%s\n", indent, comment))
	}

	propDecl := name
	if propType != "" {
		propDecl += ": " + propType
	}
	if requirements != "" {
		propDecl += " " + requirements
	}

	result.WriteString(fmt.Sprintf("%s%s\n", indent, propDecl))
}

func extractSwiftParameter(node *tree_sitter.Node, content []byte) string {
	var paramNames []string
	var paramType string

	for i := 0; i < int(node.NamedChildCount()); i++ {
		child := node.NamedChild(uint(i))
		childType := child.Kind()

		if childType == "simple_identifier" {
			paramNames = append(paramNames, getNodeText(child, content))
		} else if childType == "optional_type" || childType == "user_type" || childType == "type_identifier" {
			paramType = getNodeText(child, content)
		}
	}

	if len(paramNames) > 0 {
		paramStr := strings.Join(paramNames, " ")
		if paramType != "" {
			return paramStr + ": " + paramType
		}
		return paramStr
	}
	return ""
}

// ExtractSwiftOutline extracts Swift outline directly from the code
func ExtractSwiftOutline(root *tree_sitter.Node, content []byte) string {
	var result strings.Builder

	// Only process direct children of the source file
	for i := 0; i < int(root.NamedChildCount()); i++ {
		child := root.NamedChild(uint(i))
		processSwiftNode(child, 0, content, &result)
	}

	return result.String()
}
