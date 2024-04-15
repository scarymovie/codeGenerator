package main

import (
	"fmt"
	"text/template"
)

func processFileDefault(path string, openAPI OpenAPI) {
	tmpl, err := template.ParseFiles("templateResultDefault.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := getYamlName(path)

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			generateFile(tmpl, yamlName, operation.OperationId, "ResultDefault", path)
		}
	}
}
