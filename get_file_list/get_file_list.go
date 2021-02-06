package get_file_list

import (
	"fmt"
	"net/http"

	"github.com/oslokommune/gdrive-statistics/memory_usage"
	"google.golang.org/api/drive/v3"
)

type FileListGetter struct {
	client   *http.Client
	gdriveId string
}

func New(client *http.Client, gdriveId string) *FileListGetter {
	return &FileListGetter{
		client:   client,
		gdriveId: gdriveId,
	}
}

// GetFiles fetches files and folders from the Google Drive API
func (g *FileListGetter) GetFiles(pageCount int) ([]*DriveFile, error) {
	//goland:noinspection ALL
	srv, err := drive.New(g.client)
	if err != nil {
		return nil, fmt.Errorf("Unable to retrieve Drive client: %v", err)
	}

	allFiles := make([]*DriveFile, 0)
	nextPageToken := ""
	i := 0

	for ok := true; ok; ok = nextPageToken != "" && i < pageCount {
		i++

		if len(nextPageToken) > 0 {
			memory_usage.PrintMemUsage()
			fmt.Printf("Fetching page %d: %s\n", i, nextPageToken[len(nextPageToken)-10:])
		}

		fileList, err := srv.Files.List().
			Corpora("drive").                // Comment if using private drive
			DriveId(g.gdriveId).             // Comment if using private drive
			IncludeItemsFromAllDrives(true). // Comment if using private drive
			SupportsAllDrives(true).         // Comment if using private drive
			// IncludeItemsFromAllDrives(false). // Remove this if getting from shared drive
			PageToken(nextPageToken).
			Fields("files(id,name,parents),nextPageToken").
			PageSize(1000).
			OrderBy("folder,modifiedTime").
			Do()
		if err != nil {
			return nil, fmt.Errorf("could not list gdrive files: %w", err)
		}

		nextPageToken = fileList.NextPageToken

		files, err := g.toDriveFile(fileList.Files)
		if err != nil {
			return nil, fmt.Errorf("could not get files: %w", err)
		}

		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}

func (g *FileListGetter) toDriveFile(files []*drive.File) ([]*DriveFile, error) {
	driveFiles := make([]*DriveFile, 0)

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

func (g *FileListGetter) createDriveFile(file *drive.File) (*DriveFile, error) {
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
