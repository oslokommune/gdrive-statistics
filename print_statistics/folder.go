package print_statistics

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
)

type Folder struct {
	fileStat *convert_file_views_to_stats.FileStat
	docId    string
	docTitle string
	parent   *Folder
	children []*Folder

	viewCount       int
	uniqueViewcount int
}

func (f *Folder) String() string {
	return fmt.Sprintf("[DocId: %s] [DocTitle: %s] [Views: %d / %d]", f.docId, f.docTitle, f.viewCount, f.uniqueViewcount)
}

func newFolderFromFile(ff *convert_file_views_to_stats.FileStat) *Folder {
	return &Folder{
		fileStat:        ff,
		docId:           ff.Id,
		docTitle:        ff.DocTitle,
		viewCount:       ff.ViewCount,
		uniqueViewcount: ff.UniqueViewCount,
	}
}

func (f *Folder) AddChild(child *Folder) {
	f.children = append(f.children, child)
}
