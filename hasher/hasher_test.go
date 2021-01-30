package hasher_test

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/hasher"
	"testing"
)

func TestWhatever(t *testing.T) {
	t.Run("Test something", func(t *testing.T) {
		hash := hasher.NewHash("test")
		fmt.Printf("Hash: %s\n", hash)
	})
}
