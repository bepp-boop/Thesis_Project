package main

import (
	"fmt"
	"strings"
	"os"
	"parser-test/parser"
)
func main() {
	data, err := os.ReadFile("check.yml")
	if err != nil {
		fmt.Println(err)
		return
	}

	workflow,

}
