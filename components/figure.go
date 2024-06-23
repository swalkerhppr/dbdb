package components

import (
	"dbdb/assets"
	"time"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/yohamta/ganim8/v2"
)

func NewRandomFigure(x, y int, scale float64) *Figure {
	npcName := []string{ 
		"colo",
		"hood",
		"curt",
		"guru",
		"oarm",
		"sass",
		"shop",
	}[rand.Intn(7)]
	return NewFigure(npcName, x, y, scale)
}

func NewRandomEmployeeFigure(x, y int, scale float64) *Figure {
	npcName := []string{ 
		"hood",
		"curt",
		"guru",
		"oarm",
	}[rand.Intn(4)]
	return NewFigure(npcName, x, y, scale)
}

func NewRandomNeighborFigure(x, y int, scale float64) *Figure {
	npcName := []string{ 
		"colo",
		"sass",
		"shop",
	}[rand.Intn(3)]
	return NewFigure(npcName, x, y, scale)
}

func NewFigure(name string, x, y int, scale float64) *Figure {
	return &Figure{
		name : name,
		idleAnim : ganim8.NewAnimation(assets.Registry.Sprite(name + "-fig-idle"), time.Millisecond * 200),
		walkAnim : ganim8.NewAnimation(assets.Registry.Sprite(name + "-fig-walk"), time.Millisecond * 200),
		top : y,
		left: x,
		scale: scale,
	}
}

type Figure struct {
	name     string
	top      int
	left     int
	scale    float64

	idleAnim *ganim8.Animation 
	walkAnim *ganim8.Animation 

	walking bool
	hFlip   bool
}

func (f *Figure) CardName() string {
	return f.name + "-card"
}

func (f *Figure) Draw(screen *ebiten.Image) {
	if f.walking {
		if f.hFlip {
			f.walkAnim.Sprite().FlipH()
			f.hFlip = false
		}
		f.walkAnim.Draw(screen, ganim8.DrawOpts(float64(f.left), float64(f.top), 0, f.scale, f.scale))
		f.walkAnim.Update()
	} else {
		f.idleAnim.Draw(screen, ganim8.DrawOpts(float64(f.left), float64(f.top), 0, f.scale, f.scale))
		f.idleAnim.Update()
	}
}
