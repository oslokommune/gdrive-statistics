package convert_file_views_to_stats

import (
	"github.com/oslokommune/gdrive-statistics/get_api_data/get_file_list"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetChildren(t *testing.T) {
	t.Run("should set child relationships correctly", func(t *testing.T) {
		root := &get_file_list.FileOrFolder{
			Id:     "root",
			Parent: "",
		}

		a := &get_file_list.FileOrFolder{
			Id:     "a",
			Parent: "root",
		}

		b := &get_file_list.FileOrFolder{
			Id:     "b",
			Parent: "root",
		}

		c := &get_file_list.FileOrFolder{
			Id:     "c",
			Parent: "a",
		}

		files := []*get_file_list.FileOrFolder{root, a, b, c}

		fileStats := toFileStats("root", files)

		fsRoot := fileStats["root"]
		fsA := fileStats["a"]
		fsB := fileStats["a"]

		assert.Contains(t, fsRoot.Children, a)
		assert.Contains(t, fsRoot.Children, b)
		assert.Len(t, fsRoot.Children, 2)

		assert.Contains(t, fsA.Children, c)
		assert.Len(t, fsA.Children, 1)

		assert.Len(t, fsB.Children, 0)
	})
}
