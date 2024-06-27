package scene

import (
	"dbdb/components"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/joelschutz/stagehand"
)

type storeResults struct {
	continueButton  *components.Button
	*BaseScene
}

func CreateStoreResults(width, height int) stagehand.Scene[*State] {
	// Show contents of deck, left over money
	sr := &storeResults{
		BaseScene: NewBaseWithFill(width, height, color.Black), 
	}

	return sr
}

func (ss *storeResults) Draw(screen *ebiten.Image) {
	ebitenutil.DebugPrint(screen, "Store results")
	ss.DrawScene(screen)
	ss.continueButton.Draw(screen)
	ss.DrawIndicators(screen)
}

func (sr *storeResults) Load(s *State, controller stagehand.SceneController[*State]) {
	sr.continueButton = components.NewButton("To the deck...", 320, 272, &s.Controls, func() {
		sr.SceneManager.SwitchWithTransition(SceneMap[BuildingPhase], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
}
