package print_statistics

import "fmt"

type Folder struct {
	docId    string
	parent   *Folder
	children []*Folder

	viewCount       int
	uniqueViewcount int
}

func (f *Folder) String() string {
	return fmt.Sprintf("[DocId: %s]", f.docId)
}

func newFolder(docId string, parent *Folder, children []*Folder, views int, uniqueViews int) *Folder {
	var actualChildren []*Folder
	if children == nil {
		actualChildren = make([]*Folder, 0)
	} else {
		actualChildren = children
	}

	return &Folder{
		docId:           docId,
		parent:          parent,
		children:        actualChildren,
		viewCount:       views,
		uniqueViewcount: uniqueViews,
	}
}

func (f *Folder) AddChild(child *Folder) {
	f.children = append(f.children, child)
}
