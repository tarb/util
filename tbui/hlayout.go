package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type HLayout struct {
	Children []Element

	MinHeight int
	MinWidth  int
	Border    Border
	Padding   Padding
}

//
func (hl *HLayout) Draw(x, y int, focused Element) {
	// x and y to start drawing children
	var eX, eY int = x + hl.Padding.Left() + hl.Border.Adjust(Left), y + hl.Padding.Up() + hl.Border.Adjust(Up)

	for _, e := range hl.Children {
		e.Draw(eX, eY, focused)
		var eWidth, _ int = e.Size()
		eX += eWidth
	}
	hl.drawBorder(x, y)
}

//
func (hl *HLayout) Size() (int, int) {
	var cumulativeX, maxY int

	for _, e := range hl.Children {
		var w, h int = e.Size()

		cumulativeX += w
		if h > maxY {
			maxY = h
		}
	}

	if cumulativeX < hl.MinWidth {
		cumulativeX = hl.MinWidth
	}
	if maxY < hl.MinHeight {
		maxY = hl.MinHeight
	}

	cumulativeX += hl.Padding.Left() + hl.Padding.Right() + hl.Border.Adjust(Left) + hl.Border.Adjust(Right)
	maxY += hl.Padding.Up() + hl.Padding.Down() + hl.Border.Adjust(Up) + hl.Border.Adjust(Down)

	return cumulativeX, maxY
}

//
// func (hl *HLayout) HandleClick(mouseX, mouseY int) {
// 	fmt.Println("hlayout", hl.Border, "|", mouseX, mouseY, hl.Padding)
// }

func (hl *HLayout) drawBorder(x, y int) {
	var runes []rune = hl.Border.Runes()
	var fg, bg = hl.Border.Fg, hl.Border.Bg
	var w, h int = hl.Size()

	// x
	if hl.Border.Has(Up) {
		for i := x; i < x+w; i++ {
			termbox.SetCell(i, y, runes[0], fg, bg)
		}
	}
	if hl.Border.Has(Down) {
		for i := x; i < x+w; i++ {
			termbox.SetCell(i, y+h-1, runes[0], fg, bg)
		}
	}
	// y
	if hl.Border.Has(Left) {
		for i := y; i < y+h; i++ {
			termbox.SetCell(x, i, runes[1], fg, bg)
		}
	}
	if hl.Border.Has(Right) {
		for i := y; i < y+h; i++ {
			termbox.SetCell(x+w-1, i, runes[1], fg, bg)
		}
	}

	// corners
	if hl.Border.Has(Left | Up) {
		termbox.SetCell(x, y, runes[2], fg, bg)
	}
	if hl.Border.Has(Left | Down) {
		termbox.SetCell(x, y+h-1, runes[4], fg, bg)
	}
	if hl.Border.Has(Right | Up) {
		termbox.SetCell(x+w-1, y, runes[3], fg, bg)
	}
	if hl.Border.Has(Right | Down) {
		termbox.SetCell(x+w-1, y+h-1, runes[5], fg, bg)
	}
}

//
func (hl *HLayout) GetFocusable() []Focusable {
	var eles = make([]Focusable, 0, 10)

	for _, child := range hl.Children {
		if cont, ok := child.(Container); ok {
			eles = append(eles, cont.GetFocusable()...)
		} else if focusable, ok := child.(Focusable); ok {
			eles = append(eles, focusable)
		}
	}

	return eles
}

//
func (hl *HLayout) NextFocusable(current Focusable) Focusable {
	var eles []Focusable = hl.GetFocusable()

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
func (hl *HLayout) FocusClicked(ev termbox.Event) Focusable {
	// var w, h int = hl.Size()

	// termbox uses coords based from 1, 1 not 0, 0
	// keep it consistent for bubbling through containers
	// but -1,-1 for the HandleClick methods

	// if i ever add the Clickable interface to HLayout
	// if mouseX > 0 && mouseY > 0 && mouseX <= w && mouseY <= h {
	// 	hl.HandleClick(mouseX, mouseY)
	// }

	// normalise mouse click to this element so it can be
	// passed down to children

	// adjust for padding
	ev.MouseX, ev.MouseY = ev.MouseX-hl.Padding.Left()-hl.Border.Adjust(Left), ev.MouseY-hl.Padding.Up()-hl.Border.Adjust(Up)
	// w, h = w-(hl.Padding.Left()+hl.Padding.Right())-(hl.Border.Left()+hl.Border.Right()), h-(hl.Padding.Up()+hl.Padding.Down())-(hl.Border.Up()+hl.Border.Down())

	var sumX int
	for _, c := range hl.Children {
		var cw, ch int = c.Size()

		if ev.MouseX >= sumX && ev.MouseY >= 0 && ev.MouseX <= sumX+cw && ev.MouseY <= ch {
			// update event before passing it down
			ev.MouseX -= sumX

			if clickable, ok := c.(Clickable); ok {
				clickable.HandleClick(ev)
			}

			if cont, ok := c.(Container); ok {
				return cont.FocusClicked(ev)
			} else if foc, ok := c.(Focusable); ok {
				return foc
			}
		}

		sumX += cw
	}

	return nil
}
