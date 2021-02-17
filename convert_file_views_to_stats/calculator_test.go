package convert_file_views_to_stats

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcDirViewCount(t *testing.T) {
	t.Run("should set root view count to sum of two children", func(t *testing.T) {
		root := &FileStat{
			Id:        "root",
			ViewCount: 0,
			Parent:    nil,
		}

		a := &FileStat{
			Id:        "a",
			ViewCount: 1,
			Parent:    root,
		}

		b := &FileStat{
			Id:        "b",
			ViewCount: 2,
			Parent:    root,
		}

		fileStats := map[string]*FileStat{
			"root": root,
			"a":    a,
			"b":    b,
		}

		root.Children = []*FileStat{a, b}
		aggregateViews(root)

		assert.Equal(t, 3, fileStats["root"].ViewCount)
	})

	t.Run("should aggregate tree structure", func(t *testing.T) {
		root := &FileStat{
			Id:        "root",
			ViewCount: 0,
			Parent:    nil,
		}

		a := &FileStat{
			Id:        "a.txt",
			ViewCount: 1,
			Parent:    root,
		}

		b := &FileStat{
			Id:        "b.txt",
			ViewCount: 2,
			Parent:    root,
		}

		dir1 := &FileStat{
			Id:        "DIR1",
			ViewCount: 0,
			Parent:    root,
		}

		root.Children = []*FileStat{a, b, dir1}

		c := &FileStat{
			Id:        "c.txt",
			ViewCount: 3,
			Parent:    dir1,
		}

		dir2 := &FileStat{
			Id:        "D2",
			ViewCount: 0,
			Parent:    dir1,
		}

		dir1.Children = []*FileStat{c, dir2}

		d := &FileStat{
			Id:        "d.txt",
			ViewCount: 4,
			Parent:    dir2,
		}

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

		aggregateViews(root)

		assert.Equal(t, 10, fileStats["root"].ViewCount)
	})
}
