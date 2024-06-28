package scene

import (
	"dbdb/components"
	"dbdb/state"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type buildingPhase struct {
	hand          *components.CardHand
	discardButton *components.Button
	playButton    *components.Button
	endDayButton  *components.Button
	cardSlots     *components.CardSlots
	*BaseScene
	figures       []*components.Figure
}

func CreateBuildingPhase(width, height int) stagehand.Scene[*State] {
	sp := &buildingPhase{
		BaseScene: NewBaseWithBG(width, height, "deckbuilding.png"),
	}
	return sp
}

func (s *buildingPhase) Draw(screen *ebiten.Image) {
	if s.State.CardsLeftInDeck() < 5 {
		s.discardButton.SetText("New Hand (-2)")
	} else {
		s.discardButton.SetText("New Hand")
	}

	s.DrawScene(screen)
	s.hand.Draw(screen)
	s.discardButton.Draw(screen)
	s.playButton.Draw(screen)
	s.endDayButton.Draw(screen)
	if s.State.GameWon {
		s.SceneManager.SwitchTo(SceneMap[GameResults])
	}
	s.DrawIndicators(screen)
	s.cardSlots.Draw(screen)
	components.NewDeckIndicators(s.State, 11, 214, 0).Draw(screen)
	components.NewDeckPieces(s.State.PlankPartsBuilt, s.State.BoardPartsBuilt).Draw(screen)
	components.NewExpertiseIndicators(s.State).Draw(screen)
	if len(s.State.ActiveHelpers) > len(s.figures){
		s.figures = append(s.figures, components.NewFigure(s.State.ActiveHelpers[len(s.State.ActiveHelpers)-1].Card.CardID.RawName(), 396 + (len(s.figures) * 20 ), 110, 1.2))
	}
	for _, f := range s.figures {
		f.Draw(screen)
	}
	if s.State.Controls.KeyEnter {
		s.endDayButton.OnClick()
	}
	if s.State.Controls.KeyTab {
		s.discardButton.OnClick()
	}
	if s.State.Controls.KeySpace {
		s.playButton.OnClick()
	}
}

func (p *buildingPhase) Load(s *State, controller stagehand.SceneController[*State]) {
	s.ClearSelectedCards()
	s.Phase = state.BuildPhase
	p.hand = components.NewCardHand(s)
	p.cardSlots = components.NewCardSlots(s, 10, 186)
	s.ShuffleCards(0, 4)
	p.playButton = components.NewButton("Play", 96, 445, &s.Controls, func() {
		p.State.PlaySelected()
	})
	p.discardButton = components.NewButton("New Hand", 235, 445, &s.Controls, func() {
		if s.DiscardHand() && s.TimeLeft > 0 {
			s.TimeLeft -= 2
		}
	})
	p.endDayButton = components.NewButton("End Day", 532, 445, &s.Controls, func() {
		s.Day++
		if s.Day == s.MaxDays {
			p.SceneManager.SwitchWithTransition(SceneMap[GameResults], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
		} else {
			p.SceneManager.SwitchWithTransition(SceneMap[DayResults], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
		}
	})
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

