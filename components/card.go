package components

import (
	"dbdb/assets"
	"dbdb/state"
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/tinne26/etxt"
	"github.com/yohamta/ganim8/v2"
)

func RandomCard(left, top int) *Card {
	return NewCard(&state.CardState{CardID: state.RandomCardID()}, left, top)
}

func NewCard(cid *state.CardState, left, top int) *Card {
	if cid == nil {
		cid = &state.CardState{
			CardID: state.UnknownCard,
		}
	}
	return &Card{
		cardState : cid,
		sprite : assets.Registry.Sprite(cid.CardID.AssetName()),
		top : top,
		left : left,
	}
}

type Card struct {
	cardState *state.CardState
	sprite *ganim8.Sprite
	top int
	left int
	highlight color.Color
}

func (c *Card) Draw(screen *ebiten.Image) {
	if c.highlight != nil {
		hlImg := ebiten.NewImage(132, 196)
		hlImg.Fill(c.highlight)

		opts := &ebiten.DrawImageOptions{}
		opts.GeoM.Translate(float64(c.left - 2), float64(c.top - 2))

		screen.DrawImage(hlImg, opts)
	}
	c.sprite.Draw(screen, 0, ganim8.DrawOpts(float64(c.left), float64(c.top)))

	timeSymbol := assets.Registry.Sprite("time-symbol")
	if c.cardState.CardID.IsMaterial() {
		qualitySymbol := assets.Registry.Sprite("quality-symbol")
		// Draw Time cost and quality indicator
		qualitySymbol.Draw(screen, int(c.cardState.Quality), ganim8.DrawOpts(float64(c.left), float64(c.top)))
		for i := range 3 - int(c.cardState.Quality) {
			// TODO Center the time cost
			timeSymbol.Draw(screen, 1, ganim8.DrawOpts(float64(c.left + (i * 32) + 20), float64(c.top + 54) ))
		}

	} else if c.cardState.CardID.IsTool() {
		brokenSymbol := assets.Registry.Sprite("broken-symbol")
		useSymbol := assets.Registry.Sprite("use-symbol")
		qualitySymbol := assets.Registry.Sprite("quality-symbol")
		// Draw Time cost and Remaining uses and check broken
		qualitySymbol.Draw(screen, int(c.cardState.Quality), ganim8.DrawOpts(float64(c.left), float64(c.top)))
		if c.cardState.UsesLeft == 0 {
			brokenSymbol.Draw(screen, 0, ganim8.DrawOpts(float64(c.left), float64(c.top) ))
		} else {
			topShift := 0
			for i := range c.cardState.UsesLeft {
				if i % 2 != 0 {
					topShift = 16
				} else {
					topShift = 0
				}
				useSymbol.Draw(screen, 0, ganim8.DrawOpts(float64(c.left + (i * 16) + 8 ), float64(c.top + 48 + topShift ) ))
			}
		}

		for i := range 2 - int(c.cardState.Quality) {
			timeSymbol.Draw(screen, 1, ganim8.DrawOpts(float64(c.left + 86), float64(c.top + 36 + (i * 32)) ))
		}

	} else if c.cardState.CardID.IsHelper() {
		// Draw Money cost and favorite Tool
		if c.cardState.FavoriteTool == state.ToolGlue {
			c.cardState.FavoriteTool = state.ToolDrill
		}
		favIcon := assets.Registry.Sprite(c.cardState.FavoriteTool.IconName())
		favIcon.Draw(screen, 1, ganim8.DrawOpts(float64(c.left + 71), float64(c.top + 103)))
		r := assets.Registry.DefaultTextRenderer(24)
		r.SetAlign(etxt.YCenter, etxt.XCenter)
		r.SetColor(color.RGBA{0, 200, 50, 255})
		r.SetTarget(screen)
		r.Draw(fmt.Sprintf("$%.2f", c.cardState.MoneyCost), c.left + 64, c.top + 60)
	}
}
