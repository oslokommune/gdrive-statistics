package get_gdrive_views

import (
	"errors"
	"fmt"
	"github.com/oslokommune/gdrive-statistics/hasher"
	admin "google.golang.org/api/admin/reports/v1"
	"time"
)

func (v *GDriveViewsGetter) toViews(activities *admin.Activities) ([]*GdriveViewEvent, error) {
	docViews := make([]*GdriveViewEvent, 0)

	for _, item := range activities.Items {
		view, err := v.createDocView(item)
		if err != nil {
			return nil, err
		}

		docViews = append(docViews, view)
	}

	return docViews, nil
}

func (v *GDriveViewsGetter) createDocView(activity *admin.Activity) (*GdriveViewEvent, error) {
	itemTime, err := time.Parse(time.RFC3339Nano, activity.Id.Time)
	if err != nil {
		return nil, fmt.Errorf("unable to parse time: %w", err)
	}

	mainEvent, err := v.getMainEvent(activity)
	if err != nil {
		return nil, err
	}

	docId, err := mainEvent.GetField("doc_id")
	if err != nil {
		return nil, fmt.Errorf("get field: %w", err)
	}

	docTitle, err := mainEvent.GetField("doc_title")
	if err != nil {
		return nil, fmt.Errorf("get field: %w", err)
	}

	userHash := hasher.NewHash(activity.Actor.Email)

	// fmt.Printf("%d [%s]: %s <- %s \t\t doc_id: %s doc_title: %s \t\t\t shared_drive_id: %s \n", i, itemTime.Format(time.RFC822), mainEvent.Name, item.Actor.Email, docId, docTitle, sharedDriveId)
	view := &GdriveViewEvent{
		Time:     &itemTime,
		UserHash: userHash,
		DocId:    docId,
		DocTitle: docTitle,
	}

	return view, nil
}

func (v *GDriveViewsGetter) getMainEvent(item *admin.Activity) (eventParameters, error) {
	if len(item.Events) == 0 {
		return eventParameters{}, errors.New(fmt.Sprintf("got 0 events for item with Etag %s", item.Etag))
	}

	if len(item.Events) > 1 {
		return eventParameters{}, errors.New(fmt.Sprintf("got more than 1 event for item with Etag %s", item.Etag))
	}

	mainEvent := newEventParameters(item.Events[0])
	return mainEvent, nil
}

func (v *GDriveViewsGetter) viewsToString(views []*GdriveViewEvent) string {
	s := ""
	for _, view := range views {
		s += view.String() + "\n"
	}
	return s
}
