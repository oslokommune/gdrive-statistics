package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oslokommune/gdrive-statistics/api_data_getter"
	"github.com/oslokommune/gdrive-statistics/file_storage"
	"github.com/oslokommune/gdrive-statistics/gdrive_client"
	"github.com/oslokommune/gdrive-statistics/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_gdrive_views"
)

const Debug = true

//const Debug = false

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

	fileListGetter := get_file_list.New(client, gDriveId, storage)
	gDriveViewsGetter := get_gdrive_views.New(client, gDriveId, storage)

	apiDataGetter := api_data_getter.New(Debug, fileListGetter, gDriveViewsGetter, storage)

	err = getAndProcessApiData(apiDataGetter)
	if err != nil {
		return fmt.Errorf("get data from api: %w", err)
	}

	return nil
}

/*
	algorithm:
	1 get gdrive file and folder tree
	2 get folder tree structure, store into nice data structure
	3 get file views the last X months
	4 combine 1+2, create data structure according to spec
	5 print result
*/

func getAndProcessApiData(apiDataGetter *api_data_getter.ApiDataGetter) error {
	files, views, err := apiDataGetter.GetDataFromApi()
	if err != nil {
		return fmt.Errorf("get data from Google API(s): %w", err)
	}

	itemCountToPrint := 7

	for i := 0; i < min(itemCountToPrint, len(files)); i++ {
		fmt.Println(files[i])
	}

	fmt.Printf("File count: %d\n", len(files))

	for i := 0; i < min(itemCountToPrint, len(views)); i++ {
		fmt.Println(views[i])
	}

	fmt.Printf("View count: %d\n", len(views))

	return err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
