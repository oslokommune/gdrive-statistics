package print_statistics

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"unicode/utf8"
)

func TestColumnPrinter(t *testing.T) {
	t.Run("Should print columns", func(t *testing.T) {
		c := NewColumnPrinter([]int{10, 10})

		// When
		c.add("hello")
		c.add("there")
		c.add("you")

		// Then
		assert.Equal(t, "hello    there     you", c.get())
	})

	t.Run("Custom test", func(t *testing.T) {

		c := NewColumnPrinter([]int{45, 10})

		c.add("FOLDER")
		c.add("VIEWS")
		c.add("UNIQUE VIEWS")
		fmt.Println(c.get())

		c.reset()

		// When
		c.add("  Maler og verktøy")
		c.add("73")
		c.add("342")
		fmt.Println(c.get())
		c.reset()

		c.add("  Maler og verktoy")
		c.add("73")
		c.add("342")
		fmt.Println(c.get())
		c.reset()

		fmt.Println(utf8.RuneCountInString("Ø"))
	})

}
