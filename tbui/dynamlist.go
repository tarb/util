package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type DynamicList struct {
	Bind     func(int) Element
	BindSize func() int
	Height   int
	Padding  Padding

	selectedIdx int
	windowIdx   int
}

//
func (dl *DynamicList) Draw(x, y int, focus Element) {
	var eX, eY int = x + dl.Padding.Left(), y + dl.Padding.Up()
	var sumH int

	for i := dl.windowIdx; i < dl.BindSize(); i++ {
		var e = dl.Bind(i)
		var eH int

		if ex, ok := e.(Expandable); ok && e == dl.Bind(dl.selectedIdx) {
			_, eH = ex.ExpandSize()
		} else {
			_, eH = e.Size()
		}

		if sumH+eH >= dl.Height {
			break
		}

		// if container is active list item, grab first focusable element inside
		var listFocus Element = dl.Bind(dl.selectedIdx)
		if cont, ok := listFocus.(Container); ok {
			listFocus = cont.NextFocusable(nil)
		}

		e.Draw(eX, eY, listFocus)
		eY, sumH = eY+eH, sumH+eH
	}
}

//
func (dl *DynamicList) Size() (int, int) {
	var maxX, sumY int

	for i := dl.windowIdx; i < dl.BindSize(); i++ {
		var e = dl.Bind(i)

		var eW, eH int
		if ex, ok := e.(Expandable); ok && e == dl.Bind(dl.selectedIdx) {
			eW, eH = ex.ExpandSize()
		} else {
			eW, eH = e.Size()
		}

		if eW > maxX {
			maxX = eW
		}
		sumY += eH
		if sumY > dl.Height {
			break
		}
	}

	maxX += (dl.Padding.Left() + dl.Padding.Right())

	return maxX, dl.Padding.Up() + dl.Height + dl.Padding.Down()
}

//
func (dl *DynamicList) Handle(ev termbox.Event) {
	if ev.Key == termbox.KeyArrowUp {
		dl.scrollUp()
	} else if ev.Key == termbox.KeyArrowDown {
		dl.scrollDown()
	} else {
		// pass event on to the next focusable thing in the item
		var e = dl.Bind(dl.selectedIdx)
		if cont, ok := e.(Container); ok {
			if foc := cont.NextFocusable(nil); foc != nil {
				foc.Handle(ev)
			}
		} else if foc, ok := e.(Focusable); ok {
			foc.Handle(ev)
		}
	}
}

//
func (dl *DynamicList) HandleClick(mouseX, mouseY int) {
	// fmt.Println("list", mouseX, mouseY)
	mouseX, mouseY = mouseX-dl.Padding.Left(), mouseY-dl.Padding.Up()

	var sumY int

	for i := dl.windowIdx; i < dl.BindSize(); i++ {
		var e = dl.Bind(i)

		var cw, ch int
		if ex, ok := e.(Expandable); ok && i == dl.selectedIdx {
			cw, ch = ex.ExpandSize()
		} else {
			cw, ch = e.Size()
		}

		if mouseX >= 0 && mouseY >= sumY && mouseX < cw && mouseY < sumY+ch {
			if clickable, ok := e.(Clickable); ok {
				clickable.HandleClick(mouseX, mouseY-sumY)
			}
			if cont, ok := e.(Container); ok {
				cont.FocusClicked(mouseX, mouseY-sumY)
			}
			dl.selectedIdx = dl.windowIdx + i

			// scroll the windowIdx clicked on top|bottom element (if possible)
			if i == 0 && dl.windowIdx > 0 {
				dl.windowIdx--
			} else if dl.selectedIdx > dl.windowIdx+dl.visibleItems()-2 && dl.windowIdx+dl.visibleItems() < dl.BindSize() {
				dl.windowIdx++
			}

			return
		}

		sumY += ch
	}
}

//
func (dl *DynamicList) scrollDown() {
	var lastIndex = dl.BindSize() - 1

	if dl.selectedIdx < lastIndex {
		dl.selectedIdx++

		var idx = dl.selectedIdx
		if idx < lastIndex {
			idx++
		}

		var sumY int
		for i := idx; i >= 0; i-- {
			var e = dl.Bind(i)

			var eh int
			if ex, ok := e.(Expandable); ok && i == dl.selectedIdx {
				_, eh = ex.ExpandSize()
			} else {
				_, eh = e.Size()
			}
			sumY += eh

			if sumY >= dl.Height {
				dl.windowIdx = i + 1
				break
			}
			if i == dl.windowIdx {
				break
			}
		}
	}
}

//
func (dl *DynamicList) scrollUp() {
	if dl.selectedIdx > 0 {
		dl.selectedIdx--

		if dl.selectedIdx == dl.windowIdx && dl.selectedIdx != 0 {
			dl.windowIdx--
		}
	}
}

//
func (dl *DynamicList) visibleItems() int {
	var sumY, count int

	for i := dl.windowIdx; i < dl.BindSize(); i++ {
		var e = dl.Bind(i)

		var _, ch int
		if ex, ok := e.(Expandable); ok && i == dl.selectedIdx {
			_, ch = ex.ExpandSize()
		} else {
			_, ch = e.Size()
		}

		sumY += ch
		if sumY > dl.Height {
			break
		}
		count++
	}

	return count
}
