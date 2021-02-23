package convert_file_views_to_stats

// FileStat is a file or folder with view statistics
type FileStat struct {
	Id              string
	ViewCount       int
	UniqueViewCount int
	Parent          *FileStat
	Children        []*FileStat
	userVisits      map[string]bool
}
