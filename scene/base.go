package scene

import (
	"dbdb/assets"
	"dbdb/state"
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

	// Clear the display after each render
	ClearDisplay bool

	Stop error

	Opts *ebiten.DrawImageOptions
}

func NewBaseWithBG(width, height int, bgName string) *BaseScene {
	return &BaseScene{
		Background: assets.Registry.Image(bgName),
		Opts : ScaleOpts(),
	}
}

func NewBaseWithFill(width, height int, color color.Color) *BaseScene {
	bg := ebiten.NewImage(width, height)
	bg.Fill(color)
	return &BaseScene{
		Background: bg,
		Opts : ScaleOpts(),
	}
}

func ScaleOpts() *ebiten.DrawImageOptions {
	// images are based on a 640x480 image size
	opts := &ebiten.DrawImageOptions{}
	actW, actH := ebiten.WindowSize()
	opts.GeoM.Scale(float64(actW) / 640, float64(actH) / 480)
	return opts
}

func (b *BaseScene) DrawScene(screen *ebiten.Image) {
	screen.DrawImage(b.Background, b.Opts)
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
	s.ShuffleCards(0, 4)
	b.State = s
	b.SceneManager = controller.(*stagehand.SceneManager[*State])
}

// Unload implements stagehand.Scene.
func (b *BaseScene) Unload() *State {
	return b.State
}

// Update implements stagehand.Scene.
func (b *BaseScene) Update() error {
	return b.Stop
}
