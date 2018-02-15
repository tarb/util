package tbui

import termbox "github.com/nsf/termbox-go"

//
type If struct {
	Child Element
	Cond  func() bool
}

//
func (ifl *If) Draw(x, y int, focused Element) {
	if ifl.Cond() {
		ifl.Child.Draw(x, y, focused)
	}
}

//
func (ifl *If) Size() (int, int) {
	if ifl.Cond() {
		return ifl.Child.Size()
	}
	return 0, 0
}

//
func (ifl *If) GetFocusable() []Focusable {
	if cont, ok := ifl.Child.(Container); ok {
		return cont.GetFocusable()
	} else if focusable, ok := ifl.Child.(Focusable); ok {
		return []Focusable{focusable}
	}

	return []Focusable{}
}

//
func (ifl *If) NextFocusable(current Focusable) Focusable {
	var eles []Focusable = ifl.GetFocusable()

	// if there are focusable
	if len(eles) > 0 {
		//find the next focusable
		if current != nil {
			var curIdx int

			for i, e := range eles {
				if e == current {
					curIdx = i
				}
			}

			if curIdx < len(eles)-1 {
				return eles[curIdx+1]
			}
			return eles[0]

		}
		// if nothing is currently focused, return the first
		return eles[0]
	}
	// if there are no focusable, return nil
	return nil
}

//
func (ifl *If) FocusClicked(ev termbox.Event) Focusable {
	if ifl.Cond() {
		if clickable, ok := ifl.Child.(Clickable); ok {
			clickable.HandleClick(ev)
		}
		if cont, ok := ifl.Child.(Container); ok {
			return cont.FocusClicked(ev)
		} else if foc, ok := ifl.Child.(Focusable); ok {
			return foc
		}
	}

	return nil
}
