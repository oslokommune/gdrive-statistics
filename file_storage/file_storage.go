package file_storage

import (
	"fmt"
	"os"
	userPkg "os/user"
	"path"
)

type FileStorage struct{}

func New() *FileStorage {
	return &FileStorage{}
}

const storeFolder = ".gdrive-statistics"

func (_ *FileStorage) CreateStoreFolderIfNotExists() error {
	if _, err := os.Stat(storeFolder); os.IsNotExist(err) {
		return os.Mkdir(storeFolder, 0o744)
	}

	return nil
}

func (_ *FileStorage) GetFilepath(filename string) (string, error) {
	user, err := userPkg.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get user: %w", err)
	}

	userHomeDir := user.HomeDir
	return path.Join(userHomeDir, storeFolder, filename), nil
}

func (s *FileStorage) Save(filename string, content string) error {
	filepath, err := s.GetFilepath(filename)
	if err != nil {
		return fmt.Errorf("could not get file path: %w", err)
	}

	file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o744)
	if err != nil {
		return fmt.Errorf("could not create file: %w", err)
	}

	defer func() {
		err = file.Close()
	}()

	fmt.Printf("Writing to path: %s\n", filepath)

	_, err = file.WriteString(content)
	if err != nil {
		return fmt.Errorf("could not write to file: %w", err)
	}

	return nil
}
