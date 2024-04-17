package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func processFileAction(path string, directory string, openAPI OpenAPI) {
	tmpl, err := template.ParseFiles("templateAction.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := getYamlName(path)

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			parameters := make([]string, len(operation.Parameters))
			for i, parameter := range operation.Parameters {
				parameters[i] = fmt.Sprintf("%s $%s", parameter.Schema.Type, parameter.Name)
			}
			generateActionFile(tmpl, yamlName, operation.OperationId, "Action", parameters, path, directory)
		}
	}
}

func generateActionFile(tmpl *template.Template, yamlName string, operationId string, fileType string, parameters []string, path string, directory string) {
	fileName := fmt.Sprintf("%s/%s/%s%s.php", strings.Title(path), strings.Title(directory), fileType, strings.Title(operationId))

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

	err = tmpl.Execute(file, map[string]interface{}{
		"Module":     strings.Title(yamlName),
		"Operation":  strings.Title(operationId),
		"Parameters": strings.Join(parameters, ", "),
	})

	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		return
	}

	fmt.Printf("Generated file: %s\n", fileName)
}
