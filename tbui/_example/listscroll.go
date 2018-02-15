package main

import (
	"log"

	termbox "github.com/nsf/termbox-go"
	ui "github.com/tarb/util/tbui"
)

//
type ListExampleWindow struct {
	focus ui.Focusable
	root  *ui.HLayout
}

//
func NewListExampleWindow() *ListExampleWindow {
	var lw ListExampleWindow

	var words = []string{"alpha", "bravo", "charlie", "delta", "echo", "foxtrot", "gecko", "hotel", "igloo", "juliet", "kilo", "lama", "mike", "navi", "oscar", "oreo", "parrot"}

	lw.root = &ui.HLayout{
		Children: []ui.Element{
			&ui.VLayout{
				Children: []ui.Element{
					&ui.DynamicList{
						BindBuilder: func(i int) ui.Element {
							return &ui.Text{Text: words[i]}
						},
						BindSize: func() int { return len(words) },
						Height:   10,
					},
				},
				Border:  ui.Border{Style: ui.Thin, Side: ui.All, Bg: termbox.ColorBlack | termbox.AttrBold},
				Padding: ui.Padding{1, 2},
			},
			&ui.VLayout{
				Children: []ui.Element{
					&ui.DynamicList{
						BindBuilder: func(i int) ui.Element {
							return &ui.Text{Text: words[i]}
						},
						BindSize: func() int { return len(words) },
						Height:   10,
					},
				},
				Border:  ui.Border{Style: ui.Thin, Side: ui.All, Bg: termbox.ColorBlack | termbox.AttrBold},
				Padding: ui.Padding{1, 2},
			},
		},
	}
	return &lw
}

//
func (lw *ListExampleWindow) Listen() {
	// set up termbox
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	defer termbox.Close()

	// schedule a new paint call
	var paint = make(chan struct{}, 1)
	var schedPaint = func() {
		select {
		case paint <- struct{}{}:
		default:
		}
	}
	// termbox event poller
	var events = make(chan termbox.Event)
	go func() {
		for {
			var ev termbox.Event = termbox.PollEvent()
			if ev.Type == termbox.EventInterrupt {
				return
			}

			events <- ev
		}
	}()

	// paint, update messages and handle user events
	schedPaint()

Exit:
	for {
		//listen events, messages
		select {
		case <-paint:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			lw.root.Draw(0, 0, lw.focus)
			termbox.Flush()

		// User has fired event
		case ev := <-events:
			switch ev.Type {
			case termbox.EventMouse:
				if ev.Key == termbox.MouseLeft {
					lw.focus = lw.root.FocusClicked(ev)
				} else if ev.Key == termbox.MouseWheelDown || ev.Key == termbox.MouseWheelUp {
					lw.root.FocusClicked(ev)
				}
				schedPaint()

			case termbox.EventKey:
				switch ev.Key {
				case termbox.KeyEsc, termbox.KeyCtrlC:
					break Exit
				}

			case termbox.EventResize:
				schedPaint()

			case termbox.EventError:
				panic(ev.Err)
			}
		}
	}

	termbox.Interrupt()
}
