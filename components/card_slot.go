package components

import (
	"dbdb/assets"
	"dbdb/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/ganim8/v2"
)


type CardSlots struct {
	globalState   *state.GlobalState
	left   int
	top    int
	*state.ControlHandler
}

func NewCardSlots(g *state.GlobalState, left, top int) *CardSlots {
	return &CardSlots{
		globalState   : g,
		left   : left,
		top    : top,
		ControlHandler: &g.Controls,
	}
}

func (c *CardSlots) Draw(screen *ebiten.Image) {
	var sprite *ganim8.Sprite
	for i, slot := range c.globalState.HeldCards {
		sprite = assets.Registry.Sprite(slot.CardID.IconName())
		frame := 0
		if slot.CardID.IsEmptyCardSlot() {
			frame = 1
		}

		if !slot.CardID.IsEmptyCardSlot() && slot.Selected {
			highlight := trashCombo
			if c.globalState.SelectedCardsPlayable {
				highlight = playableCombo
			}

			hlImage := ebiten.NewImage(32, 32)
			hlImage.Fill(highlight)
			opts := &ebiten.DrawImageOptions{}
			opts.GeoM.Translate(float64(c.left + (i * 32)), float64(c.top))
			screen.DrawImage(hlImage, opts)
		}

		sprite.Draw(screen, frame, ganim8.DrawOpts(float64(c.left + (i * 32)), float64(c.top)))
	}
	for i, slot := range c.globalState.HeldCards {
		if IsMouseover(c.left + (i * 32), c.top, 32, 32, false) {
			if c.LeftClick {
				slot.Selected = !slot.Selected
				c.globalState.UpdateSelectedCardsPlayable()
			}

			if !slot.CardID.IsEmptyCardSlot() {
				NewCard(slot, c.left + ( (i + 1) * 32), c.top - 96).Draw(screen)

			} else if c.globalState.CanSlot(slot){
				if c.RightClick {
					// place the selected card in the slot
					c.globalState.HoldCard(slot)
				}
			} 
		}
	}
}
