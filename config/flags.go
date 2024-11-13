package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
)

var (
	PortNumber string
	Directory  string
)

func init() {
	flag.StringVar(&PortNumber, "port", "8080", "Port number")
	flag.StringVar(&Directory, "dir", "data", "Path to the directory")

	helpMessage := `Coffee Shop Management System

Usage:
  hot-coffee [--port <N>] [--dir <S>] 
  hot-coffee --help

Options:
  --help       Show this screen.
  --port N     Port number.
  --dir S      Path to the data directory.`

	flag.Usage = func() {
		fmt.Println(helpMessage)
	}
	flag.Parse()
}

func ValidateDirectory() error {
	// checking that --dir=path exists
	if _, err := os.Stat(Directory); os.IsNotExist(err) {
		err = os.Mkdir(Directory, 0o755)
		if err != nil {
			return errors.New("Error: " + err.Error())
		}
		// Automatically create the necessary JSON files
		if err := createEmptyJSONFiles(Directory); err != nil {
			return err
		}
	}
	// checking that '--dir=' is standard or not
	if isStandardPackage(Directory) {
		return errors.New("Error: directory(--dir=) cannot be one of the used ones {'cmd', 'config', 'internal', 'models'}.")
	}
	return nil
}

func createEmptyJSONFiles(directory string) error {
	// Define the list of files to create
	files := []string{"orders.json", "menu.json", "inventory.json"}

	for _, file := range files {
		filePath := fmt.Sprintf("%s/%s", directory, file)
		// Create and initialize each file with an empty array
		if err := createFileWithEmptyArray(filePath); err != nil {
			return err
		}
	}
	return nil
}

func createFileWithEmptyArray(filePath string) error {
	// Open the file for writing (create it if it doesn't exist)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filePath, err)
	}
	defer file.Close()

	// Write an empty array into the file
	emptyArray := []interface{}{}
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // Optional: format the JSON output nicely
	if err := encoder.Encode(emptyArray); err != nil {
		return fmt.Errorf("failed to write empty array to file %s: %v", filePath, err)
	}

	return nil
}

func isStandardPackage(packageName string) bool {
	return packageName == "cmd" || packageName == "config" || packageName == "internal" || packageName == "models"
}
