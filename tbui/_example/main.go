package main

func main() {
	// var lw = NewLoginWindow(func(u, p string, r bool) (string, error) {
	// 	// call log in function here
	// 	return "", nil
	// })
	// lw.Listen()

	window := NewListExampleWindow()
	window.Listen()
}
