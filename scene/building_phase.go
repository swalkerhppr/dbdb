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
			sp.State.TimeLeft--
		}
	})
	sp.endDayButton = components.NewButton("End Day", 532, 445, func() {
		sp.SceneManager.SwitchWithTransition(SceneMap[DayResults], stagehand.NewDurationTimedSlideTransition[*State](stagehand.TopToBottom, time.Millisecond * 500))
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
		// TODO win screen
		s.SceneManager.SwitchTo(SceneMap[MainMenu])
	}
	s.DrawIndicators(screen)
	components.NewCardSlots(s.State, 10, 186).Draw(screen)
}

func (p *buildingPhase) Load(s *State, controller stagehand.SceneController[*State]) {
	s.ClearSelectedCards()
	s.Phase = state.BuildPhase
	p.hand = components.NewCardHand(s)
	p.BaseScene.Load(s, controller)
}
