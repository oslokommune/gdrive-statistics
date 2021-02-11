package gdrive_client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/oslokommune/gdrive-statistics/file_storage"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/reports/v1"
	"google.golang.org/api/drive/v3"
)

type ClientGetter struct {
	storage *file_storage.FileStorage
}

func New(storage *file_storage.FileStorage) *ClientGetter {
	return &ClientGetter{
		storage: storage,
	}
}

func (g *ClientGetter) GetClient() (*http.Client, error) {
	credentialsFilepath, err := g.storage.GetFilepath(".google-credentials.json")
	if err != nil {
		return nil, fmt.Errorf("get credentials file: %w", err)
	}

	b, err := ioutil.ReadFile(credentialsFilepath)
	if err != nil {
		return nil, fmt.Errorf("unable to read client secret file: %w", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, admin.AdminReportsAuditReadonlyScope, drive.DriveMetadataReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("unable to parse client secret file to config: %w", err)
	}

	client, err := g.getHttpClient(config)
	if err != nil {
		return nil, fmt.Errorf("get client: %w", err)
	}

	return client, nil
}

// Retrieve a token, saves the token, then returns the generated client.
func (g *ClientGetter) getHttpClient(config *oauth2.Config) (*http.Client, error) {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.

	tokFile, err := g.storage.GetFilepath(".google-token.json")
	if err != nil {
		return nil, fmt.Errorf("get file path: %w", err)
	}

	tok, err := g.tokenFromFile(tokFile)
	if err != nil {
		tok = g.getTokenFromWeb(config)
		err = g.saveToken(tokFile, tok)
		if err != nil {
			return nil, fmt.Errorf("save token: %w", err)
		}
	}

	return config.Client(context.Background(), tok), nil
}

// Retrieves a token from a local file.
func (g *ClientGetter) tokenFromFile(file string) (*oauth2.Token, error) {
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
func (g *ClientGetter) getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
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
func (g *ClientGetter) saveToken(path string, token *oauth2.Token) error {
	fmt.Printf("Saving credential file to: %s\n", path)

	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0o600)
	if err != nil {
		return fmt.Errorf("unable to cache oauth token: %w", err)
	}

	//goland:noinspection ALL
	defer f.Close()

	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		return fmt.Errorf("encode token: %w", err)
	}

	fmt.Println("Token file saved")
	return nil
}
