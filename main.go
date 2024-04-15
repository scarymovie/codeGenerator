package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

type OpenAPI struct {
	Paths map[string]map[string]struct {
		OperationId string `yaml:"operationId"`
		Parameters  []struct {
			Name        string `yaml:"name"`
			Description string `yaml:"description"`
			In          string `yaml:"in"`
			Schema      struct {
				Type string `yaml:"type"`
			} `yaml:"schema"`
			Required bool `yaml:"required"`
		} `yaml:"parameters"`
		Responses map[string]struct {
			Description string `yaml:"description"`
			Content     map[string]struct {
				Schema struct {
					Type  string `yaml:"type"`
					Items struct {
						Ref string `yaml:"$ref"`
					} `yaml:"items"`
					Ref string `yaml:"$ref"`
				} `yaml:"schema"`
			} `yaml:"content"`
		} `yaml:"responses"`
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
	processFileDefault(path, openAPI)
	processSchemas(path, openAPI)
	processResponses(path, openAPI)
}
