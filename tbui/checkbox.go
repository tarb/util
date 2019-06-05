package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type CheckBox struct {
	Checked bool
	Symbol  rune
	Padding Padding
	Submit  func()
	Bind    *bool
}

//
func (cb *CheckBox) Draw(x, y int, focused Element) {
	x, y = x+cb.Padding.Left(), y+cb.Padding.Up()

	var checkCol = termbox.ColorBlack | termbox.AttrBold
	if focused == cb {
		checkCol = termbox.ColorWhite
	}

	var mark rune = ' '
	if cb.Checked {
		if cb.Symbol != 0 {
			mark = cb.Symbol
		} else {
			mark = 'x'
		}
	}

	termbox.SetCell(x, y, '▐', checkCol, termbox.ColorDefault)
	termbox.SetCell(x+1, y, mark, termbox.ColorRed, checkCol)
	termbox.SetCell(x+2, y, '▌', checkCol, termbox.ColorDefault)
}

//
func (cb *CheckBox) Size() (int, int) {
	return cb.Padding.Left() + 3 + cb.Padding.Right(), cb.Padding.Up() + 1 + cb.Padding.Down()
}

//
func (cb *CheckBox) Handle(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeySpace:
		cb.check()
	case termbox.KeyEnter:
		if cb.Submit != nil {
			cb.Submit()
		}
	}
}

//
func (cb *CheckBox) HandleClick(ev termbox.Event) {
	//fmt.Println("checkbox", mouseX, mouseY, cb.Padding)
	if ev.MouseX >= cb.Padding.Left() && ev.MouseX < cb.Padding.Left()+3 && ev.MouseY >= cb.Padding.Up() && ev.MouseY < cb.Padding.Up()+1 {
		cb.check()
	}
}

//
func (cb *CheckBox) check() {
	cb.Checked = !cb.Checked
	if cb.Bind != nil {
		*cb.Bind = cb.Checked
	}
}
