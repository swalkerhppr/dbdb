package components

import (
	"dbdb/assets"
	"dbdb/state"
	"math/rand"
	"time"

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

func (f *Figure) CardID() state.CardID {
	switch f.name {
	case "colo":
		return state.HelperColo
	case "hood":
		return state.HelperHood
	case "curt":
		return state.HelperCurt
	case "guru":
		return state.HelperGuru
	case "oarm":
		return state.HelperOarm
	case "sass":
		return state.HelperSass
	case "shop":
		return state.HelperShop
	}
	return state.CardID(0)
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
