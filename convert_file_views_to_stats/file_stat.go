package convert_file_views_to_stats

// FileStat is a file with view statistics
type FileStat struct {
	Id        string
	ViewCount int
	Parent    *FileStat
	Children  []*FileStat
}
