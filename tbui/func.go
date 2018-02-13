package tbui

//
type Func struct {
	F func() Element
}

//
func (f *Func) Draw(x, y int, focused Element) {
	var e = f.F()
	if ele, ok := e.(Element); ok {
		ele.Draw(x, y, focused)
	}
}

//
func (f *Func) Size() (int, int) {
	var e = f.F()
	if ele, ok := e.(Element); ok {
		return ele.Size()
	}
	return 0, 0
}
