package statistics_printer

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
)

type FolderViews struct {
	Name        string
	Views       int
	UniqueViews int
}

func Print(fileStats map[string]*convert_file_views_to_stats.FileStat) {
	folderViews := toFolderViews(fileStats)
	print(folderViews)
}

func print(folderViews []*FolderViews) {
	for _, v := range folderViews {
		fmt.Printf("%s", v.Name)
		fmt.Printf("       ")
		fmt.Printf("%d", v.Views)
		fmt.Printf("       ")
		fmt.Printf("%d", v.UniqueViews)
		fmt.Printf("\n")
	}
}
