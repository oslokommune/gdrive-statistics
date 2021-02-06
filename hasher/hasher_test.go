package hasher_test

import (
	"fmt"
	"testing"

	"github.com/oslokommune/gdrive-statistics/hasher"
)

func TestWhatever(t *testing.T) {
	t.Run("Test something", func(t *testing.T) {
		hash := hasher.NewHash("test")
		fmt.Printf("Hash: %s\n", hash)
	})
}
