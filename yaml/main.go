package main

import (
	"fmt"
	"os"
	"regexp"

	tree_sitter_yaml "github.com/tree-sitter-grammars/tree-sitter-yaml/bindings/go"
	tree_sitter "github.com/tree-sitter/go-tree-sitter"
)

func main() {
	language := tree_sitter.NewLanguage(tree_sitter_yaml.Language())
	parser := tree_sitter.NewParser()
	defer parser.Close()
	parser.SetLanguage(language)

	sourceCode, err := readFile()
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	tree := parser.Parse(sourceCode, nil)
	defer tree.Close()

	querySource := `
(
  (block_mapping_pair
    key: (flow_node
      (plain_scalar
        (string_scalar) @key))
    value: (block_node) @value)
     (#match? @key "^with$")
)

	`
	query, err := tree_sitter.NewQuery(language, querySource)
	defer query.Close()

	qc := tree_sitter.NewQueryCursor()
	defer qc.Close()

	captures := qc.Captures(query, tree.RootNode(), sourceCode)

	//print out the captures
	for match, _ := captures.Next(); match != nil; match, _ = captures.Next() {
		functionBody := match.Captures[1].Node.Utf8Text(sourceCode)
		scanForInterpolation([]byte(functionBody))
	}

}
func scanForInterpolation(sourceCode []byte) {
	// Convert to string for easier manipulation
	sourceStr := string(sourceCode)

	fmt.Println("Source Code:\n", sourceStr)
	pattern := `\$\{\{ [^\}]+ \}\}`

	// Compile the regular expression
	re := regexp.MustCompile(pattern)

	// Find all occurrences of the pattern
	matches := re.FindAllString(sourceStr, -1)

	// Print the matches (interpolated values)
	if len(matches) > 0 {
		fmt.Println("\nFound Interpolations:")
		for _, match := range matches {
			fmt.Println(match)
		}
	} else {
		fmt.Println("\nNo interpolations found.")
	}
}

func giveEnvSection() string {

func createUpdatedFile(updatedSource []byte) error {
	err := os.WriteFile("updated.yaml", updatedSource, 0644)
	if err != nil {
		return err
	}
	return nil
}

func readFile() ([]byte, error) {
	sourceCode, err := os.ReadFile("ADES101_multi.yaml")
	if err != nil {
		return nil, err
	}
	return sourceCode, nil
}

// for match, index := captures.Next(); match != nil; match, index = captures.Next() {
// 	functionName := match.Captures[0].Node.Utf8Text(sourceCode)
// 	functionBody := match.Captures[1].Node.Utf8Text(sourceCode)

// 	fmt.Printf("Capture %d: Function Name: %s, Body: %s\n", index, functionName, functionBody)
// 	// Change the capture to new
// 	if functionName == "script" {
// 		updatedSource := replaceFunctionBody(sourceCode, functionBody)
// 		tree = parser.Parse(updatedSource, nil)
// 		defer tree.Close()

// 		fmt.Println("\nUpdated Source Code:\n", string(updatedSource))

// 		err := createUpdatedFile(updatedSource)
// 		if err != nil {
// 			fmt.Println("Error creating updated file:", err)
// 			return
// 		}
// 	}
// }

//func replaceADES101(source []byte, oldBody string) []byte {
// 	// Convert to string for easier manipulation
// 	sourceStr := string(source)

// 	// Find the position of the old function body
// 	startIndex := strings.Index(sourceStr, oldBody)
// 	if startIndex == -1 {
// 		return source
// 	}

// 	// Insert the env section before the script line
// 	envSection := `
//   env:
//     NAME: ${{ inputs.name }}`

// 	withSection := `
//   with:
//     script: console.log("Hello \${process.env.NAME}")`
// 	newBodyWithEnv := envSection + withSection

// 	// Replace the old function body with the new one
// 	updatedSource := strings.Replace(sourceStr, oldBody, newBodyWithEnv, 1)

// 	// Return the updated source as a byte slice
// 	return []byte(updatedSource)
// }
