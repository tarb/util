package tbui

//
type Window interface {
	Listen() WindowResult
}

//
type WindowResult struct {
	State   WindowState
	Message string
}

//
type WindowState uint8

//
const (
	Error    WindowState = 0
	Exit                 = 1
	Back                 = 2
	Continue             = 3
)
