package print_statistics

import (
	"fmt"
	conv "github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFileStats(t *testing.T) {
	t.Run("Should find two folders", func(t *testing.T) {
		root := &conv.FileStat{Id: "root"}

		a := &conv.FileStat{Id: "a", Parent: root}
		root.Children = append(root.Children, a)

		b := &conv.FileStat{Id: "b", Parent: a}
		a.Children = append(a.Children, b)

		c := &conv.FileStat{Id: "c", Parent: b}
		b.Children = append(b.Children, c)

		d := &conv.FileStat{Id: "d", Parent: root}
		root.Children = append(root.Children, d)

		e := &conv.FileStat{Id: "c", Parent: d}
		d.Children = append(d.Children, e)

		// When
		rootFolder := toFolder(root, 8)
		printer := NewColumnPrinter([]int{20, 20})
		printFolderTree(printer, rootFolder, 0)

		// Then
		assert.Equal(t, 2, len(rootFolder.children))
		assert.Equal(t, "a", rootFolder.children[0].docId)
		assert.Equal(t, "d", rootFolder.children[1].docId)

		assert.Equal(t, "b", rootFolder.children[0].children[0].docId)
	})
}

func TestSortFolders(t *testing.T) {
	t.Run("Should indent text right", func(t *testing.T) {
		rootFile := &conv.FileStat{Id: "root"}

		a := &conv.FileStat{Id: "a", Parent: rootFile, ViewCount: 2}
		b := &conv.FileStat{Id: "b", Parent: rootFile, ViewCount: 3}
		c := &conv.FileStat{Id: "c", Parent: rootFile, ViewCount: 1}

		rootFile.Children = append(rootFile.Children, a, b, c)

		f := newFolderFromFile(rootFile)
		af := newFolderFromFile(a)
		bf := newFolderFromFile(b)
		cf := newFolderFromFile(c)

		f.AddChild(af)
		f.AddChild(bf)
		f.AddChild(cf)

		assert.Equal(t, "a", f.children[0].docId)
		assert.Equal(t, "b", f.children[1].docId)
		assert.Equal(t, "c", f.children[2].docId)

		// When
		sortFoldersByViews(f)
		fmt.Println(f.children)

		// Then
		assert.Equal(t, "b", f.children[0].docId)
		assert.Equal(t, "a", f.children[1].docId)
		assert.Equal(t, "c", f.children[2].docId)
	})
}

func TestRigtIndent(t *testing.T) {
	t.Run("Should indent text right", func(t *testing.T) {
		assert.Equal(t, " 123", rightIndent(4, "123"))
	})
}
