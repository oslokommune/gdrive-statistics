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
)

const Debug = false
const sharedDrive = true
const maxFolderDepth = 2

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

	apiDataGetter := get_api_data.New(Debug, fileListGetter, gDriveViewsGetter, storage)

	files, views, err := apiDataGetter.GetDataFromApi()
	if err != nil {
		return fmt.Errorf("get data from Google API(s): %w", err)
	}

	printData(files, views)

	_, root := convert_file_views_to_stats.CreateFileStats(gDriveId, files, views)
	print_statistics.Print(root, maxFolderDepth)

	/*
		| Mappe                   | Antall views | Antall unike views |
		| Administrasjon/         |          500 |                100 |
		| Administrasjon/Allmøter |          300 |                 50 |
		| osv                     |          osv |                osv |
	*/

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
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
