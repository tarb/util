package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type DynamicList struct {
	BindBuilder func(int) Element
	BindSize    func() int
	BindIndex   *int
	Index       int
	WindowIndex int
	Height      int
	OnChange    func(int)
}

//
func (dl *DynamicList) Draw(x, y int, focus Element) {
	var eX, eY int = x, y
	var sumH int

	for i := dl.WindowIndex; i < dl.BindSize(); i++ {
		var e = dl.BindBuilder(i)

		var listFocus Element
		if i == dl.Index {
			listFocus = e
		}

		var draw func(int, int, Element)
		var size func() (int, int)
		if ex, ok := e.(Expandable); ok && listFocus != nil {
			size = ex.ExpandSize
			draw = ex.ExpandDraw
		} else {
			size = e.Size
			draw = e.Draw
		}

		var eH int
		_, eH = size()

		if sumH+eH > dl.Height {
			break
		}

		// if container is active list item, grab first focusable element inside

		if cont, ok := e.(Container); ok {
			listFocus = cont.NextFocusable(nil)
		}

		draw(eX, eY, listFocus)
		eY, sumH = eY+eH, sumH+eH
	}
}

//
func (dl *DynamicList) Size() (int, int) {
	var maxX, sumY int

	for i := dl.WindowIndex; i < dl.BindSize(); i++ {
		var e = dl.BindBuilder(i)

		var eW, eH int
		if ex, ok := e.(Expandable); ok && i == dl.Index {
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

	return maxX, dl.Height
}

//
func (dl *DynamicList) Handle(ev termbox.Event) {
	if ev.Key == termbox.KeyArrowUp {
		dl.scrollUp()
	} else if ev.Key == termbox.KeyArrowDown {
		dl.scrollDown()
	} else {
		// pass event on to the next focusable thing in the item
		var e = dl.BindBuilder(dl.Index)
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
func (dl *DynamicList) HandleClick(ev termbox.Event) {
	// fmt.Println("list", ev.MouseX, ev.MouseY)

	switch ev.Key {
	case termbox.MouseWheelDown:
		if dl.WindowIndex+dl.visibleItems() < dl.BindSize() {
			dl.WindowIndex++
		}

	case termbox.MouseWheelUp:
		if dl.WindowIndex > 0 {
			dl.WindowIndex--
		}

	case termbox.MouseLeft:
		var sumY int

		for i := dl.WindowIndex; i < dl.BindSize(); i++ {
			var e = dl.BindBuilder(i)

			var cw, ch int
			if ex, ok := e.(Expandable); ok && i == dl.Index {
				cw, ch = ex.ExpandSize()
			} else {
				cw, ch = e.Size()
			}

			if ev.MouseX >= 0 && ev.MouseY >= sumY && ev.MouseX < cw && ev.MouseY < sumY+ch {

				// adjust the event before passing it down
				ev.MouseY -= sumY

				if clickable, ok := e.(Clickable); ok {
					clickable.HandleClick(ev)
				}
				if cont, ok := e.(Container); ok {
					cont.FocusClicked(ev)
				}

				// fire events and update bindings
				if dl.OnChange != nil && dl.Index != i { // fire registered onChange
					dl.OnChange(i)
				}
				if dl.BindIndex != nil { // update bound index
					*dl.BindIndex = i
				}
				dl.Index = i

				// scroll the WindowIndex clicked on top|bottom element (if possible)
				if i == dl.WindowIndex && dl.WindowIndex > 0 {
					dl.WindowIndex--
				} else if dl.Index > dl.WindowIndex+dl.visibleItems()-2 && dl.WindowIndex+dl.visibleItems() < dl.BindSize() {
					dl.WindowIndex++
				}

				return
			}
			sumY += ch
		}

	}

}

//
func (dl *DynamicList) scrollDown() {
	var lastIndex = dl.BindSize() - 1

	if dl.Index < lastIndex {
		dl.Index++
		if dl.BindIndex != nil {
			*dl.BindIndex = dl.Index
		}
		if dl.OnChange != nil { // fire registered onChange
			dl.OnChange(dl.Index)
		}

		var idx = dl.Index
		if idx < lastIndex {
			idx++
		}

		var sumY int
		for i := idx; i >= 0; i-- {
			var e = dl.BindBuilder(i)

			var eh int
			if ex, ok := e.(Expandable); ok && i == dl.Index {
				_, eh = ex.ExpandSize()
			} else {
				_, eh = e.Size()
			}
			sumY += eh

			if sumY >= dl.Height {
				dl.WindowIndex = i + 1
				break
			}
			if i == dl.WindowIndex {
				break
			}
		}
	}
}

//
func (dl *DynamicList) scrollUp() {
	if dl.Index > 0 {
		dl.Index--
		if dl.BindIndex != nil {
			*dl.BindIndex = dl.Index
		}
		if dl.OnChange != nil { // fire registered onChange
			dl.OnChange(dl.Index)
		}

		if dl.Index == dl.WindowIndex && dl.Index != 0 {
			dl.WindowIndex--
		}
	}
}

//
func (dl *DynamicList) visibleItems() int {
	var sumY, count int

	for i := dl.WindowIndex; i < dl.BindSize(); i++ {
		var e = dl.BindBuilder(i)

		var _, ch int
		if ex, ok := e.(Expandable); ok && i == dl.Index {
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
