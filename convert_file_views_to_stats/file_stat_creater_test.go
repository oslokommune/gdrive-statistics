package convert_file_views_to_stats_test

import (
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	t.Run("should build correct tree structure", func(t *testing.T) {
		files := []*get_file_list.FileOrFolder{
			{
				Id:     "root",
				Name:   "DummyRoot",
				Parent: "",
			},
			{
				Id:     "a",
				Name:   "a.txt",
				Parent: "root",
			},
			{
				Id:     "d1",
				Name:   "DIR1",
				Parent: "root",
			},
			{
				Id:     "b",
				Name:   "b.txt",
				Parent: "d1",
			},
		}

		// TODO next time: Put in views

		fileStats := convert_file_views_to_stats.CreateFileStats(files)

		// TODO make some smart asserts

		//root := &FileStat{
		//	Id:        "root",
		//	ViewCount: 0,
		//	Parent:    nil,
		//	Children:  nil,
		//}
		//
		//assert.Equal(t, "root", fileStats["root"].Id)
		//assert.Equal(t, 0, fileStats["root"].ViewCount)

	})
}
