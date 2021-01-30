package main

import (
	"encoding/json"
	"fmt"
	"github.com/oslokommune/gdrive-statistics/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_gdrive_views"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/reports/v1"
	"google.golang.org/api/drive/v3"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	userPkg "os/user"
	"path"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("Error while running application: %v", err)
	}
}

func run() error {
	gdriveId, ok := os.LookupEnv("GOOGLE_DRIVE_ID")
	if !ok {
		return fmt.Errorf("env need to be set: GOOGLE_DRIVE_ID")
	}

	credentialsFilepath, err := getFilepath(".google-credentials.json")
	if err != nil {
		return fmt.Errorf("could not get credentials file: %w", err)
	}

	b, err := ioutil.ReadFile(credentialsFilepath)
	if err != nil {
		return fmt.Errorf("unable to read client secret file: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, admin.AdminReportsAuditReadonlyScope, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	fmt.Println("Getting client...")
	client, err := getClient(config)
	if err != nil {
		return fmt.Errorf("could not get client: %w", err)
	}

	/* algo
	1 get folder tree structure, store into nice data structure
	2 get file views the last X months
	3 combine 1+2, create data structure according to spec
	4 print result
	*/

	// 1:

	// 2:
	views, err := get_gdrive_views.GetGdriveDocViews(client, gdriveId)
	if err != nil {
		return fmt.Errorf("error when listing drive usage: %w", err)
	}

	fmt.Println()
	fmt.Println("Google Drive views:")
	for _, view := range views {
		fmt.Println(view)
	}

	fmt.Println()
	fmt.Println("Google Drive files:")
	files, err := get_file_list.GetFiles(client, gdriveId)
	if err != nil {
		return fmt.Errorf("could not get gdrive files: %w", err)
	}

	fmt.Println(files)
	return nil
}

func getFilepath(filename string) (string, error) {
	user, err := userPkg.Current()
	if err != nil {
		return "", fmt.Errorf("unable to get user: %w", err)
	}

	userHomeDir := user.HomeDir
	return path.Join(userHomeDir, filename), nil
}

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tokFile, err := getFilepath(".google-token.json")
	if err != nil {
		return nil, fmt.Errorf("could not get file path: %w", err)
	}

	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		err = saveToken(tokFile, tok)
		if err != nil {
			return nil, fmt.Errorf("could not save token: %w", err)
		}
	}

	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	//goland:noinspection ALL
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}

	//goland:noinspection ALL
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return fmt.Errorf("could not encode token: %w", err)
	}

	fmt.Println("Token file saved")
	return nil
}
