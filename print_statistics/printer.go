package print_statistics

import (
	"fmt"
	c "github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

func Print(root *c.FileStat, maxDepth int) {
	rootFolder := toFolder(root, maxDepth)

	printer := NewColumnPrinter([]int{45, 10})
	printFolderTree(printer, rootFolder, 0)
}

func toFolder(root *c.FileStat, maxDepth int) *Folder {
	//parentDummyFolder := newFolder("root root", "root title", nil, nil, 0, 0)
	parentDummyFolder := &Folder{
		docId:    "rootOfActualRoot",
		docTitle: "parent dummy folder",
	}

	// +1 since we made root's parent to a an empty rootFolder
	findFolders(root, maxDepth+1, 0, parentDummyFolder)

	rootFolder := parentDummyFolder.children[0]
	sortFoldersByViews(rootFolder)

	return parentDummyFolder.children[0]
}

func sortFoldersByViews(folder *Folder) {
	//sortedChildren := make([]*Folder, len(folder.children))
	//
	//for i, child := range folder.children {
	//	sortedChildren[i] = child
	//}

	sort.Slice(folder.children, func(i, j int) bool {
		return folder.children[i].viewCount > folder.children[j].viewCount
		//return planets[i].Axis < planets[j].Axis
	})
}

// folders = filestats that are a parent of a child
// means we won't find bottom level folders without files, but that's okay.
func findFolders(nodeToExamine *c.FileStat, maxDepth int, currentDepth int, folderNode *Folder) {
	if currentDepth == maxDepth {
		return
	}

	if len(nodeToExamine.Children) > 0 {
		// node is a folder!
		//childFolder := newFolder(
		//	nodeToExamine.Id, nodeToExamine.DocTitle, folderNode, nil, nodeToExamine.ViewCount, nodeToExamine.UniqueViewCount)
		childFolder := newFolderFromFile(nodeToExamine)
		childFolder.parent = folderNode

		folderNode.AddChild(childFolder)

		for _, child := range nodeToExamine.Children {
			findFolders(child, maxDepth, currentDepth+1, childFolder)
		}
	}
}

func printFolderTree(printer *ColumnPrinter, f *Folder, currentDepth int) {
	if currentDepth == 0 {
		printer.add("FOLDER")
		printer.add("VIEWS")
		printer.add("UNIQUE VIEWS")
		fmt.Println(printer.get())

		printer.reset()
	}

	indent := strings.Repeat(" ", currentDepth*2)
	viewCount := rightIndent(7, strconv.Itoa(f.viewCount))
	uniqueViewCount := rightIndent(7, strconv.Itoa(f.viewCount))

	printer.add(indent + f.docTitle)
	printer.add(viewCount)
	printer.add(uniqueViewCount)

	fmt.Println(printer.get())
	printer.reset()

	for _, child := range f.children {
		printFolderTree(printer, child, currentDepth+1)
	}
}

func rightIndent(spaceForNumber int, txt string) string {
	spacesToAdd := spaceForNumber - utf8.RuneCountInString(txt)
	return strings.Repeat(" ", spacesToAdd) + txt
}
