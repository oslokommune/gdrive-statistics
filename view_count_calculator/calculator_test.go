package view_count_calculator

import (
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCalcDirViewCount(t *testing.T) {
	t.Run("should set root view count to sum of two children", func(t *testing.T) {
		root := &convert_file_views_to_stats.FileStat{
			Id:        "root",
			ViewCount: 0,
			Parent:    nil,
		}

		a := &convert_file_views_to_stats.FileStat{
			Id:        "a",
			ViewCount: 1,
			Parent:    root,
		}

		b := &convert_file_views_to_stats.FileStat{
			Id:        "b",
			ViewCount: 2,
			Parent:    root,
		}

		fileStats := map[string]*convert_file_views_to_stats.FileStat{
			"root": root,
			"a":    a,
			"b":    b,
		}

		convert_file_views_to_stats.SetChildren(fileStats)
		AggregateViews(root)

		assert.Equal(t, 3, fileStats["root"].ViewCount)
	})

	t.Run("should aggregate tree structure", func(t *testing.T) {
		root := &convert_file_views_to_stats.FileStat{
			Id:        "root",
			ViewCount: 0,
			Parent:    nil,
		}

		a := &convert_file_views_to_stats.FileStat{
			Id:        "a",
			ViewCount: 1,
			Parent:    root,
		}

		b := &convert_file_views_to_stats.FileStat{
			Id:        "b",
			ViewCount: 2,
			Parent:    root,
		}

		dir1 := &convert_file_views_to_stats.FileStat{
			Id:        "D1",
			ViewCount: 0,
			Parent:    root,
		}

		c := &convert_file_views_to_stats.FileStat{
			Id:        "c",
			ViewCount: 3,
			Parent:    dir1,
		}

		dir2 := &convert_file_views_to_stats.FileStat{
			Id:        "D2",
			ViewCount: 0,
			Parent:    dir1,
		}

		d := &convert_file_views_to_stats.FileStat{
			Id:        "d",
			ViewCount: 4,
			Parent:    dir2,
		}

		fileStats := map[string]*convert_file_views_to_stats.FileStat{
			"root": root,
			"a":    a,
			"b":    b,
			"D1":   dir1,
			"c":    c,
			"D2":   dir2,
			"d":    d,
		}

		convert_file_views_to_stats.SetChildren(fileStats) // TODO this should not be necessary, do this earlier
		AggregateViews(root)

		assert.Equal(t, 10, fileStats["root"].ViewCount)
	})
}
