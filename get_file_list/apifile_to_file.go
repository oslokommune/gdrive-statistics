package get_file_list

import (
	"fmt"
	"google.golang.org/api/drive/v3"
)

func (g *FileListGetter) toDriveFile(files []*drive.File) ([]*FileOrFolder, error) {
	driveFiles := make([]*FileOrFolder, 0)

	for _, file := range files {
		if file.Shared == false {
			driveFile, err := g.createDriveFile(file)
			if err != nil {
				return nil, fmt.Errorf("could not create drive file: %w", err)
			}

			driveFiles = append(driveFiles, driveFile)
		}
	}

	return driveFiles, nil
}

func (g *FileListGetter) createDriveFile(file *drive.File) (*FileOrFolder, error) {
	if len(file.Parents) > 1 {
		return nil, fmt.Errorf("multiple parents (%d) not supported", len(file.Parents))
	}

	parent := ""
	if len(file.Parents) == 1 {
		parent = file.Parents[0]
	}

	df := &FileOrFolder{
		Id:     file.Id,
		Name:   file.Name,
		Parent: parent,
	}

	return df, nil
}
