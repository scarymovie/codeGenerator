package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func processSchemas(path string, directory string, openAPI OpenAPI) {
	tmpl, err := template.ParseFiles("templateNwk.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	for schemaName, schema := range openAPI.Components.Schemas {
		generateSchemaFile(tmpl, path, schemaName, schema, directory)
	}
}

func getYamlName(path string) string {
	yamlName := strings.Replace(path, "src/", "", -1)
	return strings.Replace(yamlName, ".yaml", "", +1)
}

func generateSchemaFile(tmpl *template.Template, path string, schemaName string, schema struct {
	Type        string   `yaml:"type"`
	Description string   `yaml:"description"`
	Required    []string `yaml:"required"`
	Properties  map[string]struct {
		Type        string `yaml:"type"`
		Description string `yaml:"description"`
		Format      string `yaml:"format"`
	} `yaml:"properties"`
}, directory string) {
	fileName := fmt.Sprintf("%s/%s/%s.php", strings.Title(path), strings.Title(directory), strings.Title(schemaName))

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

	properties := make(map[string]map[string]string)
	for propName, prop := range schema.Properties {
		propType := prop.Type
		if propType == "integer" {
			propType = "int"
		}
		if !contains(schema.Required, propName) {
			propType += "?"
		}
		properties[propName] = map[string]string{
			"Type": propType,
			"Name": propName,
		}
	}

	moduleName := strings.Title(strings.TrimSuffix(filepath.Base(path), ".yaml"))

	err = tmpl.Execute(file, map[string]interface{}{
		"Module":     moduleName,
		"Operation":  strings.Title(schemaName),
		"Properties": properties,
	})

	if err != nil {
		fmt.Printf("Error executing template: %s\n", err)
		return
	}

	fmt.Printf("Generated file: %s\n", fileName)
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}
