package view_count_calculator

import (
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_gdrive_views"
)

type FileWithViewCount struct {
	file        *get_file_list.FileOrFolder
	views       []*get_gdrive_views.GdriveViewEvent
	uniqueViews int
}

func New(file *get_file_list.FileOrFolder) *FileWithViewCount {
	return &FileWithViewCount{
		file: file,
	}
}

func (f *FileWithViewCount) GetViewCount() int {
	return len(f.views)
}
