package tbui

import (
	"log"

	termbox "github.com/nsf/termbox-go"
)

// EventsCh termbox event poller
var EventsCh = make(chan termbox.Event)

//
func PollEvents() {
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventInterrupt {
			return
		}
		EventsCh <- ev
	}
}

//
func StopPoll() {
	DrainEvents()
	termbox.Interrupt()
}

//
func DrainEvents() {
	for {
		select {
		case <-EventsCh:
		default:
			return
		}
	}
}

//
func Init() {
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
}

//
func Close() {
	termbox.Close()
}
