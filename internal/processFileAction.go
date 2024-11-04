package internal

import (
	"fmt"
	"os"
	"strings"
	"text/template"
)

func processFileAction(path string, directory string, openAPI OpenAPI) []string {
	tmpl, err := template.ParseFiles("templateAction.txt")
	if err != nil {
		println("Error loading template: %s\n", err)
		return []string{}
	}

	yamlName := getYamlName(path)
	var actionParams []string

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			parameters := make([]string, len(operation.Parameters))
			for i, parameter := range operation.Parameters {
				parameters[i] = fmt.Sprintf("%s $%s", parameter.Schema.Type, parameter.Name)
			}
			actionParams = generateActionFile(tmpl, yamlName, operation.OperationId, "Action", parameters, path, directory)
		}
	}
	return actionParams
}

func generateActionFile(tmpl *template.Template, yamlName string, operationId string, fileType string, parameters []string, path string, directory string) []string {
	fileName := fmt.Sprintf("%s/%s/%s%s.php", strings.Title(path), strings.Title(directory), strings.Title(operationId), fileType)
	file, err := os.Create(fileName)
	if err != nil {
		println("Error creating file: %s\n", err)
		return []string{}
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
		println("Error executing template: %s\n", err)
		return []string{}
	}

	fmt.Printf("Generated file: %s\n", fileName)
	return parameters
}
