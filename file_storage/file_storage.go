package file_storage

import (
	"errors"
	"fmt"
	"io/ioutil"
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

// GetFilepath returns the application's data joined with the given filename
func (_ *FileStorage) GetFilepath(filename string) (string, error) {
	user, err := userPkg.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get user: %w", err)
	}

	userHomeDir := user.HomeDir
	return path.Join(userHomeDir, storeFolder, filename), nil
}

func (s *FileStorage) Save(filename string, content []byte) error {
	filepath, err := s.GetFilepath(filename)
	if err != nil {
		return fmt.Errorf("get file path: %w", err)
	}

	return ioutil.WriteFile(filepath, content, 0o744)
}

func (s *FileStorage) AppFileExists(filename string) (bool, error) {
	filePath, err := s.GetFilepath(filename)
	if err != nil {
		return false, fmt.Errorf("get file path: %w", err)
	}

	_, err = os.Stat(filePath)
	if err == nil {
		return true, nil
	}

	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}

	return false, err
}

func (s *FileStorage) Load(filename string) ([]byte, error) {
	filePath, err := s.GetFilepath(filename)
	if err != nil {
		return nil, fmt.Errorf("get file path: %w", err)
	}

	return ioutil.ReadFile(filePath)
}
