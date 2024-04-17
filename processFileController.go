package main

import (
	"fmt"
	"text/template"
)

func processFileController(path string, directory string, openAPI OpenAPI) {
	tmpl, err := template.ParseFiles("templateController.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := getYamlName(path)

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			generateFile(tmpl, yamlName, operation.OperationId, "Controller", path, directory)
		}
	}
}
