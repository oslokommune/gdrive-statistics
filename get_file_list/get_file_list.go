package get_file_list

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/file_storage"
	"net/http"

	"github.com/oslokommune/gdrive-statistics/memory_usage"
	"google.golang.org/api/drive/v3"
)

type FileListGetter struct {
	client   *http.Client
	gdriveId string
	storage  *file_storage.FileStorage
}

func New(client *http.Client, gdriveId string, storage *file_storage.FileStorage) *FileListGetter {
	return &FileListGetter{
		client:   client,
		gdriveId: gdriveId,
		storage:  storage,
	}
}

// GetAndStoreFiles fetches files and folders from the Google Drive API
func (g *FileListGetter) GetAndStoreFiles(filename string, pageCount int) ([]*DriveFile, error) {
	//goland:noinspection ALL
	srv, err := drive.New(g.client)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	allFiles := make([]*DriveFile, 0)
	pageToken := ""
	i := 0

	for ok := true; ok; ok = pageToken != "" && i < pageCount {
		i++

		if len(pageToken) > 0 {
			memory_usage.PrintMemUsage()
			fmt.Printf("Fetching page %d: %s\n", i, pageToken)
		}

		fileList, err := srv.Files.List().
			Corpora("drive").                // Comment if using private drive
			DriveId(g.gdriveId).             // Comment if using private drive
			IncludeItemsFromAllDrives(true). // Comment if using private drive
			SupportsAllDrives(true).         // Comment if using private drive
			// IncludeItemsFromAllDrives(false). // Remove this if getting from shared drive
			PageToken(pageToken).
			Fields("files(id,name,parents),nextPageToken").
			PageSize(1000).
			OrderBy("folder,modifiedTime").
			Do()
		if err != nil {
			return nil, fmt.Errorf("could not list gdrive files: %w", err)
		}

		pageToken = fileList.NextPageToken

		files, err := g.toDriveFile(fileList.Files)
		if err != nil {
			return nil, fmt.Errorf("could not get files: %w", err)
		}

		allFiles = append(allFiles, files...)
	}

	err = g.storage.Save(filename, g.filesToString(allFiles))
	if err != nil {
		return nil, fmt.Errorf("could not save file: %w", err)
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

func (g *FileListGetter) filesToString(files []*DriveFile) string {
	s := ""
	for _, file := range files {
		s += file.String() + "\n"
	}
	return s
}
