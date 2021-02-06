package get_file_list

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/drive/v3"
)

func TestGetFiles(t *testing.T) {
	testCases := []struct {
		name     string
		input    []*drive.File
		expected []*DriveFile
	}{
		{
			name: "Should convert file",
			input: []*drive.File{
				{
					Id:   "a",
					Name: "a.txt",
				}, {
					Id:      "b",
					Name:    "b.txt",
					Parents: []string{"parentB"},
				},
			},
			expected: []*DriveFile{
				{
					Id:     "a",
					Name:   "a.txt",
					Parent: "",
				}, {
					Id:     "b",
					Name:   "b.txt",
					Parent: "parentB",
				},
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			g := New(nil, "")
			driveFiles, err := g.toDriveFile(tc.input)
			assert.NoError(t, err)

			assert.Equal(t, tc.expected, driveFiles)
		})
	}
}
