package tbui

import (
	termbox "github.com/nsf/termbox-go"
)

//
type Button struct {
	Text    string
	Padding Padding

	Submit func()
}

//
func (b *Button) Draw(x, y int, focused Element) {
	var runes = []rune(b.Text)
	var butCol, textCol = termbox.ColorBlack, termbox.ColorWhite
	if focused == b {
		butCol, textCol = termbox.ColorWhite, termbox.ColorBlack
	}

	x, y = x+b.Padding.Left(), y+1 // so x ==0 && y ==0 is the location of the first char

	//draw background box
	for i := -b.Padding.Left(); i < len(runes)+b.Padding.Right(); i++ {
		termbox.SetCell(x+i, y-1, '▄', butCol, termbox.ColorDefault)
		termbox.SetCell(x+i, y+1, '▀', butCol, termbox.ColorDefault)
		if i >= 0 && i < len(b.Text) {
			termbox.SetCell(x+i, y, runes[i], textCol, butCol)
		} else {
			termbox.SetCell(x+i, y, ' ', butCol, butCol)
		}
	}
}

//
func (b *Button) Size() (int, int) {
	return (b.Padding.Left() + b.Padding.Right()) + len(b.Text), 3
}

//
func (b *Button) Handle(ev termbox.Event) {
	switch ev.Key {
	case termbox.KeyEnter:
		if b.Submit != nil {
			b.Submit()
		}
	}
}

//
func (b *Button) HandleClick(mouseX, mouseY int) {
	//fmt.Println("button", mouseX, mouseY, b.Padding)
	if mouseX >= b.Padding.Left() && mouseX < b.Padding.Left()+len(b.Text) && mouseY >= 0 && mouseY < 3 {
		if b.Submit != nil {
			b.Submit()
		}
	}
}
