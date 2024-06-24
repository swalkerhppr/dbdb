package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)


type Alert struct {
	box TextBox
	render bool
}

func NewAlert() *Alert {
	return &Alert{
		box: *NewTextBox("Enter Text Here!", 16, 180, 185, 320, 42),
	}
}

func (a *Alert) Show() {
	a.render = true
}

func (a *Alert) SetText(s string) {
	a.box.SetText(s)
}

func (a *Alert) SetPosition(left, top int) {
	a.box.left = left
	a.box.top = top
}

func (a *Alert) Draw(screen *ebiten.Image) {
	click := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	if a.render {
		if click {
			a.render = false
		} else {
			a.box.Draw(screen)
		}
	}
}
