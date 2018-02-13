package tbui

import termbox "github.com/nsf/termbox-go"

//
type BorderSide uint8

//
const (
	None      BorderSide = 0
	Up        BorderSide = 2
	Right     BorderSide = 4
	Down      BorderSide = 8
	Left      BorderSide = 16
	UpDown    BorderSide = Up | Down
	LeftRight BorderSide = Left | Right
	All       BorderSide = Up | Right | Down | Left
)

//
type BorderType int

//
const (
	Empty  BorderType = 0
	Thin   BorderType = 1
	Thick  BorderType = 2
	Double BorderType = 3
)

var borderTypes = [][]rune{
	[]rune{' ', ' ', ' ', ' ', ' ', ' '},
	[]rune{'─', '│', '┌', '┐', '└', '┘'},
	[]rune{'━', '┃', '┏', '┓', '┗', '┛'},
	[]rune{'═', '║', '╔', '╗', '╚', '╝'},
}

//
type Border struct {
	Side  BorderSide
	Style BorderType
	Fg    termbox.Attribute
	Bg    termbox.Attribute
}

//
func (b Border) Adjust(bs BorderSide) int {
	if b.Side&bs > 0 {
		return 1
	}
	return 0
}

//
func (b Border) Has(bs BorderSide) bool { return b.Side&bs == bs }

//
func (b Border) Runes() []rune {
	return borderTypes[b.Style]
}

//
type Align uint

//
const (
	AlignLeft   Align = 0
	AlignCenter Align = 1
	AlignRight  Align = 2
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
