package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type Text struct {
	Width   int
	Text    string
	Allign  Allign
	Padding Padding
}

//
func (t *Text) Draw(x, y int, focused Element) {
	var w int = len(t.Text)
	if t.Width != 0 && t.Width < w {
		w = t.Width
	}
	x, y = x+t.Padding.Left(), y+t.Padding.Up()

	if t.Allign == Right {
		x += (t.Width - w)
	} else if t.Allign == Center {
		x += ((t.Width - w) / 2)
	}

	for i, c := range t.Text[:w] {
		termbox.SetCell(x+i, y, c, termbox.ColorDefault, termbox.ColorDefault)
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

//
func (t *Text) Focusable() bool { return false }

//
func (t *Text) Handle(ev termbox.Event) {}

//
func (t *Text) HandleClick(mouseX, mouseY int) {
	//fmt.Println("text", t.Text, mouseX, mouseY, t.Padding)
}
