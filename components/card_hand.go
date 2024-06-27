package components

import (
	"dbdb/state"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	trashCombo    color.Color = color.RGBA{255, 127, 50, 200}
	playableCombo color.Color = color.RGBA{127, 255, 50, 200}
)


func NewCardHand(s *state.GlobalState) *CardHand {
	return &CardHand{
		state : s,
		ControlHandler: &s.Controls,
	}
}

type CardHand struct {
	state *state.GlobalState
	*state.ControlHandler
}

func (h *CardHand) Draw(screen *ebiten.Image) {
	hand := h.state.GetHand()
	cards := make([]*Card, len(hand))
	for i := range hand {
		cards[i] = NewCard(hand[i], 32 + ( i * 110 ), 230)
	}

	cardOver := -1

	for i := range cards {
		if IsMouseover(cards[i].left, cards[i].top, 128, 192, false) {
			cardOver = i
			if h.LeftClick {
				h.state.ToggleSelectCard(i)
			}
			break
		}
	}
	selectColor := trashCombo
	if h.state.SelectedCardsPlayable {
		selectColor = playableCombo
	}

	for _, card := range cards {
		if card.cardState.Selected {
			card.highlight = selectColor
		} else {
			card.highlight = nil
		}
	}

	for i := range cards {
		if i != cardOver {
			cards[i].Draw(screen)
		}
	}
	if cardOver != -1 {
		cards[cardOver].Draw(screen)
	}
	switch {
	case h.Key1:
		h.state.ToggleSelectCard(0)
	case h.Key2:
		h.state.ToggleSelectCard(1)
	case h.Key3:
		h.state.ToggleSelectCard(2)
	case h.Key4:
		h.state.ToggleSelectCard(3)
	case h.Key5:
		h.state.ToggleSelectCard(4)
	}
}
