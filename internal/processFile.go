package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

func ProcessFile(path string) {
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

	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	newDirPath := filepath.Join("./", fileName)
	err = os.MkdirAll(strings.Title(newDirPath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}
	controllerDirPath := filepath.Join(newDirPath + "/Controller")
	err = os.MkdirAll(strings.Title(controllerDirPath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}
	apiDirPath := filepath.Join(newDirPath + "/Api")
	err = os.MkdirAll(strings.Title(apiDirPath), os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}

	actionParams := processFileAction(filepath.Join(newDirPath, ""), "Api", openAPI)
	processFileController(filepath.Join(newDirPath, ""), "Controller", actionParams, openAPI)
	processFileDefault(filepath.Join(newDirPath, ""), "Api", openAPI)
	processSchemas(filepath.Join(newDirPath, ""), "Api", openAPI)
	processResponses(filepath.Join(newDirPath, ""), "Api", openAPI)
}
