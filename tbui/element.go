package tbui

import termbox "github.com/nsf/termbox-go"

//
type Element interface {
	Draw(int, int, Element)
	Size() (int, int)
}

//
type Expandable interface {
	Element
	ExpandSize() (int, int)
	ExpandDraw(int, int, Element)
}

//
type Focusable interface {
	Element
	Handle(termbox.Event)
}

//
type Clickable interface {
	Element
	HandleClick(int, int)
}

//
type Container interface {
	Element
	NextFocusable(Focusable) Focusable
	GetFocusable() []Focusable
	FocusClicked(int, int) Focusable
}
