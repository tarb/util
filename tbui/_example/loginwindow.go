package main

import (
	"context"
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

	call       func(string, string, bool) ui.WindowResult
	loginQueue chan loginAttempt
}

type loginAttempt struct {
	username string
	password string
	remember bool
}

//
func NewLoginWindow(submit func(string, string, bool) ui.WindowResult) *LoginWindow {
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
func (lw *LoginWindow) Listen() ui.WindowResult {
	// messages to write out
	var notify = make(chan string)

	var reqCtx, reqQuit = context.WithCancel(context.Background())

	// schedule a new paint call
	var paint = make(chan struct{}, 1)
	var schedPaint = func() {
		select {
		case paint <- struct{}{}:
		default:
		}
	}

	// closes on quit
	var quitter sync.Once
	var quitCh = make(chan ui.WindowResult, 1)
	var quit = func(r ui.WindowResult) {
		quitter.Do(func() {
			quitCh <- r
			close(quitCh)
			reqQuit()
		})
	}

	// process requests
	go func(ctx context.Context) {
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

				res := lw.call(a.username, a.password, a.remember)

				switch res.State {
				case ui.Continue, ui.Exit, ui.Back:
					quit(res)
				case ui.Error:
					notify <- res.Message
				}

				tick.Stop()

			case <-ctx.Done():
				return
			}
		}
	}(reqCtx)

	// paint, update messages and handle user events

	var centerPos = func() (int, int) {
		tw, th := termbox.Size()
		ew, eh := lw.root.Size()
		return (tw - ew) / 2, (th - eh) / 2
	}
	schedPaint()

	for {
		//listen events, messages
		select {
		case s := <-quitCh:
			return s

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
		case ev := <-ui.EventsCh:
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
				case termbox.KeyEsc:
					quit(ui.WindowResult{State: ui.Exit})
				case termbox.KeyTab:
					lw.focus = lw.root.NextFocusable(lw.focus)
				default:
					if lw.focus != nil {
						lw.focus.Handle(ev) //pass event down
					}
				}
				schedPaint()

			case termbox.EventResize:
				w, h := termbox.Size()
				for c := 0; c < w; c++ {
					for r := 0; r < h; r++ {
						termbox.SetCell(c, r, ' ', termbox.ColorCyan, termbox.ColorCyan)
					}
				}
				termbox.Flush()

				schedPaint()

			case termbox.EventError:
				quit(ui.WindowResult{State: ui.Error, Message: ev.Err.Error()})
			}
		}
	}

}
