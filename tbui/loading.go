package tbui

import (
	"time"

	termbox "github.com/nsf/termbox-go"
)

// the different frames of the loading animation
var frames = [5]string{"██  ", "▐█▌ ", " ██ ", " ▐█▌", "  ██"}

// LoadingTick - duration between repaints
var LoadingTick = 100 * time.Millisecond

//
type Loading struct {
	Padding Padding
}

//
func (l *Loading) Draw(x, y int, focused Element) {

	x, y = x+l.Padding.Left(), y+l.Padding.Up()
	// counts 0,1,2,3,4,4,3,2,1,0 ... repeat
	var n = time.Now().UnixNano() / int64(LoadingTick) % 10
	n = n + (n / 5 * (((n % 5) * -2) - 1))

	var i = 5
	for _, c := range frames[n] {
		termbox.SetCell(x+i, y, c, termbox.ColorWhite, termbox.ColorDefault)
		i++
	}
}

//
func (l *Loading) Size() (int, int) {
	return l.Padding.Left() + 4 + l.Padding.Right(), l.Padding.Up() + 1 + l.Padding.Down()
}
