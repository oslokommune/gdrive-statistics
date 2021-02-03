package get_file_list

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/memory_usage"
	"google.golang.org/api/drive/v3"
	"log"
	"net/http"
)

func GetFiles(client *http.Client, gdriveId string) ([]*DriveFile, error) {
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	allFiles := make([]*DriveFile, 0)
	nextPageToken := ""
	i := 0

	for ok := true; ok; ok = nextPageToken != "" && i < 10000 {
		i++

		if len(nextPageToken) > 0 {
			memory_usage.PrintMemUsage()
			fmt.Printf("Fetching page %d: %s\n", i, nextPageToken[len(nextPageToken)-10:])
		}

		fileList, err := srv.Files.List().
			Corpora("drive").                // Comment if using private drive
			DriveId(gdriveId).               // Comment if using private drive
			IncludeItemsFromAllDrives(true). // Comment if using private drive
			SupportsAllDrives(true).         // Comment if using private drive
			//IncludeItemsFromAllDrives(false). // Remove this if getting from shared drive
			PageToken(nextPageToken).
			Fields("files(id,name,parents),nextPageToken").
			PageSize(1000).
			OrderBy("folder,modifiedTime").
			Do()
		if err != nil {
			return nil, fmt.Errorf("could not list gdrive files: %w", err)
		}

		nextPageToken = fileList.NextPageToken

		files, err := toDriveFile(fileList.Files)
		if err != nil {
			return nil, fmt.Errorf("could not get files: %w", err)
		}

		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}

func toDriveFile(files []*drive.File) ([]*DriveFile, error) {
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
		return nil, fmt.Errorf("multiple parents (%d) not supported", len(file.Parents))
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
