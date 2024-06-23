package components

import (

	"github.com/hajimehoshi/ebiten/v2"
)

func IsMouseover(left, top, w, h int, center bool) bool {
	ww, wh := ebiten.WindowSize()
	xm, ym := ebiten.CursorPosition()
	if center {
		top = top - h/2
		left = left - w/2
	}
	yscale := float64(wh)/480
	top = int(float64(top) * yscale)
	h = int(float64(h) * yscale)

	xscale := float64(ww)/640
	left = int(float64(left) * xscale)
	w = int(float64(w) * xscale)


	return ym > top && ym < top + h && xm > left && xm < left + w
}
