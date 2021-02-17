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
	root := &FileStat{
		Id: "root",
	}
	fileStats["root"] = root

	for _, file := range files {
		fs := fileStats[file.Id]

		var parent *FileStat
		if file.Parent == rootLevelFile {
			parent = root
		} else {
			parent = fileStats[file.Parent]
		}

		fs.Parent = parent
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
