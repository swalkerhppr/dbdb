package scene

import (
	"dbdb/components"
	//"errors"
	"log"
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

	mm.startButton = components.NewButton("Start", 320, 176, func() {
		mm.BaseScene.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedSlideTransition[*State](stagehand.RightToLeft, time.Millisecond * 500))
	})
	mm.loadButton = components.NewButton("About", 320, 224, func() {
		log.Println("load Clicked!")
	})
	//mm.exitButton = components.NewButton("Exit", 320, 272, func() {
	//	mm.BaseScene.Stop = errors.New("User Exit")
	//})
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
