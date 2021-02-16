package view_count_calculator

type FileStat struct {
	Id        string
	ViewCount int
	Parent    *FileStat
	Children  []*FileStat
}
