package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
)

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

	// Получаем имя файла без расширения
	fileName := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))

	// Создаем новую директорию с именем файла
	newDirPath := filepath.Join("./", fileName)
	err = os.MkdirAll(newDirPath, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %s\n", err)
		return
	}

	// Обрабатываем файлы, используя новый путь к директории
	processFileController(filepath.Join(newDirPath, ""), openAPI)
	processFileAction(filepath.Join(newDirPath, ""), openAPI)
	processFileDefault(filepath.Join(newDirPath, ""), openAPI)
	processSchemas(filepath.Join(newDirPath, ""), openAPI)
	processResponses(filepath.Join(newDirPath, ""), openAPI)
}
