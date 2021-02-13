package get_gdrive_views

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/file_storage"
	"net/http"
	"time"

	"github.com/oslokommune/gdrive-statistics/memory_usage"
	admin "google.golang.org/api/admin/reports/v1"
)

type GDriveViewsGetter struct {
	client   *http.Client
	gDriveId string
	storage  *file_storage.FileStorage
}

func New(client *http.Client, gDriveId string, storage *file_storage.FileStorage) *GDriveViewsGetter {
	return &GDriveViewsGetter{
		client:   client,
		gDriveId: gDriveId,
		storage:  storage,
	}
}

// GetGdriveDocViews fetches View events from the Google Reports API
func (v *GDriveViewsGetter) GetGdriveDocViews(filename string, startTime *time.Time) ([]*GdriveViewEvent, error) {
	views, err := v.getViewsFromApi(startTime)
	if err != nil {
		return nil, fmt.Errorf("call gdrive api: %w", err)
	}

	err = v.saveToFile(filename, views)
	if err != nil {
		return nil, fmt.Errorf("save views to file: %w", err)
	}

	return views, nil
}

func (v *GDriveViewsGetter) getViewsFromApi(startTime *time.Time) ([]*GdriveViewEvent, error) {
	//goland:noinspection ALL
	srv, err := admin.New(v.client)
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve reports Client %w", err)
	}

	allViews := make([]*GdriveViewEvent, 0)
	startTimeStr := startTime.Format(time.RFC3339)

	pageToken := ""
	i := 0

	for ok := true; ok; ok = pageToken != "" && i < 1000000 {
		i++

		if len(pageToken) > 0 {
			memory_usage.PrintMemUsage()
			fmt.Printf("Fetching page %d: %s\n", i, pageToken)
		}

		activitiesListCall := srv.Activities.
			List("all", "drive").
			MaxResults(1000).
			EventName("view").
			Filters(fmt.Sprintf("shared_drive_id==%s", v.gDriveId)).
			StartTime(startTimeStr)

		if len(pageToken) > 0 {
			activitiesListCall.PageToken(pageToken)
		}

		activities, err := activitiesListCall.Do()

		if err != nil {
			return nil, fmt.Errorf("get activities: %w", err)
		}

		pageToken = activities.NextPageToken

		views, err := v.toViews(activities)
		if err != nil {
			return nil, fmt.Errorf("error getting views: %w", err)
		}

		allViews = append(allViews, views...)
	}

	return allViews, nil
}
