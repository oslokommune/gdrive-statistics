package get_gdrive_views

import (
	"fmt"
	"time"

	"github.com/oslokommune/gdrive-statistics/hasher"
)

type GdriveViewEvent struct {
	time     *time.Time
	userHash *hasher.Hash
	docId    string
	docTitle string
}

func (g GdriveViewEvent) String() string {
	return fmt.Sprintf("View [%s] [docId %s] [docTitle %s]",
		g.time.Format(time.RFC822),
		g.docId,
		g.docTitle,
	)
}
