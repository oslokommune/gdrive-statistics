package get_gdrive_views

import (
	"fmt"
	"github.com/oslokommune/gdrive-statistics/hasher"
	"time"
)

type GdriveViewEvent struct {
	time     *time.Time
	userHash *hasher.Hash
	docId    string
	docTitle string
}

func (g GdriveViewEvent) String() string {
	return fmt.Sprintf("VIEW [%s] [user hash %s] [docId %s] [docTitle %s]",
		g.time.Format(time.RFC822),
		g.userHash,
		g.docId,
		g.docTitle,
	)
}
