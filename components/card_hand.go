package components

import (
	"dbdb/state"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

var (
	trashCombo    color.Color = color.RGBA{255, 127, 50, 200}
	playableCombo color.Color = color.RGBA{127, 255, 50, 200}
)


func NewCardHand(s *state.GlobalState) *CardHand {
	return &CardHand{
		state : s,
	}
}

type CardHand struct {
	state *state.GlobalState
}

func (h *CardHand) Draw(screen *ebiten.Image) {
	cards := make([]*Card, 5)
	for i, cs := range h.state.GetHand(){
		cards[i] = NewCard(cs, 32 + ( i * 110 ), 230)
	}

	clicked := inpututil.IsMouseButtonJustPressed(ebiten.MouseButton0)
	cardOver := -1

	for i := range cards {
		if IsMouseover(cards[i].left, cards[i].top, 128, 192, false) {
			cardOver = i
			if clicked {
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
}
