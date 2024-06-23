package components

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type Drawer interface {
	Draw(*ebiten.Image)
}

func NewPopoverArea(left, top, width, height int, img Drawer, click func()) *PopoverArea {
	return &PopoverArea{
		top: top,
		left: left,
		width: width,
		height: height,
		popImage: img,
		onClick: click,
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
	pressed bool
}

func (p *PopoverArea) Draw(screen *ebiten.Image) {
	if IsMouseover(p.left, p.top, p.width, p.height, false) {
		if !p.showImg {
			p.showImg = true
		}
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0) {
			p.pressed = true
		}
		if p.pressed && inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) && p.onClick != nil {
			p.onClick()
		}
	} else if p.showImg {
		p.showImg = false
	}

	if p.showImg {
		p.popImage.Draw(screen)
	}
}
