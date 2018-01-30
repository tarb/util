package tbui

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
)

//
type VLayout struct {
	Children []Element

	Border  Border
	Padding Padding
}

//
func (vl *VLayout) Draw(x, y int, focused Element) {
	// border offset
	var bdr int
	if vl.Border != None {
		bdr = 1
		vl.drawBorder(x, y)
	}

	// x and y to start drawing children
	var eX, eY int = x + vl.Padding.Left() + bdr, y + vl.Padding.Up() + bdr

	for _, e := range vl.Children {
		var _, eHeight int = e.Size()
		e.Draw(eX, eY, focused)
		eY += eHeight
	}
}

//
func (vl *VLayout) Size() (int, int) {
	var cumulativeY, maxX, bdr int

	for _, e := range vl.Children {
		var w, h int = e.Size()

		cumulativeY += h
		if w > maxX {
			maxX = w
		}
	}

	if vl.Border != None {
		bdr = 2
	}

	cumulativeY += (vl.Padding.Up() + vl.Padding.Down()) + bdr
	maxX += (vl.Padding.Left() + vl.Padding.Right()) + bdr

	return maxX, cumulativeY
}

//
func (vl *VLayout) HandleClick(mouseX, mouseY int) {
	fmt.Println("vlayout", vl.Border, "|", mouseX, mouseY, vl.Padding)
}

//
func (vl *VLayout) drawBorder(x, y int) {
	var runes []rune = borderRunes[vl.Border]
	var w, h int = vl.Size()

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
func (vl *VLayout) GetFocusable() [][]Focusable {
	var eles = make([][]Focusable, 0, 10)

	for _, child := range vl.Children {
		if cont, ok := child.(Container); ok {
			var contCh = cont.GetFocusable()

			if _, ok := child.(Focusable); ok {
				for i := range contCh {
					contCh[i] = append(contCh[i], cont)
				}
			}

			eles = append(eles, cont.GetFocusable()...)
		} else if focusable, ok := child.(Focusable); ok {
			eles = append(eles, focusable)
		}
	}

	fmt.Println(eles)

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
	var w, h int = vl.Size()

	// termbox uses coords based from 1, 1 not 0, 0
	// keep it consistent for bubbling through containers
	// but -1,-1 for the HandleClick methods so they use 0,0

	// if i ever add the Clickable interface to VLayout
	// if mouseX > 0 && mouseY > 0 && mouseX <= w && mouseY <= h {
	// 	vl.HandleClick(mouseX, mouseY)
	// }

	// normalise mouse click to this element so it can be
	// passed down to children

	// adjust for padding
	mouseX, mouseY = mouseX-vl.Padding.Left(), mouseY-vl.Padding.Up()
	w, h = w-(vl.Padding.Left()+vl.Padding.Right()), h-(vl.Padding.Up()+vl.Padding.Down())
	// adjust for border
	if vl.Border != None {
		mouseX, mouseY = mouseX-1, mouseY-1
		w, h = w-2, h-2
	}

	var sumY int
	for _, c := range vl.Children {
		var cw, ch int = c.Size()

		if mouseX > 0 && mouseY > sumY && mouseX <= cw && mouseY <= sumY+ch {
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
