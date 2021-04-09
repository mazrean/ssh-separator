package values

type Window struct {
	height uint
	width  uint
}

func NewWindow(height uint, width uint) *Window {
	return &Window{
		height: height,
		width:  width,
	}
}

func (w *Window) Height() uint {
	return w.height
}

func (w *Window) Width() uint {
	return w.width
}
