package main

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func processFileController(path string, directory string, actionParams []string, openAPI OpenAPI) {
	tmpl, err := template.ParseFiles("templateController.go")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := getYamlName(path)

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			generateControllerFile(tmpl, yamlName, actionParams, operation.OperationId, "Controller", path, directory)
		}
	}
}

func generateControllerFile(tmpl *template.Template, yamlName string, actionParams []string, operationId string, fileType string, path string, directory string) {
	if operationId == "" {
		fmt.Println("Warning: operationId is empty. Skipping file generation.")
		return
	}
	fileName := fmt.Sprintf("%s/%s/%s%s.php", strings.Title(path), directory, strings.Title(operationId), fileType)

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

	for i, param := range actionParams {
		words := strings.Fields(param)
		if len(words) > 1 {
			actionParams[i] = strings.TrimLeft(words[1], "$")
		}
	}

	err = tmpl.Execute(file, map[string]string{
		"Module":       strings.Title(yamlName),
		"Operation":    strings.Title(operationId),
		"actionVar":    strings.ToLower(string(operationId[0])) + operationId[1:],
		"actionParams": strings.Join(actionParams, ","),
	})

	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		return
	}

	fmt.Printf("Generated file: %s\n", fileName)
}
