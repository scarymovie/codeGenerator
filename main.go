package main

import (
	"fmt"
	"io/ioutil"
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
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".yaml") {
			fmt.Println("Processing file:", path)
			processFileController(path)
			processFileAction(path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %v: %v\n", srcDir, err)
	}
}

func processFileController(path string) {
	yamlFile, err := ioutil.ReadFile(path)
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

	tmpl, err := template.ParseFiles("templateController.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := strings.Replace(path, "src/", "", -1)
	yamlName = strings.Replace(yamlName, ".yaml", "", +1)

	for _, methods := range openAPI.Paths {
		for _, operation := range methods {
			controllerFileName := fmt.Sprintf("%sController.php", strings.Title(operation.OperationId))

			file, err := os.Create(controllerFileName)
			if err != nil {
				fmt.Printf("Error creating file: %s\n", err)
				return
			}
			defer file.Close()

			err = tmpl.Execute(file, map[string]string{
				"Module":    strings.Title(yamlName),
				"Operation": strings.Title(operation.OperationId),
				"actionVar": strings.ToLower(string(operation.OperationId[0])) + operation.OperationId[1:],
			})

			if err != nil {
				fmt.Printf("Error executing template: %s\n", err)
				return
			}

			fmt.Printf("Generated file: %s\n", controllerFileName)
		}
	}
}

func processFileAction(path string) {
	yamlFile, err := ioutil.ReadFile(path)
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

	tmpl, err := template.ParseFiles("templateAction.txt")
	if err != nil {
		fmt.Printf("Error loading template: %s\n", err)
		return
	}

	yamlName := strings.Replace(path, "src/", "", -1)
	yamlName = strings.Replace(yamlName, ".yaml", "", +1)

	for _, methods := range openAPI.Paths {

		for _, operation := range methods {
			actionFileName := fmt.Sprintf("%sAction.php", strings.Title(operation.OperationId))
			file, err := os.Create(actionFileName)
			if err != nil {
				fmt.Printf("Error creating file: %s\n", err)
				return
			}
			defer file.Close()

			err = tmpl.Execute(file, map[string]string{
				"Module":    strings.Title(yamlName),
				"Operation": strings.Title(operation.OperationId),
			})

			if err != nil {
				fmt.Printf("Error executing template: %s\n", err)
				return
			}

			fmt.Printf("Generated file: %s\n", actionFileName)
		}
	}
}
