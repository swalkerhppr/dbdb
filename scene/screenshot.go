//go:build !wasm
package scene

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)


func takeScreenshot(path string, screen *ebiten.Image) {
	log.Println("Not implemented for windows yet")
}
