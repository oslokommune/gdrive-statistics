package folder_stats_calculator

import (
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_gdrive_views"
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
)

func CalculateViewStatistics([]*get_file_list.FileOrFolder, []*get_gdrive_views.GdriveViewEvent) string {
	return ""
}

func AggregateViews(node *convert_file_views_to_stats.FileStat) {
	if len(node.Children) == 0 {
		return
	}

	for _, child := range node.Children {
		AggregateViews(child)
		node.ViewCount += child.ViewCount
	}
}
