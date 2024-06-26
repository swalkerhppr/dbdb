package scene

import (
	"dbdb/components"
	"dbdb/state"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type buildingPhase struct {
	hand *components.CardHand
	discardButton *components.Button
	playButton    *components.Button
	endDayButton  *components.Button
	*BaseScene
	figures       []*components.Figure
}

func CreateBuildingPhase(width, height int) stagehand.Scene[*State] {
	sp := &buildingPhase{
		BaseScene: NewBaseWithBG(width, height, "deckbuilding.png"),
	}
	sp.playButton = components.NewButton("Play", 96, 445, func() {
		sp.State.PlaySelected()
	})
	sp.discardButton = components.NewButton("New Hand", 235, 445, func() {
		if sp.State.DiscardHand() && sp.State.TimeLeft > 0 {
			sp.State.TimeLeft -= 2
		}
	})
	sp.endDayButton = components.NewButton("End Day", 532, 445, func() {
		sp.State.Day++
		if sp.State.Day == sp.State.MaxDays {
			sp.SceneManager.SwitchWithTransition(SceneMap[GameResults], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
		} else {
			sp.SceneManager.SwitchWithTransition(SceneMap[DayResults], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
		}
	})
	return sp
}

func (s *buildingPhase) Draw(screen *ebiten.Image) {
	s.DrawScene(screen)
	s.hand.Draw(screen)
	s.discardButton.Draw(screen)
	s.playButton.Draw(screen)
	s.endDayButton.Draw(screen)
	if s.State.GameWon {
		s.SceneManager.SwitchTo(SceneMap[GameResults])
	}
	s.DrawIndicators(screen)
	components.NewCardSlots(s.State, 10, 186).Draw(screen)
	components.NewDeckIndicators(s.State, 11, 214, 0).Draw(screen)
	components.NewDeckPieces(s.State.PlankPartsBuilt, s.State.BoardPartsBuilt).Draw(screen)
	components.NewExpertiseIndicators(s.State).Draw(screen)
	if len(s.State.ActiveHelpers) > len(s.figures){
		s.figures = append(s.figures, components.NewFigure(s.State.ActiveHelpers[len(s.State.ActiveHelpers)-1].Card.CardID.RawName(), 396 + (len(s.figures) * 20 ), 110, 1.2))
	}
	for _, f := range s.figures {
		f.Draw(screen)
	}
}

func (p *buildingPhase) Load(s *State, controller stagehand.SceneController[*State]) {
	s.ClearSelectedCards()
	s.Phase = state.BuildPhase
	p.hand = components.NewCardHand(s)
	s.ShuffleCards(0, 4)
	p.BaseScene.Load(s, controller)
}

func (b *buildingPhase) Unload() *State {
	for _, c := range b.State.HeldCards {
		if !c.CardID.IsEmptyCardSlot() {
			b.State.AddCard(c)
		}
	}
	for _, h := range b.State.ActiveHelpers {
		b.State.AddCard(h.Card)
	}
	b.figures = b.figures[0:0]
	b.State.ActiveExpertise = 0
	b.State.HeldCards = b.State.HeldCards[0:0]
	b.State.ActiveHelpers = b.State.ActiveHelpers[0:0]
	return b.State
}

