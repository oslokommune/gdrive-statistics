package print_statistics

import (
	"fmt"
	c "github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"strings"
)

func Print(root *c.FileStat, maxDepth int) {
	rootFolder := toFolder(root, maxDepth)
	printFolderTree(rootFolder, 0)
	//folderViews := toFolderViews(root, 1, 0)
	//printViews(folderViews)
}

func toFolder(root *c.FileStat, maxDepth int) *Folder {
	parentDummyFolder := newFolder("root root", nil, nil, 0, 0)

	// +1 since we made root's parent to a an empty rootFolder
	findFolders(root, maxDepth+1, 0, parentDummyFolder)

	return parentDummyFolder.children[0]
}

// folders = filestats that are a parent of a child
// means we won't find bottom level folders without files, but that's okay.
func findFolders(nodeToExamine *c.FileStat, maxDepth int, currentDepth int, folderNode *Folder) {
	if currentDepth == maxDepth {
		return
	}

	if len(nodeToExamine.Children) > 0 {
		// node is a folder!
		childFolder := newFolder(
			nodeToExamine.Id, folderNode, nil, nodeToExamine.ViewCount, nodeToExamine.UniqueViewCount)
		folderNode.AddChild(childFolder)

		for _, child := range nodeToExamine.Children {
			findFolders(child, maxDepth, currentDepth+1, childFolder)
		}
	}
}

func printFolderTree(f *Folder, currentDepth int) {
	if currentDepth == 0 {
		fmt.Println("FOLDER\t\tVIEWS\t\tUNIQUE VIEWS")
	}
	emptySpace := strings.Repeat(" ", currentDepth*2)

	fmt.Print(emptySpace)
	fmt.Print(f.docId)
	fmt.Print("\t\t\t")
	fmt.Print(f.viewCount)
	fmt.Print("\t\t\t")
	fmt.Print(f.uniqueViewcount)
	fmt.Print("\n")

	//fmt.Printf("%s%s\t\t%d\t\t%d\n", emptySpace, f.docId, f.viewCount, f.uniqueViewcount)

	for _, child := range f.children {
		printFolderTree(child, currentDepth+1)
	}
}
