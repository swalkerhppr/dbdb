package state

import (
	"log"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type ControlHandler struct {
	LeftClick  bool
	RightClick bool
	LeftPress  bool
	RightPress bool

	Key1       bool
	Key2       bool
	Key3       bool
	Key4       bool
	Key5       bool
	KeySpace   bool
	KeyTab     bool
	KeyEnter   bool
	KeyEscape  bool
	KeyM       bool

	Any        bool
	Inactive   bool
	LastAction time.Time
}

func (c *ControlHandler) Update() error {
	c.LeftClick  = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	c.RightClick = inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)
	c.LeftPress  = !c.LeftClick  && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft)
	c.RightPress = !c.RightClick && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight)

	c.Key1     = inpututil.IsKeyJustReleased(ebiten.Key1)
	c.Key2     = inpututil.IsKeyJustReleased(ebiten.Key2)
	c.Key3     = inpututil.IsKeyJustReleased(ebiten.Key3)
	c.Key4     = inpututil.IsKeyJustReleased(ebiten.Key4)
	c.Key5     = inpututil.IsKeyJustReleased(ebiten.Key5)
	c.KeyEnter = inpututil.IsKeyJustReleased(ebiten.KeyEnter)
	c.KeyTab   = inpututil.IsKeyJustReleased(ebiten.KeyTab)
	c.KeySpace = inpututil.IsKeyJustReleased(ebiten.KeySpace)
	c.KeyEscape = inpututil.IsKeyJustReleased(ebiten.KeyEscape)
	c.KeyM = inpututil.IsKeyJustReleased(ebiten.KeyM)

	c.Any = c.Key1 || c.Key2 || c.Key3 || c.Key4 || c.Key5 || c.KeySpace || c.KeyTab || c.KeyEnter || c.LeftClick || c.RightClick
	if c.Any {
		c.LastAction = time.Now()
		c.Inactive = false
	} else if !c.Inactive && time.Now().Sub(c.LastAction) >= time.Second * 15 {
		log.Println("Inactivity timeout")
		c.Inactive = true
	}
	return nil
}
