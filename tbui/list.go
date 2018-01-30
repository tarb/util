package tbui

import (
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

		// if container is active list item, grab first focusable element inside
		var listFocus Element = l.Items[l.selectedIdx]
		if cont, ok := listFocus.(Container); ok {
			listFocus = cont.NextFocusable(nil)
		}

		e.Draw(eX, eY, listFocus)
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
	} else {
		// pass event on to the next focusable thing in the item
		if cont, ok := l.Items[l.selectedIdx].(Container); ok {
			if foc := cont.NextFocusable(nil); foc != nil {
				foc.Handle(ev)
			}
		} else if foc, ok := l.Items[l.selectedIdx].(Focusable); ok {
			foc.Handle(ev)
		}
	}
}

//
func (l *List) HandleClick(mouseX, mouseY int) {
	// fmt.Println("list", mouseX, mouseY)
	mouseX, mouseY = mouseX-l.Padding.Left(), mouseY-l.Padding.Up()

	var sumY int
	for i, c := range l.Items[l.windowIdx:] {
		var cw, ch int
		if ex, ok := c.(Expandable); ok && c == l.Items[l.selectedIdx] {
			cw, ch = ex.ExpandSize()
		} else {
			cw, ch = c.Size()
		}

		if mouseX >= 0 && mouseY >= sumY && mouseX < cw && mouseY < sumY+ch {
			if clickable, ok := c.(Clickable); ok {
				clickable.HandleClick(mouseX, mouseY-sumY)
			}
			if cont, ok := c.(Container); ok {
				cont.FocusClicked(mouseX, mouseY-sumY)
			}
			l.selectedIdx = l.windowIdx + i

			// scroll the windowIdx clicked on top|bottom element (if possible)
			if i == 0 && l.windowIdx > 0 {
				l.windowIdx--
			} else if l.selectedIdx > l.windowIdx+l.visibleItems()-2 && l.windowIdx+l.visibleItems() < len(l.Items) {
				l.windowIdx++
			}

			return
		}

		sumY += ch
	}
}

//
func (l *List) scrollDown() {
	if l.selectedIdx < len(l.Items)-1 {
		l.selectedIdx++

		var idx = l.selectedIdx
		if idx < len(l.Items)-1 {
			idx++
		}

		var sumY int
		for i := idx; i >= 0; i-- {
			var eh int
			if ex, ok := l.Items[i].(Expandable); ok && i == l.selectedIdx {
				_, eh = ex.ExpandSize()
			} else {
				_, eh = l.Items[i].Size()
			}
			sumY += eh

			if sumY >= l.Height {
				l.windowIdx = i + 1
				break
			}
			if i == l.windowIdx {
				break
			}
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

//
func (l *List) visibleItems() int {
	var sumY, count int

	for _, c := range l.Items[l.windowIdx:] {
		var _, ch int
		if ex, ok := c.(Expandable); ok && c == l.Items[l.selectedIdx] {
			_, ch = ex.ExpandSize()
		} else {
			_, ch = c.Size()
		}

		sumY += ch
		if sumY > l.Height {
			break
		}
		count++
	}

	return count
}
