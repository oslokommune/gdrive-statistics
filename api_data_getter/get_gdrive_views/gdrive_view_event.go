package get_gdrive_views

import (
	"fmt"
	"time"

	"github.com/oslokommune/gdrive-statistics/hasher"
)

type GdriveViewEvent struct {
	DocId    string
	UserHash *hasher.Hash
	Time     *time.Time
	docTitle string
}

func (g GdriveViewEvent) String() string {
	return fmt.Sprintf("View [%s] [UserHash %s] [DocId %s]",
		g.Time.Format(time.RFC822),
		g.UserHash,
		g.DocId,
	)
}
