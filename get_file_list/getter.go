package get_file_list

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
	"net/http"
)

func GetFiles(client *http.Client, gdriveId string) ([]*DriveFile, error) {
	//goland:noinspection GoDeprecation
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// Remove comments if using shared gdrive
	fileList, err := srv.Files.List().
		//Corpora("drive").
		//DriveId(gdriveId).
		//IncludeItemsFromAllDrives(true).
		//SupportsAllDrives(true).
		IncludeItemsFromAllDrives(false). // Remove this if getting from shared drive
		Fields("files(id,name,parents)").
		PageSize(5).
		OrderBy("folder,modifiedTime").
		Do()
	if err != nil {
		return nil, fmt.Errorf("could not list gdrive files: %w", err)
	}

	//data, err := getSomeSmartDataStructure(fileList)
	//if err != nil {
	//	return nil, fmt.Errorf("error getting views: %w", err)
	//}

	files, err := getFiles(fileList.Files)
	if err != nil {
		return nil, fmt.Errorf("could not get files: %w", err)
	}

	return files, nil
}

func getFiles(files []*drive.File) ([]*DriveFile, error) {
	driveFiles := make([]*DriveFile, 0)

	for _, file := range files {
		if file.Shared == false {
			driveFile, err := createDriveFile(file)
			if err != nil {
				return nil, fmt.Errorf("could not create drive file: %w", err)
			}

			driveFiles = append(driveFiles, driveFile)
		}
	}

	return driveFiles, nil
}

func createDriveFile(file *drive.File) (*DriveFile, error) {
	if len(file.Parents) > 1 {
		return nil, fmt.Errorf("multiple parents not supported: %w")
	}

	parent := ""
	if len(file.Parents) == 1 {
		parent = file.Parents[0]
	}

	df := &DriveFile{
		Id:     file.Id,
		Name:   file.Name,
		Parent: parent,
	}

	return df, nil
}
