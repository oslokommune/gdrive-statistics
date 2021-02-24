package print_statistics

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
)

type Folder struct {
	docId    string
	docTitle string
	parent   *Folder
	children []*Folder

	viewCount       int
	uniqueViewcount int
}

func (f *Folder) String() string {
	return fmt.Sprintf("[DocId: %s]", f.docId)
}

func newFolderFromFile(ff *convert_file_views_to_stats.FileStat) *Folder {
	return &Folder{
		docId:           ff.Id,
		docTitle:        ff.DocTitle,
		viewCount:       ff.ViewCount,
		uniqueViewcount: ff.UniqueViewCount,
	}
}

func (f *Folder) AddChild(child *Folder) {
	f.children = append(f.children, child)
}
