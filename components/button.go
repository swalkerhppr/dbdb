package components

import (
	"dbdb/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"github.com/yohamta/ganim8/v2"
)

func NewButton(text string, x, y int, click func()) *Button {
	r := assets.Registry.DefaultTextRenderer(16)
	r.SetAlign(etxt.YCenter, etxt.XCenter)
	return &Button{
		sprite       : assets.Registry.Sprite("MainButton"),
		textRenderer : r,

		centerX : x,
		centerY : y,
		text    : text,
		OnClick : click,
	}
}

type Button struct {
	text         string
	textRenderer *etxt.Renderer
	sprite       *ganim8.Sprite
	centerX      int
	centerY      int
	pressed      bool
	OnClick      func()
}

func (b *Button) Draw(screen *ebiten.Image) {
	b.textRenderer.SetTarget(screen)
	clicked := ebiten.IsMouseButtonPressed(ebiten.MouseButton0)
	frameIdx := 1

	if IsMouseover(b.centerX, b.centerY, b.sprite.W(), b.sprite.H(), true) {
		frameIdx = 0
		if clicked {
			b.pressed = true
		} else if b.pressed && b.OnClick != nil {
			b.OnClick()
			b.pressed = false
		}
	}
	if !clicked {
		b.pressed = false
	}
	if b.pressed {
		frameIdx = 2
	}

	opts := ganim8.DrawOpts(float64(b.centerX), float64(b.centerY), 0, 1, 1, .5, .5)
	ganim8.DrawSpriteWithOpts(screen, b.sprite, frameIdx, opts, nil)
	b.textRenderer.Draw(b.text, b.centerX, b.centerY)
}
