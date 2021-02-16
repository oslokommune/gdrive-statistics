package view_count_calculator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetChildren(t *testing.T) {
	t.Run("should set child relationships correctly", func(t *testing.T) {
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

		c := &FileStat{
			Id:        "c",
			ViewCount: 2,
			Parent:    a,
		}

		fileStats := map[string]*FileStat{
			"root": root,
			"a":    a,
			"b":    b,
			"c":    c,
		}

		SetChildren(fileStats)

		assert.Contains(t, root.Children, a)
		assert.Contains(t, root.Children, b)
		assert.Len(t, root.Children, 2)

		assert.Contains(t, a.Children, c)
		assert.Len(t, a.Children, 1)

		assert.Len(t, b.Children, 0)
	})
}
