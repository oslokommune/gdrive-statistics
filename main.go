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

const Debug = false

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error while running application: %v", err)
	}
}

/*
	algorithm:
	1 get gdrive file and folder tree
	2 get folder tree structure, store into nice data structure
	3 get file views the last X months
	4 combine 1+2, create data structure according to spec
	5 print result
*/

func run() error {
	gDriveId, ok := os.LookupEnv("GOOGLE_DRIVE_ID")
	if !ok {
		return fmt.Errorf("env need to be set: GOOGLE_DRIVE_ID")
	}

	storage := file_storage.New()
	clientGetter := gdrive_client.New(storage)

	err := storage.CreateStoreFolderIfNotExists()
	if err != nil {
		return fmt.Errorf("could not create directory: %w", err)
	}

	fmt.Println("Getting client...")
	client, err := clientGetter.GetClient()
	if err != nil {
		return fmt.Errorf("could not get client: %w", err)
	}

	fileListGetter := get_file_list.New(client, gDriveId)
	gDriveViewsGetter := get_gdrive_views.New(client, gDriveId)

	apiDataGetter := api_data_getter.New(storage, Debug, fileListGetter, gDriveViewsGetter)

	err = apiDataGetter.Run()
	if err != nil {
		return fmt.Errorf("could not get data from Google API(s): %w", err)
	}

	return nil
}
