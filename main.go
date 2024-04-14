package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

type OpenAPI struct {
	Paths map[string]map[string]struct {
		OperationId string `yaml:"operationId"`
	} `yaml:"paths"`
	Components struct {
		Schemas map[string]struct {
			Type        string   `yaml:"type"`
			Description string   `yaml:"description"`
			Required    []string `yaml:"required"`
			Properties  map[string]struct {
				Type        string `yaml:"type"`
				Description string `yaml:"description"`
				Format      string `yaml:"format"`
			} `yaml:"properties"`
		} `yaml:"schemas"`
	} `yaml:"components"`
}

func main() {
	srcDir := "./src"
	err := processDirectory(srcDir)
	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", srcDir, err)
	}
}

func processDirectory(srcDir string) error {
	return filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			fmt.Println("Processing file:", path)
			processFile(path)
		}

		return nil
	})
}

func processFile(path string) {
	yamlFile, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading YAML file: %s\n", err)
		return
	}

	var openAPI OpenAPI

	err = yaml.Unmarshal(yamlFile, &openAPI)
	if err != nil {
		fmt.Printf("Error parsing YAML file: %s\n", err)
		return
	}

	processFileController(path, openAPI)
	processFileAction(path, openAPI)
}

func processFileController(path string, openAPI OpenAPI) {
	tmpl, err := template.ParseFiles("templateController.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := getYamlName(path)

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			generateFile(tmpl, yamlName, operation.OperationId, "Controller")
		}
	}
}

func processFileAction(path string, openAPI OpenAPI) {
	tmpl, err := template.ParseFiles("templateAction.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := getYamlName(path)

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			generateFile(tmpl, yamlName, operation.OperationId, "Action")
		}
	}
}

func getYamlName(path string) string {
	yamlName := strings.Replace(path, "src/", "", -1)
	return strings.Replace(yamlName, ".yaml", "", +1)
}

func generateFile(tmpl *template.Template, yamlName string, operationId string, fileType string) {
	fileName := fmt.Sprintf("%s%s.php", strings.Title(operationId), fileType)

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
