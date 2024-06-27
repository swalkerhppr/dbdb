package components

import (
	"dbdb/state"

	"github.com/hajimehoshi/ebiten/v2"
)


type Alert struct {
	box TextBox
	render bool
	*state.ControlHandler
}

func NewAlert(c *state.ControlHandler) *Alert {
	return &Alert{
		box: *NewTextBox("Enter Text Here!", 16, 180, 185, 320, 42),
		ControlHandler: c,
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
	if a.render {
		if a.Any {
			a.render = false
		} else {
			a.box.Draw(screen)
		}
	}
}
