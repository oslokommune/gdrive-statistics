package get_file_list

import (
	"fmt"
)

// FileOrFolder is a file or folder
type FileOrFolder struct {
	Id     string
	Name   string
	Parent string
}

func (f *FileOrFolder) String() string {
	parent := ""
	if f.HasParent() {
		parent = fmt.Sprintf(" [parent %s]", f.Parent)
	}

	return fmt.Sprintf("File [id %s] [name %s]%s",
		f.Id,
		f.Name,
		parent,
	)
}

func (f *FileOrFolder) HasParent() bool {
	return len(f.Parent) > 0
}
