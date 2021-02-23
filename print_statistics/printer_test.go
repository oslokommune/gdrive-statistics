package print_statistics

import (
	conv "github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFileStats(t *testing.T) {
	t.Run("Should find two folders", func(t *testing.T) {
		root := &conv.FileStat{Id: "root"}

		a := conv.FileStat{Id: "a", Parent: root}
		root.Children = append(root.Children, &a)

		b := &conv.FileStat{Id: "b", Parent: &a}
		a.Children = append(a.Children, b)

		c := &conv.FileStat{Id: "c", Parent: b}
		b.Children = append(b.Children, c)

		d := &conv.FileStat{Id: "d", Parent: root}
		root.Children = append(root.Children, d)

		e := &conv.FileStat{Id: "c", Parent: d}
		d.Children = append(d.Children, e)

		// When
		rootFolder := toFolder(root, 8)
		printFolderTree(rootFolder, 0)

		// Then
		assert.Equal(t, 2, len(rootFolder.children))
		assert.Equal(t, "a", rootFolder.children[0].docId)
		assert.Equal(t, "d", rootFolder.children[1].docId)

		assert.Equal(t, "b", rootFolder.children[0].children[0].docId)

	})
}
