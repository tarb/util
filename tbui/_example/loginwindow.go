package main

import (
	"log"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
	ui "github.com/tarb/util/tbui"
)

//
type LoginWindow struct {
	username string
	password string
	remember bool

	focus  ui.Focusable
	root   *ui.VLayout
	notice *ui.HLayout

	call       func(string, string, bool) (string, error)
	loginQueue chan loginAttempt
}

type loginAttempt struct {
	username string
	password string
	remember bool
}

//
func NewLoginWindow(submit func(string, string, bool) (string, error)) *LoginWindow {
	var lw LoginWindow

	lw.call = submit
	lw.loginQueue = make(chan loginAttempt)

	lw.notice = &ui.HLayout{MinHeight: 1, Padding: ui.Padding{0, 0, 2}}
	lw.root = &ui.VLayout{
		Children: []ui.Element{
			lw.notice,

			&ui.Text{Text: "Username", Width: 8},
			&ui.EditBox{Width: 26, Submit: lw.submit, Bind: &lw.username, Padding: ui.Padding{0, 1}},

			&ui.Text{Text: "Password", Width: 8},
			&ui.EditBox{Width: 26, HideContent: true, Submit: lw.submit, Bind: &lw.password, Padding: ui.Padding{0, 1}},

			&ui.HLayout{Children: []ui.Element{
				&ui.Text{Text: "Remember?"},
				&ui.CheckBox{Submit: lw.submit, Bind: &lw.remember},
			}, Padding: ui.Padding{0, 0, 0, 9}},

			&ui.HLayout{Children: []ui.Element{
				&ui.Button{Text: "Login", Submit: lw.submit, Padding: ui.Padding{0, 5}},
			}, Padding: ui.Padding{2, 0, 0, 7}},
		},
		Padding: ui.Padding{1, 2},
		Border:  ui.Border{Style: ui.Thin, Side: ui.All, Bg: termbox.ColorBlack | termbox.AttrBold},
	}
	return &lw
}

//
func (lw *LoginWindow) submit() {
	select {
	case lw.loginQueue <- loginAttempt{username: lw.username, password: lw.password, remember: lw.remember}:
	default:
	}
}

//
func (lw *LoginWindow) Listen() {
	// set up termbox
	if err := termbox.Init(); err != nil {
		log.Fatal(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	defer termbox.Close()

	// closes on quit
	var quitter sync.Once
	var quit = make(chan struct{})

	// messages to write out
	var notify = make(chan string)

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
	go func() {
		var centerPos = func() (int, int) {
			tw, th := termbox.Size()
			ew, eh := lw.root.Size()
			return (tw - ew) / 2, (th - eh) / 2
		}
		schedPaint()

		for {
			//listen events, messages
			select {
			case <-paint:
				termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
				var xOffset, yOffset int = centerPos()
				lw.root.Draw(xOffset, yOffset, lw.focus)
				termbox.Flush()

			// LoginResult has come back
			case msg := <-notify:
				lw.notice.Children = []ui.Element{&ui.Text{Text: msg, Width: 28, Align: ui.AlignCenter}}
				schedPaint()

			// User has fired event
			case ev := <-events:
				switch ev.Type {
				case termbox.EventMouse:
					if ev.Key == termbox.MouseLeft {
						var xOffset, yOffset int = centerPos()
						ev.MouseX, ev.MouseY = ev.MouseX-xOffset, ev.MouseY-yOffset
						lw.focus = lw.root.FocusClicked(ev)
						schedPaint()
					}
				case termbox.EventKey:
					switch ev.Key {
					case termbox.KeyEsc, termbox.KeyCtrlC:
						quitter.Do(func() { close(quit) })
					case termbox.KeyTab:
						lw.focus = lw.root.NextFocusable(lw.focus)
					default:
						if lw.focus != nil {
							lw.focus.Handle(ev) //pass event down
						}
					}
					schedPaint()

				case termbox.EventResize:
					schedPaint()

				case termbox.EventError:
					panic(ev.Err)
				}

			case <-quit:
				return
			}

		}
	}()

	// process requests
	for {
		select {
		case a := <-lw.loginQueue:
			// add loading widget here
			var tick = time.NewTicker(100 * time.Millisecond)
			lw.notice.Children = []ui.Element{&ui.Loading{Padding: ui.Padding{0, 7}}}
			go func() {
				for range tick.C {
					schedPaint()
				}
			}()

			if msg, err := lw.call(a.username, a.password, a.remember); err == nil {
				quitter.Do(func() { close(quit) })
			} else {
				notify <- msg
			}
			tick.Stop()

		case <-quit:
			termbox.Interrupt()
			return
		}
	}

}
