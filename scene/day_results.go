package scene

import (
	"dbdb/components"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type dayResults struct {
	continueButton *components.Button
	scoreText string
	*BaseScene
}

func CreateDayResults(width, height int) stagehand.Scene[*State] {
	// If day > 5 you lose. If all requirements are met you win!
	gr := &dayResults{
		BaseScene: NewBaseWithFill(width, height, color.RGBA{R: 0x73, G: 0x7b, B: 0x9a, A: 0xFF}), 
	}
	return gr
}

func (gr *dayResults) Load(s *State, controller stagehand.SceneController[*State]) {
	gr.continueButton = components.NewButton("Back to it", 320, 404, &s.Controls, func() {
		gr.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	gr.scoreText = createScoreText(s)
	gr.BaseScene.Load(s, controller)
}

func (gr *dayResults) Draw(screen *ebiten.Image) {
	gr.DrawScene(screen)
	text := "Pheew...\nAll decked out.\nTomorrow is another day.\n\n\n\n\n\n\n\n\n"
	components.NewTextBox(text, 24, 18, 16, 608, 448).Draw(screen)
	components.NewDeckIndicators(gr.State, 8, 164, 163).Draw(screen)
	components.NewTextBox(gr.scoreText, 18, 191, 240, 256, 128).Draw(screen)
	gr.continueButton.Draw(screen)
}
