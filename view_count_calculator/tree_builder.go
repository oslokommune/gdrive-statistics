package view_count_calculator

func SetChildren(stats map[string]*FileStat) {
	for _, fs := range stats {
		if fs.Parent != nil {
			fs.Parent.Children = append(fs.Parent.Children, fs)
		}
	}
}
