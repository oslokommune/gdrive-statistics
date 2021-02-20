package get_gdrive_views

import (
	"encoding/json"
	"fmt"
)

func (v *GDriveViewsGetter) saveToFile(filename string, views []*GdriveViewEvent) error {
	jsonData, err := json.MarshalIndent(views, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	err = v.storage.Save(filename, jsonData)
	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func (v *GDriveViewsGetter) LoadFromFile(filename string) ([]*GdriveViewEvent, error) {
	jsonData, err := v.storage.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("load file: %w", err)
	}

	var data []GdriveViewEvent
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	var list []*GdriveViewEvent
	for _, item := range data {
		var separateItem = item

		// We have to add &separateItem and not &item, because item's value is updated by the range function for
		// every iteration. https://stackoverflow.com/questions/48826460/using-pointers-in-a-for-loop
		list = append(list, &separateItem)
	}

	return list, nil
}
