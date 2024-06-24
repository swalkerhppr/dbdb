package scene

import (
	"dbdb/components"
	"dbdb/state"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/joelschutz/stagehand"
)

type storeShop struct {
	plankCard       *state.CardState
	boardCard       *state.CardState
	nailCard        *state.CardState
	screwCard       *state.CardState

	toolBought int

	continueButton  *components.Button
	repairButton    *components.Button
	*BaseScene
}

func CreateStoreShop(width, height int) stagehand.Scene[*State] {
	ss := &storeShop{
		BaseScene: NewBaseWithFill(width, height, color.Black), 
	}
	ss.AdjustAlertPosition(230, 250)

	ss.continueButton = components.NewButton("Finish", 483, 320, func() {
		ss.SceneManager.SwitchWithTransition(SceneMap[StoreResults], stagehand.NewDurationTimedSlideTransition[*State](stagehand.RightToLeft, time.Millisecond * 500))
	})
	ss.repairButton = components.NewButton("Repair Tools", 483, 360, func() {
		for _, c := range ss.State.Deck {
			if c.UsesLeft < 4 {
				total := float32(4 - c.UsesLeft) * 50
				if total < ss.State.MoneyLeft {
					ss.State.MoneyLeft -= total
					c.UsesLeft = 4
				} else {
					break
				}
			}
		}
	})
	return ss
}

func (ss *storeShop) Load(s *State, controller stagehand.SceneController[*State]) {
	ss.plankCard = &state.CardState{
		CardID:       state.MaterialPlank,
		Quality:      s.ChosenStore.StoreQuality,
	}
	ss.boardCard = &state.CardState{
		CardID:       state.MaterialBoard,
		Quality:      s.ChosenStore.StoreQuality,
	}
	ss.nailCard = &state.CardState{
		CardID:       state.MaterialNail,
		Quality:      s.ChosenStore.StoreQuality,
	}
	ss.screwCard = &state.CardState{
		CardID:       state.MaterialScrew,
		Quality:      s.ChosenStore.StoreQuality,
	}
	ss.BaseScene.Load(s, controller)
}

func (ss *storeShop) Draw(screen *ebiten.Image) {
	ss.DrawScene(screen)
	ss.continueButton.Draw(screen)
	ss.repairButton.Draw(screen)
	components.NewIndicatorWithColor("money-symbol", fmt.Sprintf("%d", int(ss.State.MoneyLeft)), 260, 10, color.RGBA{0, 200, 50, 255}).Draw(screen)
	components.NewCard(ss.plankCard, 40, 50).Draw(screen)
	components.NewCard(ss.boardCard, 170, 50).Draw(screen)
	components.NewCard(ss.nailCard, 300, 50).Draw(screen)
	components.NewCard(ss.screwCard, 430, 50).Draw(screen)

	var tool *state.CardState
	if ss.toolBought < 2 {
		tool = ss.State.ChosenStore.AvailableTools[ss.toolBought]
		components.NewCard(tool, 100, 250).Draw(screen)
	}

	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButton0) {
		if components.IsMouseover(40, 50, 128, 198, false) {
			ss.State.ChosenStore.PlankStock = ss.State.BuyCard(ss.plankCard, ss.State.ChosenStore.PlankPrice, ss.State.ChosenStore.PlankStock)
		} else if components.IsMouseover(170, 50, 128, 198, false) {
			ss.State.ChosenStore.BoardStock = ss.State.BuyCard(ss.boardCard, ss.State.ChosenStore.BoardPrice, ss.State.ChosenStore.BoardStock)
		} else if components.IsMouseover(300, 50, 128, 198, false) {
			ss.State.ChosenStore.NailStock = ss.State.BuyCard(ss.nailCard, ss.State.ChosenStore.NailPrice, ss.State.ChosenStore.NailStock)
		} else if components.IsMouseover(430, 50, 128, 198, false) {
			ss.State.ChosenStore.ScrewStock = ss.State.BuyCard(ss.screwCard, ss.State.ChosenStore.ScrewPrice, ss.State.ChosenStore.ScrewStock)
		} else if components.IsMouseover(100, 250, 128, 198, false) && tool != nil {
				switch tool.Quality {
				case state.OneStar:
					// Shouldn't ever happen
				case state.TwoStar:
					ss.State.BuyCard(tool, 100, 1)
					ss.toolBought++
				case state.ThreeStar:
					ss.State.BuyCard(tool, 200, 1)
					ss.toolBought++
			}
		}
	}
}
