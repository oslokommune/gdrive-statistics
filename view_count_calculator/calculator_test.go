package view_count_calculator

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

		SetChildren(fileStats)
		AggregateViews(root)

		assert.Equal(t, 3, fileStats["root"].ViewCount)
	})

	t.Run("should aggregate tree structure", func(t *testing.T) {
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

		dir1 := &FileStat{
			Id:        "D1",
			ViewCount: 0,
			Parent:    root,
		}

		c := &FileStat{
			Id:        "c",
			ViewCount: 3,
			Parent:    dir1,
		}

		dir2 := &FileStat{
			Id:        "D2",
			ViewCount: 0,
			Parent:    dir1,
		}

		d := &FileStat{
			Id:        "d",
			ViewCount: 4,
			Parent:    dir2,
		}

		fileStats := map[string]*FileStat{
			"root": root,
			"a":    a,
			"b":    b,
			"D1":   dir1,
			"c":    c,
			"D2":   dir2,
			"d":    d,
		}

		SetChildren(fileStats)
		AggregateViews(root)

		assert.Equal(t, 10, fileStats["root"].ViewCount)
	})
}
