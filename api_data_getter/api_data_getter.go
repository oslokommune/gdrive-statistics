package api_data_getter

import (
	"fmt"
	"time"

	"github.com/oslokommune/gdrive-statistics/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_gdrive_views"
)

type ApiDataGetter struct {
	debug             bool
	fileListGetter    *get_file_list.FileListGetter
	gDriveViewsGetter *get_gdrive_views.GDriveViewsGetter
}

func New(
	debug bool,
	fileListGetter *get_file_list.FileListGetter,
	gDriveViewsGetter *get_gdrive_views.GDriveViewsGetter,
) *ApiDataGetter {
	return &ApiDataGetter{
		debug:             debug,
		fileListGetter:    fileListGetter,
		gDriveViewsGetter: gDriveViewsGetter,
	}
}

func (g *ApiDataGetter) Run() error {
	_, err := g.getFilesAndFolders()
	if err != nil {
		return fmt.Errorf("could not show files and folders: %w", err)
	}

	fmt.Println()

	_, err2 := g.GetViewEvents()
	if err2 != nil {
		return fmt.Errorf("could not show view events: %w", err2)
	}

	return nil
}

func (g *ApiDataGetter) getFilesAndFolders() ([]*get_file_list.DriveFile, error) {
	var pageCount int
	if g.debug {
		pageCount = 1
	} else {
		pageCount = 1000000
	}

	fmt.Printf("Getting files and folders (pageCount=%d) ...\n", pageCount)
	files, err := g.fileListGetter.GetFiles(pageCount)

	if err != nil {
		return nil, fmt.Errorf("could not get gdrive files: %w", err)
	}

	for i := 0; i < g.min(3, len(files)); i++ {
		fmt.Println(files[i])
	}

	fmt.Printf("File count: %d\n", len(files))

	return files, nil
}

func (g *ApiDataGetter) GetViewEvents() ([]*get_gdrive_views.GdriveViewEvent, error) {
	var startTime time.Time
	if g.debug {
		startTime = time.Now().AddDate(0, 0, -1)
	} else {
		startTime = time.Now().AddDate(0, -1, 0)
	}

	fmt.Printf("Getting view events (startTime=%s)...\n", startTime.Format(time.RFC3339))
	views, err := g.gDriveViewsGetter.GetGdriveDocViews(&startTime)

	if err != nil {
		return nil, fmt.Errorf("error when listing drive usage: %w", err)
	}

	for i := 0; i < g.min(3, len(views)); i++ {
		fmt.Println(views[i])
	}

	fmt.Printf("View count: %d\n", len(views))

	return views, nil
}

func (g *ApiDataGetter) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
