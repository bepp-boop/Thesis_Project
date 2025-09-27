package main

import (
	"fmt"
	"strings"

	tree_sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_go "github.com/tree-sitter/tree-sitter-go/bindings/go"
)

func main() {
	language := tree_sitter.NewLanguage(tree_sitter_go.Language())
	parser := tree_sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	sourceCode := []byte(`
			package main
	
			import "fmt"
	
			func main() { fmt.Println("Hello, world!") }
		`)

	tree := parser.Parse(sourceCode, nil)
	defer tree.Close()

	query, err := tree_sitter.NewQuery(
		language,
		`
			(function_declaration
				name: (identifier) @function.name
				body: (block) @function.block
			)
			`,
	)
	if err != nil {
		panic(err)
	}
	defer query.Close()

	qc := tree_sitter.NewQueryCursor()
	defer qc.Close()

	captures := qc.Captures(query, tree.RootNode(), sourceCode)

	for match, index := captures.Next(); match != nil; match, index = captures.Next() {
		functionName := match.Captures[0].Node.Utf8Text(sourceCode)
		functionBody := match.Captures[1].Node.Utf8Text(sourceCode)

		fmt.Printf("Capture %d: Function Name: %s, Body: %s\n", index, functionName, functionBody)
		// Change the capture to new
		if functionName == "main" {
			updatedSource := replaceFunctionBody(sourceCode, functionBody, "fmt.Println(\"Goodbye, World\")")
			tree = parser.Parse(updatedSource, nil)
			defer tree.Close()

			fmt.Println("\nUpdated Source Code:\n", string(updatedSource))

			// Capture the updated function body
			captures = qc.Captures(query, tree.RootNode(), updatedSource)
			fmt.Println("\nUpdated Captures:")
			for match, index := captures.Next(); match != nil; match, index = captures.Next() {
				functionName := match.Captures[0].Node.Utf8Text(updatedSource)
				functionBody := match.Captures[1].Node.Utf8Text(updatedSource)
				fmt.Printf("Capture %d: Function Name: %s, Body: %s\n", index, functionName, functionBody)
			}
		}
	}
}

func replaceFunctionBody(source []byte, oldBody string, newBody string) []byte {
	// Convert to string for easier manipulation
	sourceStr := string(source)

	// Replace the old function body with the new one
	updatedSource := strings.Replace(sourceStr, oldBody, newBody, 1)

	// Return the updated source as a byte slice
	return []byte(updatedSource)
}
