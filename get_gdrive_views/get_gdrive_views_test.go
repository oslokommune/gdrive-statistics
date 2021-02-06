package get_gdrive_views_test

import (
	"testing"
)

func TestWhatever(t *testing.T) {
	testCases := []struct {
		name string
	}{
		{
			name: "Should work",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
		})
	}
}
