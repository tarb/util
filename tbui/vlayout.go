package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type VLayout struct {
	Children []Element

	MinWidth  int
	MinHeight int
	Border    Border
	Padding   Padding
}

//
func (vl *VLayout) Draw(x, y int, focused Element) {
	// x and y to start drawing children
	var eX, eY int = x + vl.Padding.Left() + vl.Border.Adjust(Left), y + vl.Padding.Up() + vl.Border.Adjust(Up)

	for _, e := range vl.Children {
		e.Draw(eX, eY, focused)
		var _, eHeight int = e.Size()
		eY += eHeight
	}
	vl.drawBorder(x, y)
}

//
func (vl *VLayout) Size() (int, int) {
	var cumulativeY, maxX int

	for _, e := range vl.Children {
		var w, h int = e.Size()

		cumulativeY += h
		if w > maxX {
			maxX = w
		}
	}

	if maxX < vl.MinWidth {
		maxX = vl.MinWidth
	}
	if cumulativeY < vl.MinHeight {
		cumulativeY = vl.MinHeight
	}

	cumulativeY += vl.Padding.Up() + vl.Padding.Down() + vl.Border.Adjust(Up) + vl.Border.Adjust(Down)
	maxX += vl.Padding.Left() + vl.Padding.Right() + vl.Border.Adjust(Left) + vl.Border.Adjust(Right)

	return maxX, cumulativeY
}

//
// func (vl *VLayout) HandleClick(mouseX, mouseY int) {
// 	fmt.Println("vlayout", vl.Border, "|", mouseX, mouseY, vl.Padding)
// }

func (vl *VLayout) drawBorder(x, y int) {
	var runes []rune = vl.Border.Runes()
	var fg, bg = vl.Border.Fg, vl.Border.Bg
	var w, h int = vl.Size()

	// x
	if vl.Border.Has(Up) {
		for i := x; i < x+w; i++ {
			termbox.SetCell(i, y, runes[0], fg, bg)
		}
	}
	if vl.Border.Has(Down) {
		for i := x; i < x+w; i++ {
			termbox.SetCell(i, y+h-1, runes[0], fg, bg)
		}
	}
	// y
	if vl.Border.Has(Left) {
		for i := y; i < y+h; i++ {
			termbox.SetCell(x, i, runes[1], fg, bg)
		}
	}
	if vl.Border.Has(Right) {
		for i := y; i < y+h; i++ {
			termbox.SetCell(x+w-1, i, runes[1], fg, bg)
		}
	}

	// corners
	if vl.Border.Has(Left | Up) {
		termbox.SetCell(x, y, runes[2], fg, bg)
	}
	if vl.Border.Has(Left | Down) {
		termbox.SetCell(x, y+h-1, runes[4], fg, bg)
	}
	if vl.Border.Has(Right | Up) {
		termbox.SetCell(x+w-1, y, runes[3], fg, bg)
	}
	if vl.Border.Has(Right | Down) {
		termbox.SetCell(x+w-1, y+h-1, runes[5], fg, bg)
	}
}

//
func (vl *VLayout) GetFocusable() []Focusable {
	var eles = make([]Focusable, 0, 10)

	for _, child := range vl.Children {
		if cont, ok := child.(Container); ok {
			eles = append(eles, cont.GetFocusable()...)
		} else if focusable, ok := child.(Focusable); ok {
			eles = append(eles, focusable)
		}
	}

	return eles
}

//
func (vl *VLayout) NextFocusable(current Focusable) Focusable {
	var eles []Focusable = vl.GetFocusable()

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
func (vl *VLayout) FocusClicked(mouseX, mouseY int) Focusable {
	// var w, h int = vl.Size()

	// adjust for padding
	mouseX, mouseY = mouseX-vl.Padding.Left()-vl.Border.Adjust(Left), mouseY-vl.Padding.Up()-vl.Border.Adjust(Up)
	// w, h = w-(vl.Padding.Left()+vl.Padding.Right()), h-(vl.Padding.Up()+vl.Padding.Down())

	var sumY int
	for _, c := range vl.Children {
		var cw, ch int = c.Size()

		if mouseX >= 0 && mouseY >= sumY && mouseX < cw && mouseY < sumY+ch {
			if clickable, ok := c.(Clickable); ok {
				clickable.HandleClick(mouseX, mouseY-sumY)
			}
			if cont, ok := c.(Container); ok {
				return cont.FocusClicked(mouseX, mouseY-sumY)
			} else if foc, ok := c.(Focusable); ok {
				return foc
			}
		}

		sumY += ch
	}

	return nil
}
