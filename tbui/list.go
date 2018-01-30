package tbui

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

//
type List struct {
	Items   []Element
	Width   int
	Height  int
	Padding Padding

	selectedIdx int
	windowIdx   int
}

//
func (l *List) Draw(x, y int, focus Element) {
	var eX, eY int = x + l.Padding.Left(), y + l.Padding.Up()
	var sumH int

	for _, e := range l.Items[l.windowIdx:] {
		var eH int
		if ex, ok := e.(Expandable); ok && e == l.Items[l.selectedIdx] {
			_, eH = ex.ExpandSize()
		} else {
			_, eH = e.Size()
		}

		if sumH+eH >= l.Height {
			break
		}

		e.Draw(eX, eY, l.Items[l.selectedIdx])
		eY, sumH = eY+eH, sumH+eH
	}
}

//
func (l *List) Size() (int, int) {
	var maxX int

	for _, e := range l.Items {
		var eW, _ int = e.Size()

		if eW > maxX {
			maxX = eW
		}
	}

	maxX += (l.Padding.Left() + l.Padding.Right())

	return maxX, l.Padding.Up() + l.Height + l.Padding.Down()
}

//
func (l *List) Handle(ev termbox.Event) {
	if ev.Key == termbox.KeyArrowUp {
		l.scrollUp()
	} else if ev.Key == termbox.KeyArrowDown {
		l.scrollDown()
	}
}

//
func (l *List) HandleClick(mouseX, mouseY int) {

}

//
func (l *List) Focusable() bool { return true }

//
func (l *List) scrollDown() {
	var lastIndex int = len(l.Items) - 1

	fmt.Println(l.selectedIdx, "<", lastIndex)
	if l.selectedIdx < lastIndex {
		l.selectedIdx++

		var numVisibleItems int
		var sumY int

		for _, e := range l.Items[l.windowIdx:] {
			var eh int
			if ex, ok := e.(Expandable); ok && e == l.Items[l.selectedIdx] {
				_, eh = ex.ExpandSize()
			} else {
				_, eh = e.Size()
			}

			if sumY+eh >= l.Height {
				break
			}

			sumY += eh
			numVisibleItems++
			// fmt.Println("eh", eh, "sumY", sumY, "numVisibleItems", numVisibleItems)
		}

		// fmt.Println("numVisibleItems", numVisibleItems)

		if l.selectedIdx == l.windowIdx+numVisibleItems-1 && l.windowIdx+numVisibleItems != len(l.Items) {
			l.windowIdx++
		}
	}
}

//
func (l *List) scrollUp() {
	if l.selectedIdx > 0 {
		l.selectedIdx--

		if l.selectedIdx == l.windowIdx && l.selectedIdx != 0 {
			l.windowIdx--
		}
	}
}
