package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func processResponses(path string, directory string, openAPI OpenAPI) {
	for _, methods := range openAPI.Paths {
		for _, response := range methods {
			for _, responses := range response.Responses {
				for _, content := range responses.Content {
					schemaName := strings.Split(content.Schema.Ref, "/")[len(strings.Split(content.Schema.Ref, "/"))-1]
					if schemaName == "NwkError" {
						continue
					}

					if content.Schema.Type == "array" {
						tmpl, err := template.ParseFiles("templateResult200array.txt")
						if err != nil {
							fmt.Printf("Error loading template: %s\n", err)
							return
						}
						generateResponseFile(tmpl, path, content.Schema.Items.Ref, directory)
					} else if content.Schema.Ref != "" {
						tmpl, err := template.ParseFiles("templateResult200.txt")
						if err != nil {
							fmt.Printf("Error loading template: %s\n", err)
							return
						}
						generateResponseFile(tmpl, path, content.Schema.Ref, directory)
					}
				}
			}
		}
	}
}

func generateResponseFile(tmpl *template.Template, path string, ref string, directory string) {
	schemaName := strings.Split(ref, "/")[len(strings.Split(ref, "/"))-1]
	fileName := fmt.Sprintf("%s/%s/%sResult200.php", strings.Title(path), strings.Title(directory), strings.Title(schemaName))

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

	moduleName := strings.Title(strings.TrimSuffix(filepath.Base(path), ".yaml"))

	err = tmpl.Execute(file, map[string]interface{}{
		"Module":      moduleName,
		"Operation":   strings.Title(schemaName),
		"NwkDto":      schemaName,
		"NwkDtoArray": schemaName,
	})

	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		return
	}

	fmt.Printf("Generated file: %s\n", fileName)
}
