package scene

import (
	"dbdb/components"
	"dbdb/state"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
)

type buildingPhase struct {
	hand *components.CardHand
	discardButton *components.Button
	playButton    *components.Button
	*BaseScene
}

func CreateBuildingPhase(width, height int) stagehand.Scene[*State] {
	sp := &buildingPhase{
		BaseScene: NewBaseWithBG(width, height, "deckbuilding.png"),
	}
	sp.discardButton = components.NewButton("Discard", 500, 455, func() {
		sp.BaseScene.State.DiscardHand()
	})
	sp.playButton = components.NewButton("Play", 128, 455, func() {
		sp.BaseScene.State.PlaySelected()
	})
	return sp
}

func (s *buildingPhase) Draw(screen *ebiten.Image) {
	s.DrawScene(screen)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton1) {
		s.SceneManager.SwitchTo(SceneMap[MainMenu])
	}
	s.hand.Draw(screen)
	s.discardButton.Draw(screen)
	s.playButton.Draw(screen)
}

func (p *buildingPhase) Load(s *State, controller stagehand.SceneController[*State]) {
	s.ClearSelectedCards()
	s.Phase = state.BuildPhase
	p.hand = components.NewCardHand(s)
	p.BaseScene.Load(s, controller)
}
