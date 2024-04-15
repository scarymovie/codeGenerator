package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func generateFile(tmpl *template.Template, yamlName string, operationId string, fileType string, path string) {
	if operationId == "" {
		fmt.Println("Warning: operationId is empty. Skipping file generation.")
		return
	}
	fileName := fmt.Sprintf("%s/%s%s.php", path, strings.Title(operationId), fileType)

	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("Error creating file: %s\n", err)
		return
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Printf("Error closing file: %s\n", err)
		}
	}(file)

	err = tmpl.Execute(file, map[string]string{
		"Module":    strings.Title(yamlName),
		"Operation": strings.Title(operationId),
		"actionVar": strings.ToLower(string(operationId[0])) + operationId[1:],
	})

	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		return
	}

	fmt.Printf("Generated file: %s\n", fileName)
}
