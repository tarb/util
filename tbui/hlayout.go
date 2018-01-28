package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type HLayout struct {
	Children []Element

	Border  Border
	Padding Padding
}

//
func (hl *HLayout) Draw(x, y int, focused Element) {
	// border offset
	var bdr int
	if hl.Border != None {
		bdr = 1
		hl.drawBorder(x, y)
	}

	// x and y to start drawing children
	var eX, eY int = x + hl.Padding.Left() + bdr, y + hl.Padding.Up() + bdr

	for _, e := range hl.Children {
		var eWidth, _ int = e.Size()
		e.Draw(eX, eY, focused)
		eX += eWidth
	}
}

//
func (hl *HLayout) Size() (int, int) {
	var cumulativeX, maxY, bdr int

	for _, e := range hl.Children {
		var w, h int = e.Size()

		cumulativeX += w
		if h > maxY {
			maxY = h
		}
	}

	if hl.Border != None {
		bdr = 2
	}

	cumulativeX += (hl.Padding.Left() + hl.Padding.Right()) + bdr
	maxY += (hl.Padding.Up() + hl.Padding.Down()) + bdr

	return cumulativeX, maxY
}

//
func (hl *HLayout) Handle(ev termbox.Event) {}

//
func (hl *HLayout) HandleClick(mouseX, mouseY int) {
	//fmt.Println("hlayout", hl.Border, "|", mouseX, mouseY, hl.Padding)
}

//
func (hl *HLayout) Focusable() bool { return false }

func (hl *HLayout) drawBorder(x, y int) {
	var runes []rune = borderRunes[hl.Border]
	var w, h int = hl.Size()

	// x
	for i := x + 1; i < x+w-1; i++ {
		termbox.SetCell(i, y, runes[0], termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(i, y+h-1, runes[0], termbox.ColorDefault, termbox.ColorDefault)
	}
	// y
	for i := y + 1; i < y+h-1; i++ {
		termbox.SetCell(x, i, runes[1], termbox.ColorDefault, termbox.ColorDefault)
		termbox.SetCell(x+w-1, i, runes[1], termbox.ColorDefault, termbox.ColorDefault)
	}
	// corners
	termbox.SetCell(x, y, runes[2], termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x+w-1, y, runes[3], termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x, y+h-1, runes[4], termbox.ColorDefault, termbox.ColorDefault)
	termbox.SetCell(x+w-1, y+h-1, runes[5], termbox.ColorDefault, termbox.ColorDefault)
}

//
func (hl *HLayout) GetFocusable() []Element {
	var eles = make([]Element, 0, 10)

	if hl.Focusable() {
		eles = append(eles, hl)
	}

	for _, child := range hl.Children {
		if cont, ok := child.(Container); ok {
			eles = append(eles, cont.GetFocusable()...)
		} else if child.Focusable() {
			eles = append(eles, child)
		}
	}

	return eles
}

//
func (hl *HLayout) NextFocusable(current Element) Element {
	var eles []Element = hl.GetFocusable()

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
func (hl *HLayout) FocusClicked(mouseX, mouseY int) Element {
	var w, h int = hl.Size()

	// termbox uses coords based from 1, 1 not 0, 0
	// keep it consistent for bubbling through containers
	// but -1,-1 for the HandleClick methods

	if mouseX > 0 && mouseY > 0 && mouseX <= w && mouseY <= h {
		hl.HandleClick(mouseX, mouseY)
	}

	// normalise mouse click to this element so it can be
	// passed down to children

	// adjust for padding
	mouseX, mouseY = mouseX-hl.Padding.Left(), mouseY-hl.Padding.Up()
	w, h = w-(hl.Padding.Left()+hl.Padding.Right()), h-(hl.Padding.Up()+hl.Padding.Down())
	// adjust for border
	if hl.Border != None {
		mouseX, mouseY = mouseX-1, mouseY-1
		w, h = w-2, h-2
	}

	var sumX int
	for _, c := range hl.Children {
		var cw, ch int = c.Size()

		if mouseX > sumX && mouseY > 0 && mouseX <= sumX+cw && mouseY <= ch {
			if cont, ok := c.(Container); ok {
				return cont.FocusClicked(mouseX-sumX, mouseY)
			}
			c.HandleClick(mouseX-sumX, mouseY) //needs to be normalized
			return c
		}

		sumX += cw
	}

	return nil
}
