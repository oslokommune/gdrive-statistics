package convert_file_views_to_stats

import "github.com/oslokommune/gdrive-statistics/hasher"

type UniqueViewCalculator struct {
	docToUsers map[string]map[string]bool
}

func NewUniqueViewCalculator() *UniqueViewCalculator {
	return &UniqueViewCalculator{
		docToUsers: make(map[string]map[string]bool),
	}
}

func (c *UniqueViewCalculator) addViewForDocument(docId string, userHash *hasher.Hash) {
	if _, ok := c.docToUsers[docId]; !ok {
		c.docToUsers[docId] = make(map[string]bool)
	}

	c.docToUsers[docId][userHash.String()] = true
}

func (c *UniqueViewCalculator) getUniqueViewsForDocument(docId string) int {
	return len(c.docToUsers[docId])
}
