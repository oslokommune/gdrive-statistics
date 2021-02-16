package view_count_calculator

import (
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_gdrive_views"
)

func CalculateViewStatistics([]*get_file_list.FileOrFolder, []*get_gdrive_views.GdriveViewEvent) string {
	return ""
}

func AggregateViews(node *FileStat) {
	if len(node.Children) == 0 {
		return
	}

	for _, child := range node.Children {
		AggregateViews(child)
		node.ViewCount += child.ViewCount
	}
}
