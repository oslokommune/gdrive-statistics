package convert_file_views_to_stats

import (
	"github.com/oslokommune/gdrive-statistics/get_api_data/get_file_list"
	"github.com/oslokommune/gdrive-statistics/get_api_data/get_gdrive_views"
	"github.com/oslokommune/gdrive-statistics/hasher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateFileStats(t *testing.T) {
	t.Run("should create file stats", func(t *testing.T) {
		// Given
		/* Directory structure:
		   root/					V:11,UQ: 7
		     a.txt joe				V:1, UQ: 1
		     d1/					V:6, UQ: 3
		       b.txt joe, bob		V:5, UQ: 2
		       c.txt bob			V:1, UQ: 1
		     d2/					V:3, UQ: 3
		       d.txt bob			V:1, UQ: 1
		       d3/					V:2, UQ: 2
		         e.txt joe, bob		V:2, UQ: 2
		*/
		rootId := "root"
		files := []*get_file_list.FileOrFolder{
			{Id: "a", Name: "a.txt", Parent: rootId},
			{Id: "d1", Name: "DIR1", Parent: rootId},
			{Id: "b", Name: "b.txt", Parent: "d1"},
			{Id: "c", Name: "c.txt", Parent: "d1"},
			{Id: "d2", Name: "DIR2", Parent: rootId},
			{Id: "d", Name: "d.txt", Parent: "d2"},
			{Id: "d3", Name: "DIR3", Parent: "d2"},
			{Id: "e", Name: "e.txt", Parent: "d3"},
		}

		views := []*get_gdrive_views.GdriveViewEvent{
			{DocId: "a", UserHash: hasher.NewHash("joe"), Time: nil},

			{DocId: "b", UserHash: hasher.NewHash("joe"), Time: nil},
			{DocId: "b", UserHash: hasher.NewHash("joe"), Time: nil},
			{DocId: "b", UserHash: hasher.NewHash("joe"), Time: nil},
			{DocId: "b", UserHash: hasher.NewHash("bob"), Time: nil},
			{DocId: "b", UserHash: hasher.NewHash("bob"), Time: nil},

			{DocId: "c", UserHash: hasher.NewHash("bob"), Time: nil},

			{DocId: "d", UserHash: hasher.NewHash("bob"), Time: nil},

			{DocId: "e", UserHash: hasher.NewHash("joe"), Time: nil},
			{DocId: "e", UserHash: hasher.NewHash("bob"), Time: nil},
		}

		// When
		fileStats := CreateFileStats(rootId, files, views)

		// Then verify views of individual files
		assert.Equal(t, 1, fileStats["a"].ViewCount, "a")
		assert.Equal(t, 5, fileStats["b"].ViewCount, "b")
		assert.Equal(t, 1, fileStats["c"].ViewCount, "c")
		assert.Equal(t, 1, fileStats["d"].ViewCount, "d")
		assert.Equal(t, 2, fileStats["e"].ViewCount, "e")

		// Verify aggregated use of folders
		assert.Equal(t, 2, fileStats["d3"].ViewCount, "d3")
		assert.Equal(t, 3, fileStats["d2"].ViewCount, "d2")
		assert.Equal(t, 6, fileStats["d1"].ViewCount, "d1")
		assert.Equal(t, 10, fileStats[rootId].ViewCount, "root")

		// Verify unique views for individual files
		assert.Equal(t, 1, fileStats["a"].UniqueViewCount, "a")
		assert.Equal(t, 2, fileStats["b"].UniqueViewCount, "b")
		assert.Equal(t, 1, fileStats["c"].UniqueViewCount, "c")
		assert.Equal(t, 1, fileStats["d"].UniqueViewCount, "d")
		assert.Equal(t, 2, fileStats["e"].UniqueViewCount, "e")

		// Verify unique views for folders
		assert.Equal(t, 2, fileStats["d3"].UniqueViewCount)
		assert.Equal(t, 3, fileStats["d2"].UniqueViewCount)
		assert.Equal(t, 3, fileStats["d1"].UniqueViewCount)
		assert.Equal(t, 7, fileStats[rootId].UniqueViewCount)
	})
}

