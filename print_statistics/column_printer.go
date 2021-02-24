package print_statistics

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type ColumnPrinter struct {
	columnSpaces  []int
	currentColumn int
	line          string
}

func NewColumnPrinter(columnSpaces []int) *ColumnPrinter {
	return &ColumnPrinter{
		columnSpaces: append(columnSpaces, 0),
	}
}

func (c *ColumnPrinter) reset() {
	c.line = ""
	c.currentColumn = 0
}

func (c *ColumnPrinter) add(txt string) {
	c.line += txt

	if c.currentColumn >= len(c.columnSpaces) {
		fmt.Println("WARN: adding text beyound defined columns")
		return
	}

	columnSpaces := c.columnSpaces[c.currentColumn]
	textLength := utf8.RuneCountInString(txt)

	spaceCount := columnSpaces - textLength
	if c.currentColumn == 0 {
		// As the first text is 1-indexed, we have to remove 1 space. If we don't, we'll get this failed test:
		//            1234567890123456789012
		// expected: "hello    there     you"
		// actual  : "hello     there     you"
		spaceCount--
	}

	if spaceCount < 0 {
		spaceCount = 0
	}

	c.line += strings.Repeat(" ", spaceCount)

	c.currentColumn++

	//c.previousText = txt
	//c.currentCoord += len(txt)
}

func (c *ColumnPrinter) get() string {
	return c.line
}
