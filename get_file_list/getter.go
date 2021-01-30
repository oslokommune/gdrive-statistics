package get_file_list

import (
	"fmt"
	"google.golang.org/api/drive/v3"
	"log"
	"net/http"
)

func GetFiles(client *http.Client, gdriveId string) ([]*drive.File, error) {
	//goland:noinspection GoDeprecation
	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	fileList, err := srv.Files.List().
		Corpora("drive").
		DriveId(gdriveId).
		IncludeItemsFromAllDrives(true).
		SupportsAllDrives(true).
		PageSize(3).
		// TODO: Ordering og/eller pageToken
		Do()
	if err != nil {
		return nil, fmt.Errorf("could not list gdrive files: %w", err)
	}

	fmt.Println()
	fmt.Printf("Got files. Count: %d\n", len(fileList.Files))

	//data, err := getSomeSmartDataStructure(fileList)
	//if err != nil {
	//	return nil, fmt.Errorf("error getting views: %w", err)
	//}

	return nil, nil
}
