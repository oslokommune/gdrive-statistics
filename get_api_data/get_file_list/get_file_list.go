package get_file_list

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/file_storage"
	"net/http"

	"github.com/oslokommune/gdrive-statistics/calc_memory_usage"
	"google.golang.org/api/drive/v3"
)

type FileListGetter struct {
	client      *http.Client
	gdriveId    string
	storage     *file_storage.FileStorage
	sharedDrive bool
}

func New(client *http.Client, gdriveId string, storage *file_storage.FileStorage, sharedDrive bool) *FileListGetter {
	return &FileListGetter{
		client:      client,
		gdriveId:    gdriveId,
		storage:     storage,
		sharedDrive: sharedDrive,
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
			calc_memory_usage.PrintMemUsage()
			fmt.Printf("Fetching page %d: %s\n", i, pageToken)
		}

		fileListRequest := srv.Files.List().
			PageToken(pageToken).
			Fields("files(id,name,parents),nextPageToken").
			PageSize(1000).
			OrderBy("folder,modifiedTime")

		if g.sharedDrive {
			fileListRequest.
				Corpora("drive").
				DriveId(g.gdriveId).
				IncludeItemsFromAllDrives(true).
				SupportsAllDrives(true)
		} else {
			fileListRequest.IncludeItemsFromAllDrives(false)
		}

		fileList, err := fileListRequest.Do()
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
