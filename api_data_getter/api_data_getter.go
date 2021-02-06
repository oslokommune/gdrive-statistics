package api_data_getter

import (
	"fmt"

	"github.com/oslokommune/gdrive-statistics/file_storage"
	"github.com/oslokommune/gdrive-statistics/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_gdrive_views"
)

type ApiDataGetter struct {
	storage           *file_storage.FileStorage
	debug             bool
	fileListGetter    *get_file_list.FileListGetter
	gDriveViewsGetter *get_gdrive_views.GDriveViewsGetter
}

func New(
	storage *file_storage.FileStorage,
	debug bool,
	fileListGetter *get_file_list.FileListGetter,
	gDriveViewsGetter *get_gdrive_views.GDriveViewsGetter,
) *ApiDataGetter {
	return &ApiDataGetter{
		storage:           storage,
		debug:             debug,
		fileListGetter:    fileListGetter,
		gDriveViewsGetter: gDriveViewsGetter,
	}
}

func (g *ApiDataGetter) Run() error {
	// 1
	err := g.getFilesAndFolders()
	if err != nil {
		return fmt.Errorf("could not show files and folders: %w", err)
	}

	fmt.Println()

	// 3
	err = g.getViewEvents()
	if err != nil {
		return fmt.Errorf("could not show view events: %w", err)
	}

	return nil
}

func (g *ApiDataGetter) getFilesAndFolders() error {
	var pageCount int
	if g.debug {
		pageCount = 1
	} else {
		pageCount = 1000000
	}

	fmt.Printf("Getting files and folders (pageCount=%d) ...\n", pageCount)
	files, err := g.fileListGetter.GetFiles(pageCount)

	if err != nil {
		return fmt.Errorf("could not get gdrive files: %w", err)
	}

	for i := 0; i < g.min(3, len(files)); i++ {
		fmt.Println(files[i])
	}

	fmt.Printf("\nFile count: %d\n", len(files))

	err = g.storage.Save("files.txt", g.filesToString(files))
	if err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	return nil
}

func (g *ApiDataGetter) filesToString(files []*get_file_list.DriveFile) string {
	s := ""
	for _, file := range files {
		s += file.String() + "\n"
	}
	return s
}

func (g *ApiDataGetter) getViewEvents() error {
	var pageCount int
	if g.debug {
		pageCount = 1
	} else {
		pageCount = 1000000
	}

	fmt.Printf("Getting view events (pageCount=%d)...\n", pageCount)

	views, err := g.gDriveViewsGetter.GetGdriveDocViews(pageCount)

	if err != nil {
		return fmt.Errorf("error when listing drive usage: %w", err)
	}

	for i := 0; i < g.min(3, len(views)); i++ {
		fmt.Println(views[i])
	}

	fmt.Printf("View count: %d\n", len(views))

	err = g.storage.Save("views.txt", g.viewsToString(views))
	if err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	return nil
}

func (g *ApiDataGetter) viewsToString(views []*get_gdrive_views.GdriveViewEvent) string {
	s := ""
	for _, view := range views {
		s += view.String() + "\n"
	}
	return s
}

func (g *ApiDataGetter) min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
