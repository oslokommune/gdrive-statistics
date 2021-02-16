package convert_file_views_to_stats

func SetChildren(stats map[string]*FileStat) {
	for _, fs := range stats {
		if fs.Parent != nil {
			fs.Parent.Children = append(fs.Parent.Children, fs)
		}
	}
}
