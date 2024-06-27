package components

import (
	"dbdb/state"

	"github.com/hajimehoshi/ebiten/v2"
)

type Drawer interface {
	Draw(*ebiten.Image)
}

func NewPopoverArea(left, top, width, height int, img Drawer, c *state.ControlHandler, click func()) *PopoverArea {
	return &PopoverArea{
		top: top,
		left: left,
		width: width,
		height: height,
		popImage: img,
		onClick: click,
		ControlHandler: c,
	}
}

type PopoverArea struct {
	top int
	left int
	height int
	width int

	popImage Drawer

	onClick func()

	showImg bool
	*state.ControlHandler
}

func (p *PopoverArea) Draw(screen *ebiten.Image) {
	if IsMouseover(p.left, p.top, p.width, p.height, false) {
		if !p.showImg {
			p.showImg = true
		}
		if p.LeftClick && p.onClick != nil {
			p.onClick()
		}
	} else if p.showImg {
		p.showImg = false
	}

	if p.showImg {
		p.popImage.Draw(screen)
	}
}
