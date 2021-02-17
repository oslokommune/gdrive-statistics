package convert_file_views_to_stats_test

import (
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_file_list"
	"github.com/oslokommune/gdrive-statistics/api_data_getter/get_gdrive_views"
	"github.com/oslokommune/gdrive-statistics/convert_file_views_to_stats"
	"github.com/oslokommune/gdrive-statistics/hasher"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJoin(t *testing.T) {
	t.Run("should build correct tree structure", func(t *testing.T) {
		files := []*get_file_list.FileOrFolder{
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
			{
				Id:     "c",
				Name:   "c.txt",
				Parent: "d1",
			},
		}

		views := []*get_gdrive_views.GdriveViewEvent{
			{
				DocId:    "a",
				UserHash: hasher.NewHash("joe"),
				Time:     nil,
			},
			{
				DocId:    "a",
				UserHash: hasher.NewHash("bob"),
				Time:     nil,
			},
			{
				DocId:    "a",
				UserHash: hasher.NewHash("bob"),
				Time:     nil,
			},
			{
				DocId:    "b",
				UserHash: hasher.NewHash("joe"),
				Time:     nil,
			},
			{
				DocId:    "b",
				UserHash: hasher.NewHash("joe"),
				Time:     nil,
			},
			{
				DocId:    "b",
				UserHash: hasher.NewHash("joe"),
				Time:     nil,
			},
			{
				DocId:    "b",
				UserHash: hasher.NewHash("bob"),
				Time:     nil,
			},
			{
				DocId:    "b",
				UserHash: hasher.NewHash("bob"),
				Time:     nil,
			},
			{
				DocId:    "c",
				UserHash: hasher.NewHash("bob"),
				Time:     nil,
			},
		}

		fileStats := convert_file_views_to_stats.CreateFileStats("root", files, views)

		// Verify views of individual files
		assert.Equal(t, 3, fileStats["a"].ViewCount)
		assert.Equal(t, 5, fileStats["b"].ViewCount)
		assert.Equal(t, 1, fileStats["c"].ViewCount)

		// Verify aggregated use of folders
		assert.Equal(t, 6, fileStats["d1"].ViewCount)   // 6 = total views of b and c
		assert.Equal(t, 9, fileStats["root"].ViewCount) // 6 = total views of a + b and c

		// Verify unique views for folders
		//assert.Equal(t, 6, fileStats["d1"].UniqueViewCount) // 6 = total views of b and c
	})
}
