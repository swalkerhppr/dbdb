package components

import (
	"dbdb/assets"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"github.com/yohamta/ganim8/v2"
)


type Indicator struct {
	sprite *ganim8.Sprite
	text   string
	left   int
	top    int
	textColor color.Color
}

func NewIndicator(name, text string, left, top int) *Indicator {
	return &Indicator{
		sprite : assets.Registry.Sprite(name),
		text   : text,
		left   : left,
		top    : top,
		textColor: color.Black,
	}
}

func NewIndicatorWithColor(name, text string, left, top int, textColor color.Color) *Indicator {
	return &Indicator{
		sprite : assets.Registry.Sprite(name),
		text   : text,
		left   : left,
		top    : top,
		textColor: textColor,
	}
}

func (i *Indicator) Draw(screen *ebiten.Image) {
	i.sprite.Draw(screen, 0, ganim8.DrawOpts(float64(i.left), float64(i.top)))

	r := assets.Registry.DefaultTextRenderer(19)
	r.SetAlign(etxt.YCenter, etxt.Left)
	r.SetColor(i.textColor)
	r.SetTarget(screen)
	r.Draw(i.text, i.left + i.sprite.W() - 7, i.top + (i.sprite.H()/2))
	
}
