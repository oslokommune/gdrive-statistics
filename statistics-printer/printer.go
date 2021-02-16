package statistics_printer

import (
	"fmt"
)

type FolderViews struct {
	Name        string
	Views       int
	UniqueViews int
}

func Print(folderViews []*FolderViews) {
	for _, v := range folderViews {
		fmt.Printf("%s")
		fmt.Printf("       ")
		fmt.Printf("%d", v.Views)
		fmt.Printf("       ")
		fmt.Printf("%d", v.UniqueViews)
		fmt.Printf("\n")
	}
}
