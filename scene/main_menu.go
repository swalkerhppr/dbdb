package scene

import (
	"dbdb/components"
	"dbdb/state"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type mainMenu struct {
	startButton  *components.Button
	loadButton   *components.Button
	exitButton   *components.Button
	npcFigure    *components.Figure
	*BaseScene
}

func CreateMainMenu(width, height int) stagehand.Scene[*State] {
	mm := &mainMenu{
		npcFigure: components.NewRandomFigure(250, 330, 1),
		BaseScene: NewBaseWithBG(width, height, "mainmenu.png"), 
	}

	mm.startButton = components.NewButton("Normal", 320, 176, func() {
		mm.State = state.InitialState()
		mm.State.MaxDays = 5
		mm.State.MoneyLeft = 2000
		mm.State.Deck = state.InitialDeck(0)
		mm.BaseScene.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	mm.loadButton = components.NewButton("Hard", 320, 224, func() {
		mm.State = state.InitialState()
		mm.State.MaxDays = 4
		mm.State.MoneyLeft = 1000
		mm.State.Deck = state.InitialDeck(1)
		mm.BaseScene.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	mm.exitButton = components.NewButton("Impossible", 320, 272, func() {
		mm.State = state.InitialState()
		mm.State.MaxDays = 3
		mm.State.MoneyLeft = 500
		mm.State.Deck = state.InitialDeck(2)
		mm.BaseScene.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	return mm
}

func (mm *mainMenu) Load(s *State, controller stagehand.SceneController[*State]) {
	mm.npcFigure = components.NewRandomFigure(250, 330, 1)
	mm.BaseScene.Load(s, controller)
}

func (m *mainMenu) Draw(screen *ebiten.Image) {
	m.DrawScene(screen)

	m.startButton.Draw(screen)
	m.loadButton.Draw(screen)
	m.exitButton.Draw(screen)
	m.npcFigure.Draw(screen)
}
