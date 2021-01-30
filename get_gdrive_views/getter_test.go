package get_gdrive_views

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/hasher"
	"testing"
	"time"
)

func TestWhatever(t *testing.T) {
	t.Run("Test something", func(t *testing.T) {
		now := time.Now()

		view := GdriveViewEvent{
			time:     &now,
			userHash: hasher.NewHash("HELLO"),
			docId:    "SOMEDOCID",
		}

		fmt.Println(view)
	})
}
