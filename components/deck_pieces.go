package components

import (
	"dbdb/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/ganim8/v2"
)

type DeckPieces struct {
	numPlanksBuilt int
	numBoardsBuilt int
	sprite *ganim8.Sprite
}

func NewDeckPieces(planks, boards int) *DeckPieces {
	return &DeckPieces{
		numPlanksBuilt: planks,
		numBoardsBuilt: boards,
		sprite: assets.Registry.Sprite("deck-tiles"),
	}
}

func (d *DeckPieces) Draw(screen *ebiten.Image) {
	// | 0:lt      1:lt-con   2:rt-con   3:rt-stop |
	// | 4:l-plank 5:board    6:board    7:r-plank |
	// | 8:l-gnd-1 9:l-gnd-2 10:r-gnd-1 11:r-gnd2  |

	switch d.numPlanksBuilt {
	case 5:
		d.sprite.Draw(screen, 3,  ganim8.DrawOpts(614, 103))
		d.sprite.Draw(screen, 7,  ganim8.DrawOpts(614, 119))
		d.sprite.Draw(screen, 7,  ganim8.DrawOpts(614, 135))
		d.sprite.Draw(screen, 11, ganim8.DrawOpts(614, 151))
		fallthrough
	case 4:
		d.sprite.Draw(screen, 2,  ganim8.DrawOpts(550, 103))
		d.sprite.Draw(screen, 7,  ganim8.DrawOpts(550, 119))
		d.sprite.Draw(screen, 7,  ganim8.DrawOpts(550, 135))
		d.sprite.Draw(screen, 10, ganim8.DrawOpts(550, 151))
		fallthrough
	case 3:
		d.sprite.Draw(screen, 2,  ganim8.DrawOpts(502, 103))
		d.sprite.Draw(screen, 7,  ganim8.DrawOpts(502, 119))
		d.sprite.Draw(screen, 7,  ganim8.DrawOpts(502, 135))
		d.sprite.Draw(screen, 7, ganim8.DrawOpts(502, 151))
		fallthrough
	case 2:
		d.sprite.Draw(screen, 1,  ganim8.DrawOpts(454, 103))
		d.sprite.Draw(screen, 4,  ganim8.DrawOpts(454, 119))
		d.sprite.Draw(screen, 4,  ganim8.DrawOpts(454, 135))
		d.sprite.Draw(screen, 4, ganim8.DrawOpts(454, 151))
		fallthrough
	case 1:
		d.sprite.Draw(screen, 1,  ganim8.DrawOpts(406, 103))
		d.sprite.Draw(screen, 4,  ganim8.DrawOpts(406, 119))
		d.sprite.Draw(screen, 4,  ganim8.DrawOpts(406, 135))
		d.sprite.Draw(screen, 9, ganim8.DrawOpts(406, 151))
	}

	switch d.numBoardsBuilt {
	case 5:
		d.sprite.Draw(screen, 5,  ganim8.DrawOpts(582, 103))
		d.sprite.Draw(screen, 6,  ganim8.DrawOpts(598, 103))
		d.sprite.Draw(screen, 5,  ganim8.DrawOpts(566, 103))
		d.sprite.Draw(screen, 6,  ganim8.DrawOpts(582, 103))
		fallthrough
	case 4:
		d.sprite.Draw(screen, 5,  ganim8.DrawOpts(518, 103))
		d.sprite.Draw(screen, 6,  ganim8.DrawOpts(534, 103))
		fallthrough
	case 3:
		d.sprite.Draw(screen, 5,  ganim8.DrawOpts(470, 103))
		d.sprite.Draw(screen, 6,  ganim8.DrawOpts(486, 103))
		fallthrough
	case 2:
		d.sprite.Draw(screen, 5,  ganim8.DrawOpts(422, 103))
		d.sprite.Draw(screen, 6,  ganim8.DrawOpts(438, 103))
		fallthrough
	case 1:
		d.sprite.Draw(screen, 5,  ganim8.DrawOpts(374, 103))
		d.sprite.Draw(screen, 6,  ganim8.DrawOpts(390, 103))
	}
}
