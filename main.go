package main

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/file_saver"
	"github.com/oslokommune/gdrive-statistics/gdrive_client"
	"github.com/oslokommune/gdrive-statistics/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_gdrive_views"
	"log"
	"net/http"
	"os"
)

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
	gdriveId, ok := os.LookupEnv("GOOGLE_DRIVE_ID")
	if !ok {
		return fmt.Errorf("env need to be set: GOOGLE_DRIVE_ID")
	}

	fmt.Println("Getting client...")
	client, err := gdrive_client.GetClient()
	if err != nil {
		return fmt.Errorf("could not get client: %w", err)
	}

	// 1
	//err = getFilesAndFolders(client, gdriveId)
	//if err != nil {
	//	return fmt.Errorf("could not show files and folders: %w", err)
	//}

	// 3
	err = getViewEvents(client, gdriveId)
	if err != nil {
		return fmt.Errorf("could not show view events: %w", err)
	}

	return nil
}

func getFilesAndFolders(client *http.Client, gdriveId string) error {
	fmt.Println()
	fmt.Println("Getting files and folders...")

	files, err := get_file_list.GetFiles(client, gdriveId)
	if err != nil {
		return fmt.Errorf("could not get gdrive files: %w", err)
	}

	for i := 0; i < min(3, len(files)); i++ {
		fmt.Println(files[i])
	}

	fmt.Printf("\nFile count: %d\n", len(files))

	err = file_saver.Save(filesToString(files), "files.txt")
	if err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	return nil
}

func filesToString(files []*get_file_list.DriveFile) string {
	s := ""
	for _, file := range files {
		s += file.String() + "\n"
	}
	return s
}

func getViewEvents(client *http.Client, gdriveId string) error {
	fmt.Println()
	fmt.Println("Getting view events...")

	views, err := get_gdrive_views.GetGdriveDocViews(client, gdriveId)
	if err != nil {
		return fmt.Errorf("error when listing drive usage: %w", err)
	}

	for i := 0; i < min(3, len(views)); i++ {
		fmt.Println(views[i])
	}

	fmt.Printf("View count: %d\n", len(views))

	err = file_saver.Save(viewsToString(views), "views.txt")
	if err != nil {
		return fmt.Errorf("could not save file: %w", err)
	}

	return nil
}

func viewsToString(views []*get_gdrive_views.GdriveViewEvent) string {
	s := ""
	for _, view := range views {
		s += view.String() + "\n"
	}
	return s
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
