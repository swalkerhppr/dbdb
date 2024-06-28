package scene

import (
	"dbdb/assets"
	"dbdb/components"
	"dbdb/state"
	"fmt"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type State = state.GlobalState

type BaseScene struct {
	State        *State
	SceneManager *stagehand.SceneManager[*State]
	Background   *ebiten.Image
	Alert        *components.Alert

	// Clear the display after each render
	ClearDisplay bool

	Stop error

	Opts *ebiten.DrawImageOptions
}

func NewBaseWithBG(width, height int, bgName string) *BaseScene {
	return &BaseScene{
		Background : assets.Registry.Image(bgName),
		Opts       : ScaleOpts(),
	}
}

func NewBaseWithFill(width, height int, color color.Color) *BaseScene {
	bg := ebiten.NewImage(width, height)
	bg.Fill(color)
	return &BaseScene{
		Background : bg,
		Opts       : ScaleOpts(),
	}
}

func ScaleOpts() *ebiten.DrawImageOptions {
	// images are based on a 640x480 image size
	//opts := &ebiten.DrawImageOptions{}
	//actW, actH := ebiten.WindowSize()
	//opts.GeoM.Scale(float64(actW) / 640, float64(actH) / 480)
	return nil
}

func (b *BaseScene) AdjustAlertPosition(left, top int) {
	b.Alert.SetPosition(left, top)
}

func (b *BaseScene) DrawScene(screen *ebiten.Image) {
	screen.DrawImage(b.Background, b.Opts)
	if b.State.ShowAlert {
		log.Printf("Showing alert: %s", b.State.AlertText)
		b.Alert.SetText(b.State.AlertText)
		b.Alert.Show()
		b.State.ShowAlert = false
	}
	b.Alert.Draw(screen)
}

func (b *BaseScene) DrawIndicators(screen *ebiten.Image) {
	components.NewIndicator("time-symbol", fmt.Sprintf("x%d", b.State.TimeLeft), 400, 435).Draw(screen)
	components.NewIndicator("money-symbol", fmt.Sprintf("%d", int(b.State.MoneyLeft)), 10, 5).Draw(screen)
	components.NewIndicator("card-symbol", fmt.Sprintf(" %d/%d", b.State.CardsLeftInDeck(), len(b.State.Deck)), 315, 435).Draw(screen)

	if b.State.Controls.Inactive {
		b.State.ShowHint = true
	} 
}

// Draw implements stagehand.Scene
func (b *BaseScene) Draw(screen *ebiten.Image) {
	panic("unimplemented!")
}

// Layout implements stagehand.Scene.
func (b *BaseScene) Layout(outsideWidth, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

// Load implements stagehand.Scene.
func (b *BaseScene) Load(s *State, controller stagehand.SceneController[*State]) {
	log.Printf("Loading State: %+v", s)
	b.State = s
	b.Alert = components.NewAlert(&s.Controls)
	b.SceneManager = controller.(*stagehand.SceneManager[*State])
}

// Unload implements stagehand.Scene.
func (b *BaseScene) Unload() *State {
	return b.State
}

// Update implements stagehand.Scene.
func (b *BaseScene) Update() error {
	b.State.Controls.Update()
	return b.Stop
}
