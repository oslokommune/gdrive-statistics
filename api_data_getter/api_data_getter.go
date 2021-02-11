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

func (g *ApiDataGetter) Run() (string, string, error) {
	gdriveFiledataFilename, err := g.getAndStoreFilesAndFolders()
	if err != nil {
		return "", "", fmt.Errorf("get and store files and folders: %w", err)
	}

	fmt.Println(gdriveFiledataFilename)

	viewsFilename, err := g.getAndStoreViewEvents()
	if err != nil {
		return "", "", fmt.Errorf("get and store view events: %w", err)
	}

	fmt.Println(viewsFilename)

	return gdriveFiledataFilename, viewsFilename, err
}

func (g *ApiDataGetter) getAndStoreFilesAndFolders() (string, error) {
	filename := "files.json"
	fileExists, err := g.storage.AppFileExists(filename)
	if err != nil {
		return "", fmt.Errorf("app file exists: %w", err)
	}

	if fileExists {
		fmt.Printf("File %s already exists, skipping API call to fetch GDrive files and folders\n", filename)
	} else {
		err := g.GetAndStoreFilesAndFolders(filename)
		if err != nil {
			return "", fmt.Errorf("show files and folders: %w", err)
		}
	}

	return filename, nil
}

func (g *ApiDataGetter) getAndStoreViewEvents() (string, error) {
	filename := "views.json"
	fileExists, err := g.storage.AppFileExists(filename)
	if err != nil {
		return "", fmt.Errorf("app file exists: %w", err)
	}

	if fileExists {
		fmt.Printf("File %s already exists, skipping API call to fetch GDrive views\n", filename)
	} else {
		err := g.GetAndStoreViewEvents(filename)
		if err != nil {
			return "", fmt.Errorf("show view events: %w", err)
		}
	}

	return filename, nil
}

func (g *ApiDataGetter) GetAndStoreFilesAndFolders(filename string) error {
	var pageCount int
	if g.debug {
		pageCount = 1
	} else {
		pageCount = 1000000
	}

	fmt.Printf("Getting files and folders (pageCount=%d) ...\n", pageCount)
	files, err := g.fileListGetter.GetAndStoreFiles(filename, pageCount)

	if err != nil {
		return fmt.Errorf("get gdrive files: %w", err)
	}

	for i := 0; i < g.min(3, len(files)); i++ {
		fmt.Println(files[i])
	}

	fmt.Printf("File count: %d\n", len(files))

	return nil
}

func (g *ApiDataGetter) GetAndStoreViewEvents(filename string) error {
	var startTime time.Time
	if g.debug {
		startTime = time.Now().AddDate(0, 0, -2)
	} else {
		startTime = time.Now().AddDate(0, -3, 0)
	}

	fmt.Printf("Getting view events (startTime=%s)...\n", startTime.Format(time.RFC3339))
	views, err := g.gDriveViewsGetter.GetGdriveDocViews(filename, &startTime)

	if err != nil {
		return fmt.Errorf("error when listing drive usage: %w", err)
	}

	for i := 0; i < g.min(3, len(views)); i++ {
		fmt.Println(views[i])
	}

	fmt.Printf("View count: %d\n", len(views))

	return nil
}

func (g *ApiDataGetter) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
