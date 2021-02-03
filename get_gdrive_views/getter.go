package get_gdrive_views

import (
	"errors"
	"fmt"
	"github.com/oslokommune/gdrive-statistics/hasher"
	"github.com/oslokommune/gdrive-statistics/memory_usage"
	admin "google.golang.org/api/admin/reports/v1"
	"net/http"
	"time"
)

// GetGdriveDocViews returns a slice of Google Drive View events
func GetGdriveDocViews(client *http.Client, gdriveId string) ([]*GdriveViewEvent, error) {
	srv, err := admin.New(client)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve reports Client %w", err)
	}

	allViews := make([]*GdriveViewEvent, 0)
	startTime := time.Now().
		//AddDate(0, -3, 0).
		AddDate(0, 0, -1).
		Format(time.RFC3339)

	nextPageToken := ""
	i := 0

	for ok := true; ok; ok = nextPageToken != "" && i < 3 {
		i++

		if len(nextPageToken) > 0 {
			memory_usage.PrintMemUsage()
			fmt.Printf("Fetching page %d: %s\n", i, nextPageToken[len(nextPageToken)-10:])
		}

		activities, err := srv.Activities.
			List("all", "drive").
			MaxResults(1000).
			EventName("view").
			Filters(fmt.Sprintf("shared_drive_id==%s", gdriveId)).
			StartTime(startTime).
			Do()
		if err != nil {
			return nil, fmt.Errorf("could not get activities: %w", err)
		}

		nextPageToken = activities.NextPageToken

		views, err := getViews(activities)
		if err != nil {
			return nil, fmt.Errorf("error getting views: %w", err)
		}

		allViews = append(allViews, views...)
	}

	return allViews, nil
}

func getViews(activities *admin.Activities) ([]*GdriveViewEvent, error) {
	docViews := make([]*GdriveViewEvent, 0)

	for _, item := range activities.Items {
		view, err := createDocView(item)
		if err != nil {
			return nil, err
		}

		docViews = append(docViews, view)
	}

	return docViews, nil
}

func createDocView(activity *admin.Activity) (*GdriveViewEvent, error) {
	itemTime, err := time.Parse(time.RFC3339Nano, activity.Id.Time)
	if err != nil {
		return nil, fmt.Errorf("unable to parse time: %w", err)
	}

	mainEvent, err := getMainEvent(activity)
	if err != nil {
		return nil, err
	}

	docId, err := mainEvent.GetField("doc_id")
	if err != nil {
		return nil, fmt.Errorf("could not get field: %w", err)
	}

	docTitle, err := mainEvent.GetField("doc_title")
	if err != nil {
		return nil, fmt.Errorf("could not get field: %w", err)
	}

	userHash := hasher.NewHash(activity.Actor.Email)

	//fmt.Printf("%d [%s]: %s <- %s \t\t doc_id: %s doc_title: %s \t\t\t shared_drive_id: %s \n", i, itemTime.Format(time.RFC822), mainEvent.Name, item.Actor.Email, docId, docTitle, sharedDriveId)
	view := &GdriveViewEvent{
		time:     &itemTime,
		userHash: userHash,
		docId:    docId,
		docTitle: docTitle,
	}

	return view, nil
}

func getMainEvent(item *admin.Activity) (eventParameters, error) {
	if len(item.Events) == 0 {
		return eventParameters{}, errors.New(fmt.Sprintf("got 0 events for item with Etag %s", item.Etag))
	}

	if len(item.Events) > 1 {
		return eventParameters{}, errors.New(fmt.Sprintf("got more than 1 event for item with Etag %s", item.Etag))
	}

	mainEvent := newEventParameters(item.Events[0])
	return mainEvent, nil
}
