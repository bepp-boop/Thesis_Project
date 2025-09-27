	for match, _ := captures.Next(); match != nil; match, _ = captures.Next() {
		for i, capture := range match.Captures {
			node := capture.Node
			text := node.Utf8Text(sourceCode)
			fmt.Printf("Capture %d: %s\n", i, text)
		}
		fmt.Println("---")
	}
