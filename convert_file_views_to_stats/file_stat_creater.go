package convert_file_views_to_stats

import (
	"github.com/oslokommune/gdrive-statistics/get_api_data/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_api_data/get_gdrive_views"
)

func CreateFileStats(rootLevelFileId string, files []*get_file_list.FileOrFolder, views []*get_gdrive_views.GdriveViewEvent) (map[string]*FileStat, *FileStat) {
	fileStats := toFileStats(rootLevelFileId, files, views)

	root := fileStats[rootLevelFileId]
	aggregateViews(root)

	return fileStats, root
}

func toFileStats(rootLevelFile string, files []*get_file_list.FileOrFolder, views []*get_gdrive_views.GdriveViewEvent) map[string]*FileStat {
	fileStats := make(map[string]*FileStat)

	views = stripViewsThatDontHaveAFile(files, views)
	mergeFilesAndViewsToFileStats(files, fileStats, views)
	setParentsAndChildren(files, fileStats, rootLevelFile)

	return fileStats
}

// stripViewsThatDontHaveAFile removes views that doesn't have a corresponding file. These views are most likely for
// files that are not shared with the organization.
func stripViewsThatDontHaveAFile(files []*get_file_list.FileOrFolder, views []*get_gdrive_views.GdriveViewEvent) []*get_gdrive_views.GdriveViewEvent {
	fileMap := make(map[string]bool)
	for _, file := range files {
		fileMap[file.Id] = true
	}

	viewsWithFile := make([]*get_gdrive_views.GdriveViewEvent, 0)
	for _, view := range views {
		if _, ok := fileMap[view.DocId]; ok {
			viewsWithFile = append(viewsWithFile, view)

		}
	}

	return viewsWithFile
}

func mergeFilesAndViewsToFileStats(files []*get_file_list.FileOrFolder, fileStats map[string]*FileStat, views []*get_gdrive_views.GdriveViewEvent) {
	for _, file := range files {
		fileStat := FileStat{
			Id:        file.Id,
			ViewCount: 0,
			Parent:    nil,
			Children:  nil,
		}

		fileStats[fileStat.Id] = &fileStat
	}

	calc := NewUniqueViewCalculator()

	for _, view := range views {
		fileStats[view.DocId].ViewCount++
		calc.addViewForDocument(view.DocId, view.UserHash)
	}

	for docId, fs := range fileStats {
		fs.UniqueViewCount = calc.getUniqueViewsForDocument(docId)
	}
}

// setParentsAndChildren sets the parent and children pointers of the FileStats.
// The function also creates a top level root FileStat, because files from Gdrive doesn't contain the root element
// itself, just pointers to it.
func setParentsAndChildren(files []*get_file_list.FileOrFolder, fileStats map[string]*FileStat, rootLevelFileId string) {
	root := &FileStat{
		Id: rootLevelFileId,
	}
	fileStats[rootLevelFileId] = root

	for _, file := range files {
		fs := fileStats[file.Id]

		var parent *FileStat
		if file.Parent == rootLevelFileId {
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
		node.UniqueViewCount += child.UniqueViewCount
	}
}
