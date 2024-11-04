package main

import (
	"fmt"
	"phpGenerator/internal"
)

func main() {
	srcDir := "./src"
	paths, err := internal.ProcessDirectory(srcDir)
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", srcDir, err)
		return
	}

	for _, path := range paths {
		fmt.Println("Processing file:", path)
		internal.ProcessFile(path)
	}
}
