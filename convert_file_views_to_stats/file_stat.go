package convert_file_views_to_stats

import (
	"fmt"
)

// FileStat is a file or folder with view statistics
type FileStat struct {
	Id              string
	DocTitle        string
	ViewCount       int
	UniqueViewCount int
	Parent          *FileStat
	Children        []*FileStat
	userVisits      map[string]bool
}

func (f *FileStat) String() string {
	return fmt.Sprintf("[Id %s] [Children: %d]", f.Id, len(f.Children))
}
