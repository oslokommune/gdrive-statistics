package convert_file_views_to_stats

import (
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_gdrive_views"
)

func CreateFileStats(rootLevelFileId string, files []*get_file_list.FileOrFolder, views []*get_gdrive_views.GdriveViewEvent) map[string]*FileStat {
	fileStats := toFileStats(rootLevelFileId, files)

	root := fileStats[rootLevelFileId]
	aggregateViews(root)

	return fileStats
}

func toFileStats(rootLevelFile string, files []*get_file_list.FileOrFolder) map[string]*FileStat {
	fileStats := make(map[string]*FileStat)

	convertFilesToFileStats(files, fileStats)
	setParentsAndChildren(files, fileStats, rootLevelFile)

	return fileStats
}

func convertFilesToFileStats(files []*get_file_list.FileOrFolder, fileStats map[string]*FileStat) {
	for _, file := range files {
		fileStat := FileStat{
			Id:        file.Id,
			ViewCount: 0,
			Parent:    nil,
			Children:  nil,
		}

		fileStats[fileStat.Id] = &fileStat
	}
}

func setParentsAndChildren(files []*get_file_list.FileOrFolder, fileStats map[string]*FileStat, rootLevelFile string) {
	for _, file := range files {
		fs := fileStats[file.Id]

		if file.Parent == rootLevelFile {
			// fileStats[rootLevelFile] doesn't exist, as it is the root, so we're skipping it
			continue
		}

		fs.Parent = fileStats[file.Parent]
		fs.Parent.Children = append(fs.Parent.Children, fs)
	}
}

func aggregateViews(node *FileStat) {
	if len(node.Children) == 0 {
		return
	}

	for _, child := range node.Children {
		aggregateViews(child)
		node.ViewCount += child.ViewCount
	}
}
