package main

import (
	"fmt"
	"log"
	"os"

	"github.com/oslokommune/gdrive-statistics/api_data_getter"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/gdrive_client"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_gdrive_views"
	"github.com/oslokommune/gdrive-statistics/file_storage"
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

	files, views, err := apiDataGetter.GetDataFromApi()
	if err != nil {
		return fmt.Errorf("get data from Google API(s): %w", err)
	}

	printData(files, views)

	//fileViewStatistics := view_count_calculator.CalculateViewStatistics(files, views)
	// for every view
	//   doc = view.doc
	//   stats[doc].totalViewCount++
	//   stats[doc][view.userHash] = 1
	// func unique(doc): len(doc.viewUserHashes)

	// file id:xy123, views: 83, unique views: 50
	// file id:abc567, views: 40, unique views: 2

	//fileTree := create_file_tree.Create(files)
	// FileTree (id, name, parent)
	// parent
	// file
	// children

	// count views per folder:
	// for every view
	//   folder = get_folder_of(view.docId)
	//   viewCount[folder]

	// print stuff
	// iterer gjennom fileTree
	//   if item is folder: print("$item - $viewStatistics")

	// something something

	/*
		| Mappe                   | Antall views | Antall unike views |
		| Administrasjon/         |          500 |                100 |
		| Administrasjon/Allmøter |          300 |                 50 |
		| osv                     |          osv |                osv |
	*/

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
