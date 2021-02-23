package convert_file_views_to_stats

// FileStat is a file or folder with view statistics
type FileStat struct {
	Id              string
	ViewCount       int
	UniqueViewCount int
	Parent          *FileStat
	Children        []*FileStat

	userVisits map[string]bool
}

//func (s *FileStat) AddViewForUser(userHash *hasher.Hash) {
//	if s.userVisits == nil {
//		s.userVisits = make(map[string]bool)
//	}
//	s.userVisits[userHash.String()] = true
//}
//
//func (s *FileStat) UniqueViewCount() int {
//	return len(s.userVisits)
//}
