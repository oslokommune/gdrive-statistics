package get_file_list

import (
	"encoding/json"
	"fmt"
)

func (g *FileListGetter) saveToFile(filename string, files []*FileOrFolder) error {
	jsonData, err := json.MarshalIndent(files, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json: %w", err)
	}

	err = g.storage.Save(filename, jsonData)
	if err != nil {
		return fmt.Errorf("save file: %w", err)
	}

	return nil
}

func (g *FileListGetter) LoadFromFile(filename string) ([]*FileOrFolder, error) {
	jsonData, err := g.storage.Load(filename)
	if err != nil {
		return nil, fmt.Errorf("load file: %w", err)
	}

	var data []FileOrFolder
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	var list []*FileOrFolder
	for _, item := range data {
		var separateItem = item

		// We have to add &separateItem and not &item, because item's value is updated by the range function for
		// every iteration. https://stackoverflow.com/questions/48826460/using-pointers-in-a-for-loop
		list = append(list, &separateItem)
	}

	return list, nil
}
