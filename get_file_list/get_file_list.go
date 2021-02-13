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
func (g *FileListGetter) GetAndStoreFiles(filename string, pageCount int) ([]*FileOrFolder, error) {
	files, err := g.getFilesFromApi(pageCount)
	if err != nil {
		return nil, fmt.Errorf("call gdrive api: %w", err)
	}

	err = g.saveToFile(filename, files)
	if err != nil {
		return nil, fmt.Errorf("save file list to file: %w", err)
	}

	return files, nil
}

func (g *FileListGetter) getFilesFromApi(pageCount int) ([]*FileOrFolder, error) {
	//goland:noinspection ALL
	srv, err := drive.New(g.client)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve Drive client: %v", err)
	}

	allFiles := make([]*FileOrFolder, 0)
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
			return nil, fmt.Errorf("list gdrive files: %w", err)
		}

		pageToken = fileList.NextPageToken

		files, err := g.toDriveFile(fileList.Files)
		if err != nil {
			return nil, fmt.Errorf("get files: %w", err)
		}

		allFiles = append(allFiles, files...)
	}

	return allFiles, nil
}
