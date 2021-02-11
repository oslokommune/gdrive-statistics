package get_gdrive_views

import (
	"fmt"
	"time"

	"github.com/oslokommune/gdrive-statistics/hasher"
)

type GdriveViewEvent struct {
	Time     *time.Time
	userHash *hasher.Hash
	DocId    string
	docTitle string
}

func (g GdriveViewEvent) String() string {
	return fmt.Sprintf("View [%s] [docId %s] [docTitle %s]",
		g.Time.Format(time.RFC822),
		g.DocId,
		g.docTitle,
	)
}
