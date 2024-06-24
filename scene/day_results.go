package scene

import (
	"dbdb/components"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joelschutz/stagehand"
)

type dayResults struct {
	continueButton  *components.Button
	*BaseScene
}

func CreateDayResults(width, height int) stagehand.Scene[*State] {
	// Show current state: cards in deck, left over money, deck progress
	// If day > 5 you lose. If all requirements are met you win!
	sr := &dayResults{
		BaseScene: NewBaseWithFill(width, height, color.Black), 
	}

	sr.continueButton = components.NewButton("Next Day", 320, 272, func() {
		sr.SceneManager.SwitchWithTransition(SceneMap[ChooseStore], stagehand.NewDurationTimedSlideTransition[*State](stagehand.RightToLeft, time.Millisecond * 500))
	})
	return sr
}

func (ss *dayResults) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Day results")
	ss.DrawScene(screen)
	ss.continueButton.Draw(screen)
	ss.DrawIndicators(screen)
}
