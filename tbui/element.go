package tbui

import termbox "github.com/nsf/termbox-go"

//
type Element interface {
	Draw(int, int, Element)
	Size() (int, int)
	Handle(termbox.Event)
	HandleClick(int, int)
	Focusable() bool
}

//
type Container interface {
	NextFocusable(Element) Element
	GetFocusable() []Element
	FocusClicked(int, int) Element
}
