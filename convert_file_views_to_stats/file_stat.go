package convert_file_views_to_stats

type FileStat struct {
	Id        string
	ViewCount int
	Parent    *FileStat
	Children  []*FileStat
}
