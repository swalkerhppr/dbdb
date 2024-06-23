package components

import (
	"dbdb/assets"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"github.com/yohamta/ganim8/v2"
)


func NewTextBox(text string, fontSize, left, top, width, height int) *TextBox {
	return &TextBox{
		text   : text,
		fontSize: fontSize,
		left   : left,
		top    : top,
		width  : width,
		height : height,
		tileSprite : assets.Registry.Sprite("textbox-tiles"),
	}
}

type TextBox struct {
	text   string
	fontSize int

	left   int
	top    int
	width  int
	height int
	tileSprite *ganim8.Sprite
}
func (t *TextBox) SetText(text string) {
	t.text = text
}

func (t *TextBox) Draw(screen *ebiten.Image) {
	// Draw Text Box
	// _______
	// |0 1 2|
	// |3 4 5|
	// |6 7 8|
	// -------
	tilesWide := t.width / 32
	tilesHigh := t.height / 32
	// top edge
	for x := range tilesWide - 2 {
		t.tileSprite.Draw(screen, 1, ganim8.DrawOpts(float64(t.left + ((x + 1) * 32)), float64(t.top)))
	}
	// bottom edge
	for x := range tilesWide - 2 {
		t.tileSprite.Draw(screen, 7, ganim8.DrawOpts(float64(t.left + ((x + 1) * 32)), float64(t.top + t.height - 32)))
	}
	// center pieces
	for y := range tilesHigh - 2 {
		// left wall
		t.tileSprite.Draw(screen, 3, ganim8.DrawOpts(float64(t.left), float64(t.top + (y + 1) * 32)))
		// right wall
		t.tileSprite.Draw(screen, 5, ganim8.DrawOpts(float64(t.left + t.width - 32), float64(t.top + (y + 1) * 32)))
		for x := range tilesWide - 2 {
			t.tileSprite.Draw(screen, 4, ganim8.DrawOpts(float64(t.left + ((x + 1) * 32)), float64(t.top + ((y + 1) * 32))))
		}
	}
	// corners
	t.tileSprite.Draw(screen, 0, ganim8.DrawOpts(float64(t.left), float64(t.top)))
	t.tileSprite.Draw(screen, 2, ganim8.DrawOpts(float64(t.left + t.width - 32), float64(t.top)))
	t.tileSprite.Draw(screen, 6, ganim8.DrawOpts(float64(t.left), float64(t.top + t.height - 32)))
	t.tileSprite.Draw(screen, 8, ganim8.DrawOpts(float64(t.left + t.width - 32), float64(t.top + t.height - 32)))

	r := assets.Registry.DefaultTextRenderer(t.fontSize)
	r.SetAlign(etxt.YCenter, etxt.XCenter)
	r.SetColor(color.Black)
	r.SetTarget(screen)
	r.Draw(t.text, t.left + t.width/2, t.top + t.height / 2)
}
