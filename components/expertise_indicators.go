package components

import (
	"dbdb/state"

	"github.com/hajimehoshi/ebiten/v2"
)

type ExpertiseIndicators struct {
	s *state.GlobalState
}

func NewExpertiseIndicators(s *state.GlobalState) *ExpertiseIndicators {
	return &ExpertiseIndicators{
		s: s,
	}
}

func (e *ExpertiseIndicators) Draw(screen *ebiten.Image) {
	row := 0
	col := 0
	for i := state.ExpertiseTradesman; i < state.EmptyCardSlot; i = i << 1{
		if e.s.ActiveExpertise & i != 0 {
			NewIndicator(i.SymbolName(), "", 10 + (col * 32), 80 + (row * 32)).Draw(screen)
		}
		col++
		if col == 4 {
			col = 0
			row++
		}
	}
	row = 0
	col = 0
	for i := state.ExpertiseTradesman; i < state.EmptyCardSlot; i = i << 1{
		if e.s.ActiveExpertise & i != 0 {
			NewPopoverArea(10 + (col * 32), 64 + (row * 32), 32, 32,
				NewTextBox(i.HelpDescription(), 15, 42 + (col *32), 80 + (row * 32), 288, 50),
				nil,
			).Draw(screen)
		}
		col++
		if col == 4 {
			col = 0
			row++
		}
	}
}