func TestSetParentsAndChildren(t *testing.T) {
	t.Run("should build correct tree structure from filestats", func(t *testing.T) {
		// Given
		rootId := "myroot"
		files := []*get_file_list.FileOrFolder{
			{Id: "a", Name: "a.txt", Parent: rootId},
			{Id: "d1", Name: "DIR1", Parent: rootId},
			{Id: "b", Name: "b.txt", Parent: "d1"},
			{Id: "c", Name: "c.txt", Parent: "d1"},
			{Id: "d2", Name: "DIR2", Parent: "d1"},
			{Id: "e", Name: "e.txt", Parent: "d2"},
		}

		fileStats := make(map[string]*FileStat)

		fileStats["a"] = &FileStat{Id: "a", ViewCount: 0, Parent: nil, Children: nil}
		fileStats["d1"] = &FileStat{Id: "c", ViewCount: 0, Parent: nil, Children: nil}
		fileStats["b"] = &FileStat{Id: "b", ViewCount: 0, Parent: nil, Children: nil}
		fileStats["c"] = &FileStat{Id: "c", ViewCount: 3, Parent: nil, Children: nil}
		fileStats["d2"] = &FileStat{Id: "c", ViewCount: 0, Parent: nil, Children: nil}
		fileStats["e"] = &FileStat{Id: "c", ViewCount: 0, Parent: nil, Children: nil}

		// When
		setParentsAndChildren(files, fileStats, rootId)
		//fileStats := convert_file_views_to_stats.CreateFileStats("root", files, views)

		// Then
		assert.Equal(t, rootId, fileStats["a"].Parent.Id)
		assert.Equal(t, rootId, fileStats["d1"].Parent.Id)

		assert.Equal(t, fileStats["d1"], fileStats["b"].Parent)
		assert.Equal(t, fileStats["d1"], fileStats["c"].Parent)
		assert.Equal(t, fileStats["d1"], fileStats["d2"].Parent)
		assert.Equal(t, fileStats["d2"], fileStats["e"].Parent)

		assert.Equal(t, 3, len(fileStats["d1"].Children))

		assert.Contains(t, fileStats["d1"].Children, fileStats["b"])
		assert.Contains(t, fileStats["d1"].Children, fileStats["c"])
		assert.Contains(t, fileStats["d1"].Children, fileStats["d2"])

		assert.Equal(t, 1, len(fileStats["d2"].Children))
		assert.Contains(t, fileStats["d2"].Children, fileStats["e"])
	})

	t.Run("should build correct tree structure from files", func(t *testing.T) {
		rootId := "myroot"
		a := &get_file_list.FileOrFolder{Id: "a", Parent: rootId}
		b := &get_file_list.FileOrFolder{Id: "b", Parent: rootId}
		c := &get_file_list.FileOrFolder{Id: "c", Parent: "a"}

		files := []*get_file_list.FileOrFolder{a, b, c}
		var views []*get_gdrive_views.GdriveViewEvent

		fileStats := toFileStats(rootId, files, views)

		fsRoot := fileStats[rootId]

		assert.Contains(t, fsRoot.Children, fileStats["a"])
		assert.Contains(t, fsRoot.Children, fileStats["b"])
		assert.Len(t, fsRoot.Children, 2)

		assert.Contains(t, fileStats["a"].Children, fileStats["c"])
		assert.Len(t, fileStats["a"].Children, 1)

		assert.Len(t, fileStats["b"].Children, 0)
	})
}

func TestViewCount(t *testing.T) {
	t.Run("should set root view count to sum of its children", func(t *testing.T) {
		root := &FileStat{Id: "root", ViewCount: 0, Parent: nil}
		a := &FileStat{Id: "a", ViewCount: 1, Parent: root}
		b := &FileStat{Id: "b", ViewCount: 2, Parent: root}
		fileStats := map[string]*FileStat{"root": root, "a": a, "b": b}

		root.Children = []*FileStat{a, b}
		aggregateViews(root)

		assert.Equal(t, 3, fileStats["root"].ViewCount)
	})

	t.Run("should sum view counts for files and folders", func(t *testing.T) {
		// Given
		root := &FileStat{Id: "root", ViewCount: 0, Parent: nil}

		a := &FileStat{Id: "a.txt", ViewCount: 1, Parent: root}
		b := &FileStat{Id: "b.txt", ViewCount: 2, Parent: root}
		dir1 := &FileStat{Id: "DIR1", ViewCount: 0, Parent: root}

		root.Children = []*FileStat{a, b, dir1}
		c := &FileStat{Id: "c.txt", ViewCount: 3, Parent: dir1}

		dir2 := &FileStat{Id: "D2", ViewCount: 0, Parent: dir1}
		dir1.Children = []*FileStat{c, dir2}

		d := &FileStat{Id: "d.txt", ViewCount: 4, Parent: dir2}
		dir2.Children = []*FileStat{d}

		fileStats := map[string]*FileStat{
			root.Id: root,
			a.Id:    a,
			b.Id:    b,
			"DIR1":  dir1,
			"c":     c,
			"D2":    dir2,
			"d":     d,
		}

		// When
		aggregateViews(root)

		// Then
		assert.Equal(t, 10, fileStats["root"].ViewCount)
	})
}

func TestUniqueViewCount(t *testing.T) {
	t.Run("should calculate unique view counts for file", func(t *testing.T) {
		// Given
		rootId := "root"
		files := []*get_file_list.FileOrFolder{
			{Id: "a", Name: "a.txt", Parent: rootId},
			{Id: "b", Name: "b.txt", Parent: rootId},
			{Id: "c", Name: "c.txt", Parent: rootId},
		}
		fileStats := make(map[string]*FileStat)
		views := []*get_gdrive_views.GdriveViewEvent{
			{DocId: "a", UserHash: hasher.NewHash("joe"), Time: nil},
			{DocId: "b", UserHash: hasher.NewHash("joe"), Time: nil},
			{DocId: "b", UserHash: hasher.NewHash("bob"), Time: nil},
			{DocId: "c", UserHash: hasher.NewHash("bob"), Time: nil},
			{DocId: "c", UserHash: hasher.NewHash("bob"), Time: nil},
		}

		// When
		mergeFilesAndViewsToFileStats(files, fileStats, views)

		// Then
		assert.Equal(t, 1, fileStats["a"].ViewCount)
		assert.Equal(t, 1, fileStats["a"].UniqueViewCount)

		assert.Equal(t, 2, fileStats["b"].ViewCount)
		assert.Equal(t, 2, fileStats["b"].UniqueViewCount)

		assert.Equal(t, 2, fileStats["c"].ViewCount)
		assert.Equal(t, 1, fileStats["c"].UniqueViewCount)
	})

	t.Run("should set root unique view count to sum of its children's unique views", func(t *testing.T) {
		// Given
		root := &FileStat{Id: "root", ViewCount: 0, Parent: nil}
		a := &FileStat{Id: "a", ViewCount: 1, UniqueViewCount: 1, Parent: root}
		b := &FileStat{Id: "b", ViewCount: 2, UniqueViewCount: 2, Parent: root}

		root.Children = []*FileStat{a, b}

		// When
		aggregateViews(root)

		// Then
		assert.Equal(t, 3, root.UniqueViewCount)
	})
}
