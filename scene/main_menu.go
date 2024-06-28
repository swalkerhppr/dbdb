package scene

import (
	"dbdb/assets"
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
	muteButton   *components.Button
	*BaseScene
}

func CreateMainMenu(width, height int) stagehand.Scene[*State] {
	mm := &mainMenu{
		npcFigure: components.NewRandomFigure(250, 330, 1),
		BaseScene: NewBaseWithBG(width, height, "mainmenu.png"), 
	}

	return mm
}

func (mm *mainMenu) Load(s *State, controller stagehand.SceneController[*State]) {
	mm.startButton = components.NewButton("Normal", 320, 176, &s.Controls, func() {
		mm.State = state.InitialState()
		mm.State.MaxDays = 5
		mm.State.MoneyLeft = 2000
		mm.State.Deck = state.InitialDeck(0)
		mm.BaseScene.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	mm.loadButton = components.NewButton("Hard", 320, 224, &s.Controls, func() {
		mm.State = state.InitialState()
		mm.State.MaxDays = 4
		mm.State.MoneyLeft = 1000
		mm.State.Deck = state.InitialDeck(1)
		mm.BaseScene.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	mm.exitButton = components.NewButton("Impossible", 320, 272, &s.Controls, func() {
		mm.State = state.InitialState()
		mm.State.MaxDays = 3
		mm.State.MoneyLeft = 1000
		mm.State.Deck = state.InitialDeck(2)
		mm.BaseScene.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	bgm := assets.Registry.Sound("background")
	mm.muteButton = components.NewButton("BGM OFF", 71, 458, &s.Controls, func() {
		switch int(bgm.Volume() * 100) {
		case 30:
			bgm.SetVolume(0)
			mm.muteButton.SetText("BGM ON")
		case 15:
			bgm.SetVolume(0.3)
			mm.muteButton.SetText("BGM OFF")
		case 0:
			bgm.SetVolume(0.15)
			mm.muteButton.SetText("BGM UP")
		default:
			bgm.SetVolume(0)
			mm.muteButton.SetText("BGM ON")
		}
	})
	mm.npcFigure = components.NewRandomFigure(250, 330, 1)
	mm.BaseScene.Load(s, controller)
	mm.Update()
}

func (m *mainMenu) Draw(screen *ebiten.Image) {
	m.DrawScene(screen)

	m.startButton.Draw(screen)
	m.loadButton.Draw(screen)
	m.exitButton.Draw(screen)
	m.muteButton.Draw(screen)
	m.npcFigure.Draw(screen)

	if m.State.Controls.Key1 {
		m.startButton.OnClick()
	}
	if m.State.Controls.Key1 {
		m.loadButton.OnClick()
	}
	if m.State.Controls.Key1 {
		m.exitButton.OnClick()
	}
}
