package main

import (
	"dbdb/assets"
	"dbdb/scene"
	"dbdb/state"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

const (
	screenWidth  = 640
	screenHeight = 480
	//screenWidth  = 960
	//screenHeight = 720
)

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Deck Building Deck Builder")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	assets.InitRegistry()
	scene.InitializeScenes(screenWidth, screenHeight)

	globalState := &state.GlobalState{}

	startAt := scene.MainMenu

	sm := stagehand.NewSceneManager[*scene.State](scene.SceneMap[startAt], globalState)

	assets.Registry.Sound("background").Play()

	if err := ebiten.RunGame(sm); err != nil {
		log.Fatal(err)
	}
}
