package main

import (
	"fmt"
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
	1 get folder tree structure, store into nice data structure
	2 get file views the last X months
	3 combine 1+2, create data structure according to spec
	4 print result
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
	err = showFilesAndFolders(client, gdriveId)
	if err != nil {
		return fmt.Errorf("could not show files and folders: %w", err)
	}

	// 2
	//err = printViewEvents(client, gdriveId)
	//if err != nil {
	//	return fmt.Errorf("could not show view events: %w", err)
	//}

	return nil
}

func showFilesAndFolders(client *http.Client, gdriveId string) error {
	fmt.Println()
	fmt.Println("Getting files and folders...")

	files, err := get_file_list.GetFiles(client, gdriveId)
	if err != nil {
		return fmt.Errorf("could not get gdrive files: %w", err)
	}

	fmt.Printf("Google Drive files (%d):\n", len(files))
	for _, view := range files {
		fmt.Println(view)
	}

	return nil
}

func printViewEvents(client *http.Client, gdriveId string) error {
	fmt.Println()
	fmt.Println("Getting view events...")

	views, err := get_gdrive_views.GetGdriveDocViews(client, gdriveId)
	if err != nil {
		return fmt.Errorf("error when listing drive usage: %w", err)
	}

	fmt.Printf("Google Drive views (%d):\n", len(views))
	for _, view := range views {
		fmt.Println(view)
	}

	return nil
}
