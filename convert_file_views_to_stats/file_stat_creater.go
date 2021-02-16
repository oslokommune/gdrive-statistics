package convert_file_views_to_stats

import (
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_gdrive_views"
)

func CreateFileStats(files []*get_file_list.FileOrFolder, views []*get_gdrive_views.GdriveViewEvent) map[string]*FileStat {
	fileStats := toFileStats(files)
	// replace root parent gdrive id with nil
	SetChildren(fileStats)
	return fileStats
}

func toFileStats([]*get_file_list.FileOrFolder) map[string]*FileStat {
	return nil
}
