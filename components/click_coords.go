package components

import (
	"fmt"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)


func NewClickCoords() *ClickCoords {
	return &ClickCoords{}
}

type ClickCoords struct {
	coordString string
	clickedString string
}

func (p *ClickCoords) Draw(screen *ebiten.Image) {
	x, y := ebiten.CursorPosition()
	p.coordString = fmt.Sprintf("(%d, %d)", x, y)
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		p.clickedString += fmt.Sprintf(" clicked (%d, %d)", x, y)
	}
	ebitenutil.DebugPrint(screen, p.coordString + p.clickedString) 

}
