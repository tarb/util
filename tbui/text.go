package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type Text struct {
	Text    string
	DText   func() string
	Width   int
	Align   Align
	Padding Padding
	BgCol   termbox.Attribute
	TextCol termbox.Attribute
}

//
func (t *Text) Draw(x, y int, focused Element) {
	// background shading
	var cw, rh int = t.Size()
	for r := y; r < y+rh; r++ {
		for c := x; c < x+cw; c++ {
			termbox.SetCell(c, r, ' ', t.BgCol, t.BgCol)
		}
	}

	if t.DText != nil {
		t.Text = t.DText()
	}

	var w int = len(t.Text)
	if t.Width != 0 && t.Width < w {
		w = t.Width
	}
	x, y = x+t.Padding.Left(), y+t.Padding.Up()

	if t.Align == AlignRight {
		x += (t.Width - w)
	} else if t.Align == AlignCenter {
		x += ((t.Width - w) / 2)
	}

	for i, c := range t.Text[:w] {
		termbox.SetCell(x+i, y, c, t.TextCol, t.BgCol)
	}
}

//
func (t *Text) Size() (int, int) {
	var w int = len(t.Text)
	if t.Width != 0 {
		w = t.Width
	}
	return t.Padding.Left() + t.Padding.Right() + w, t.Padding.Up() + t.Padding.Down() + 1
}
