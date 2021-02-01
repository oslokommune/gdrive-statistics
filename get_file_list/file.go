package get_file_list

import (
	"fmt"
)

type DriveFile struct {
	Name   string
	Id     string
	Parent string
}

func (f *DriveFile) String() string {
	parent := ""
	if f.HasParent() {
		parent = fmt.Sprintf(" [parent %s]", f.Parent)
	}

	return fmt.Sprintf("DriveFile [id %s] [name %s]%s",
		f.Id,
		f.Name,
		parent,
	)
}

func (f *DriveFile) HasParent() bool {
	return len(f.Parent) > 0
}
