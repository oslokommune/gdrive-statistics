package main

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"github.com/oslokommune/gdrive-statistics/file_storage"
	"github.com/oslokommune/gdrive-statistics/get_api_data"
	"github.com/oslokommune/gdrive-statistics/get_api_data/gdrive_client"
	"github.com/oslokommune/gdrive-statistics/get_api_data/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_api_data/get_gdrive_views"
	"github.com/oslokommune/gdrive-statistics/print_statistics"
	"log"
	"os"
	"time"
)

const Debug = false
const sharedDrive = true
const maxFolderDepth = 1

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error while running application: %v", err)
	}
}

func run() error {
	gDriveId, ok := os.LookupEnv("GOOGLE_DRIVE_ID")
	if !ok {
		return fmt.Errorf("env need to be set: GOOGLE_DRIVE_ID")
	}

	storage := file_storage.New()
	clientGetter := gdrive_client.New(storage)

	err := storage.CreateStoreFolderIfNotExists()
	if err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	fmt.Println("Getting client...")
	client, err := clientGetter.GetClient()
	if err != nil {
		return fmt.Errorf("get client: %w", err)
	}

	fileListGetter := get_file_list.New(client, gDriveId, storage, sharedDrive)
	gDriveViewsGetter := get_gdrive_views.New(client, gDriveId, storage)

	var startTime time.Time
	if Debug {
		startTime = time.Now().AddDate(0, 0, -2)
	} else {
		startTime = time.Now().AddDate(0, -3, 0)
	}
	apiDataGetter := get_api_data.New(Debug, fileListGetter, gDriveViewsGetter, storage)

	files, views, err := apiDataGetter.GetDataFromApi(startTime)
	if err != nil {
		return fmt.Errorf("get data from Google API(s): %w", err)
	}

	printData(files, views)

	_, root := convert_file_views_to_stats.CreateFileStats(gDriveId, files, views)
	print_statistics.Print(root, maxFolderDepth)

	return nil
}

func printData(files []*get_file_list.FileOrFolder, views []*get_gdrive_views.GdriveViewEvent) {
	itemCountToPrint := 7

	for i := 0; i < min(itemCountToPrint, len(files)); i++ {
		fmt.Println(files[i])
	}

	fmt.Printf("File count: %d\n", len(files))

	for i := 0; i < min(itemCountToPrint, len(views)); i++ {
		fmt.Println(views[i])
	}

	fmt.Printf("View count: %d\n", len(views))
	fmt.Println(`This view count may not equal to the root folder's view acount below. This is because the total
view count includes views for files that are hidden, i.e. not shared with other users (I think). More specifically,
these are views that doesn't have a corresponding file, as the Gdrive file API doesn't return these hidden files,

Also, if you show multiple levels in the folder tree below, you might wonder why the views of a folder doesn't equal
the sum of views of its children. This is because the folder itself contains files that have views.`)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
