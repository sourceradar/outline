package languages

import (
	"strings"

	sitter "github.com/tree-sitter/go-tree-sitter"
)

// getNodeText extracts the text of a node from the source content
func getNodeText(node *sitter.Node, content []byte) string {
	return string(content[node.StartByte():node.EndByte()])
}

// getNodeLineNumber returns the line number (1-indexed) of a node's start position
func getNodeLineNumber(node *sitter.Node) uint {
	return node.StartPosition().Row + 1
}

// findDocComment finds and aggregates documentation comments preceding a node
func findDocComment(node *sitter.Node, content []byte, language string) string {
	if node.Parent() == nil {
		return ""
	}

	var comment string
	currentNode := node.PrevNamedSibling()

	for currentNode != nil {
		nodeType := currentNode.Kind()

		if strings.Contains(nodeType, "comment") {
			text := getNodeText(currentNode, content)
			text = strings.TrimSpace(text)
			if comment == "" {
				comment = text
			} else {
				comment = text + "\n" + comment
			}

			currentNode = currentNode.PrevNamedSibling()
		} else {
			break
		}
	}

	return comment
}