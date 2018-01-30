package main

import (
	"fmt"

	termbox "github.com/nsf/termbox-go"
	ui "github.com/tarb/util/tbui"
)

func main() {
	if err := termbox.Init(); err != nil {
		panic(err)
	}
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)
	defer termbox.Close()

	var uname, email, pswd *ui.EditBox
	var save *ui.CheckBox

	var submit = func() {
		// do anything here
		fmt.Printf("User: %s, Email: %s, Pwd: %s, Save: %v\n", uname.Text(), email.Text(), pswd.Text(), save.Checked)
	}

	uname = &ui.EditBox{Width: 20, Padding: ui.Padding{0, 1}, Submit: submit}
	email = &ui.EditBox{Width: 20, Padding: ui.Padding{0, 1}, Submit: submit}
	pswd = &ui.EditBox{Width: 20, Padding: ui.Padding{0, 1}, HideContent: true, Submit: submit}
	save = &ui.CheckBox{Padding: ui.Padding{1, 0}, Submit: submit}

	var focus ui.Focusable
	// var window = &ui.VLayout{ // outer window, taking up whole screen - padding adjusted on window resize
	// 	Children: []ui.Element{
	// 		&ui.VLayout{ //inner window, with border
	// 			Children: []ui.Element{
	// 				&ui.Text{Text: "Sample Login", Padding: ui.Padding{1}, Width: 30, Allign: ui.Center},
	// 				&ui.HLayout{Children: []ui.Element{&ui.Text{Text: "Username", Padding: ui.Padding{1}, Width: 8}, uname}},
	// 				&ui.HLayout{Children: []ui.Element{&ui.Text{Text: "Email", Padding: ui.Padding{1}, Width: 8, Allign: ui.Right}, email}},
	// 				&ui.HLayout{Children: []ui.Element{&ui.Text{Text: "Password", Padding: ui.Padding{1}, Width: 8}, pswd}},
	// 				&ui.HLayout{Children: []ui.Element{&ui.Text{Text: "Save Details?", Padding: ui.Padding{1}, Width: 22, Allign: ui.Right}, save}},
	// 				&ui.HLayout{Children: []ui.Element{&ui.Button{Text: "Login", Padding: ui.Padding{0, 3}, Submit: submit}}, Padding: ui.Padding{2, 0, 0, 12}},
	// 			},
	// 			Padding: ui.Padding{2, 4},
	// 			Border:  ui.Thin,
	// 		},
	// 	},
	// }

	var window = &ui.VLayout{
		Children: []ui.Element{
			&ui.List{
				Items: []ui.Element{
					&ui.EditBox{PlaceHolder: "alpha", Padding: ui.Padding{1}, Width: 8},
					// &ui.VLayout{Children: []ui.Element{
					// 	&ui.EditBox{PlaceHolder: "bravo", Padding: ui.Padding{1}, Width: 8},
					// },
					&ui.EditBox{PlaceHolder: "bravo", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "charlie", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "delta", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "echo", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "foxtrot", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "geko", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "hotel", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "alpha", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "bravo", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "charlie", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "delta", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "echo", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "foxtrot", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "geko", Padding: ui.Padding{1}, Width: 8},
					&ui.EditBox{PlaceHolder: "hotel", Padding: ui.Padding{1}, Width: 8},
				},
				Height: 25,
			},
		},
	}

	var paint = func() {
		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		window.Draw(0, 0, focus)
		termbox.Flush()
	}

	var resizeWindow = func() {
		tw, th := termbox.Size()
		ew, eh := window.Children[0].Size()
		window.Padding = ui.Padding{(th - eh) / 2, (tw - ew) / 2}
	}

	resizeWindow()

mainloop:
	for {
		paint()

		switch ev := termbox.PollEvent(); ev.Type {

		case termbox.EventMouse:
			if ev.Key == termbox.MouseLeft {
				focus = window.FocusClicked(ev.MouseX, ev.MouseY)
			}

		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop //quit
			case termbox.KeyTab:
				focus = window.NextFocusable(focus) //change focus
			default:
				if focus != nil {
					focus.Handle(ev) //pass event down
				}
			}

		case termbox.EventResize:
			resizeWindow()

		case termbox.EventError:
			panic(ev.Err)
		}

	}
}
