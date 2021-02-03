package file_saver

import (
	"fmt"
	"os"
	"path"
)

func Save(content string, filename string) error {
	filepath, err := getFilepath(filename)
	if err != nil {
		return fmt.Errorf("could not get file path: %w", err)
	}

	dirname := "result"
	if _, err := os.Stat(dirname); os.IsNotExist(err) {
		err = os.Mkdir(dirname, 0744)
		if err != nil {
			return fmt.Errorf("could not create directory: %w", err)
		}
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0744)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	defer file.Close()

	fmt.Printf("Writing to path: %s", filepath)

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}

func getFilepath(filename string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("could not get working directory: %w", err)
	}

	return path.Join(currentDir, "result", filename), nil
}
