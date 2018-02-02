package tbui

import termbox "github.com/nsf/termbox-go"

//
type Border int

//
const (
	None   Border = 0
	Thin   Border = 1
	Thick  Border = 2
	Double Border = 3
)

var borderRunes = [][]rune{
	nil,
	[]rune{'─', '│', '┌', '┐', '└', '┘'},
	[]rune{'━', '┃', '┏', '┓', '┗', '┛'},
	[]rune{'═', '║', '╔', '╗', '╚', '╝'},
}

//
type Allign uint

//
const (
	Left   Allign = 0
	Center Allign = 1
	Right  Allign = 2
)

//
type Padding []int

//
func (p Padding) Left() int {
	if len(p) == 0 {
		return 0
	}
	return p[[4]int{0, 1, 1, 3}[len(p)-1]]
}

//
func (p Padding) Right() int {
	if len(p) == 0 {
		return 0
	}
	return p[[4]int{0, 1, 1, 1}[len(p)-1]]
}

//
func (p Padding) Up() int {
	if len(p) == 0 {
		return 0
	}
	return p[[4]int{0, 0, 0, 0}[len(p)-1]]
}

//
func (p Padding) Down() int {
	if len(p) == 0 {
		return 0
	}
	return p[[4]int{0, 0, 2, 2}[len(p)-1]]
}

//
const (
	ColBackground = termbox.ColorBlack
	ColText       = termbox.ColorWhite
	ColBorder     = termbox.ColorBlack | termbox.AttrBold
	ColAccent     = termbox.ColorRed | termbox.AttrBold
)
