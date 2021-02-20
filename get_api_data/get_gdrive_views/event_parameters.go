package get_gdrive_views

import (
	"fmt"

	admin "google.golang.org/api/admin/reports/v1"
)

type eventParameters struct {
	data map[string]string
}

func newEventParameters(event *admin.ActivityEvents) eventParameters {
	data := make(map[string]string)

	for _, param := range event.Parameters {
		data[param.Name] = param.Value
	}

	return eventParameters{
		data: data,
	}
}

func (e *eventParameters) GetField(key string) (string, error) {
	val, found := e.data[key]
	if !found {
		return "", fmt.Errorf("field doesn't exist: %s", key)
	}

	return val, nil
}
