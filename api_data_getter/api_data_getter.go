package api_data_getter

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/file_storage"
	"time"

	"github.com/oslokommune/gdrive-statistics/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_gdrive_views"
)

type ApiDataGetter struct {
	debug             bool
	fileListGetter    *get_file_list.FileListGetter
	gDriveViewsGetter *get_gdrive_views.GDriveViewsGetter
	storage           *file_storage.FileStorage
}

func New(
	debug bool,
	fileListGetter *get_file_list.FileListGetter,
	gDriveViewsGetter *get_gdrive_views.GDriveViewsGetter,
	storage *file_storage.FileStorage,
) *ApiDataGetter {
	return &ApiDataGetter{
		debug:             debug,
		fileListGetter:    fileListGetter,
		gDriveViewsGetter: gDriveViewsGetter,
		storage:           storage,
	}
}

func (g *ApiDataGetter) GetDataFromApi() ([]*get_file_list.FileOrFolder, []*get_gdrive_views.GdriveViewEvent, error) {
	files, err := g.getFilesAndFolders()
	if err != nil {
		return nil, nil, fmt.Errorf("get and store files and folders: %w", err)
	}

	views, err := g.getViewEvents()
	if err != nil {
		return nil, nil, fmt.Errorf("get and store view events: %w", err)
	}

	return files, views, err
}

func (g *ApiDataGetter) getFilesAndFolders() ([]*get_file_list.FileOrFolder, error) {
	filename := "files.json"
	fileExists, err := g.storage.AppFileExists(filename)
	if err != nil {
		return nil, fmt.Errorf("app file exists: %w", err)
	}

	if !fileExists {
		return g.GetAndStoreFilesAndFolders(filename)
	} else {
		fmt.Printf("File %s already exists, skipping API call to fetch GDrive files and folders\n", filename)
		return g.fileListGetter.LoadFromFile(filename)
	}
}

func (g *ApiDataGetter) GetAndStoreFilesAndFolders(filename string) ([]*get_file_list.FileOrFolder, error) {
	var pageCount int
	if g.debug {
		pageCount = 1
	} else {
		pageCount = 1000000
	}

	fmt.Printf("Getting files and folders (pageCount=%d) ...\n", pageCount)
	files, err := g.fileListGetter.GetAndStoreFiles(filename, pageCount)

	if err != nil {
		return nil, fmt.Errorf("get gdrive files: %w", err)
	}

	return files, nil
}

func (g *ApiDataGetter) getViewEvents() ([]*get_gdrive_views.GdriveViewEvent, error) {
	filename := "views.json"
	fileExists, err := g.storage.AppFileExists(filename)
	if err != nil {
		return nil, fmt.Errorf("app file exists: %w", err)
	}

	if !fileExists {
		return g.getAndStoreViewEvents(filename)
	} else {
		fmt.Printf("File %s already exists, skipping API call to fetch GDrive views\n", filename)
		return g.gDriveViewsGetter.LoadFromFile(filename)
	}
}

func (g *ApiDataGetter) getAndStoreViewEvents(filename string) ([]*get_gdrive_views.GdriveViewEvent, error) {
	var startTime time.Time
	if g.debug {
		startTime = time.Now().AddDate(0, 0, -2)
	} else {
		startTime = time.Now().AddDate(0, -3, 0)
	}

	fmt.Printf("Getting view events (startTime=%s)...\n", startTime.Format(time.RFC3339))
	views, err := g.gDriveViewsGetter.GetGdriveDocViews(filename, &startTime)

	if err != nil {
		return nil, fmt.Errorf("error when listing drive usage: %w", err)
	}

	return views, nil
}
