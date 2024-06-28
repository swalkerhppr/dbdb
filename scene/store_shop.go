package scene

import (
	"dbdb/assets"
	"dbdb/components"
	"dbdb/state"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
	"github.com/tinne26/etxt"
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
		BaseScene: NewBaseWithBG(width, height, "store_shop.png"), 
	}

	return ss
}

func (ss *storeShop) Load(s *State, controller stagehand.SceneController[*State]) {
	ss.continueButton = components.NewButton("Get Building!", 483, 320, &s.Controls, func() {
		ss.SceneManager.SwitchWithTransition(SceneMap[BuildingPhase], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	ss.repairButton = components.NewButton("Repair Tools", 483, 360, &s.Controls, func() {
		for _, c := range ss.State.Deck {
			if c.CardID.IsTool() && c.UsesLeft < 4 {
				total := float32(4 - c.UsesLeft) * 25
				if total < ss.State.MoneyLeft {
					ss.State.MoneyLeft -= total
					if int(total/2) < ss.State.TimeLeft {
						ss.State.TimeLeft -= int(total / 2)
					}
					c.UsesLeft = 4
				} else {
					break
				}
			}
		}
	})

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
	s.ChosenStore.AvailableTools = make([]*state.CardState, 2)
	toolID := state.RandomToolID()
	s.ChosenStore.AvailableTools[0] = &state.CardState{
		CardID:   toolID,
		Quality:  toolID.ToolQuality(),
		UsesLeft: 4,
	}
	toolID = state.RandomToolID()
	s.ChosenStore.AvailableTools[1] = &state.CardState{
		CardID:   toolID,
		Quality:  toolID.ToolQuality(),
		UsesLeft: 4,
	}
	s.Phase = state.StorePhase
	ss.toolBought = 0

	s.ClearSelectedCards()
	ss.BaseScene.Load(s, controller)
	ss.AdjustAlertPosition(230, 250)
}

func (ss *storeShop) Draw(screen *ebiten.Image) {
	moneyGreen := color.RGBA{0, 200, 50, 255}
	ss.DrawScene(screen)
	ss.continueButton.Draw(screen)
	ss.repairButton.Draw(screen)
	components.NewIndicatorWithColor("money-symbol", fmt.Sprintf("%d", int(ss.State.MoneyLeft)), 260, 10, moneyGreen).Draw(screen)
	components.NewCard(ss.plankCard, 40, 50).Draw(screen)
	components.NewCard(ss.boardCard, 170, 50).Draw(screen)
	components.NewCard(ss.nailCard, 300, 50).Draw(screen)
	components.NewCard(ss.screwCard, 430, 50).Draw(screen)

	components.NewIndicator("time-symbol", fmt.Sprintf("x %d", ss.State.TimeLeft), 400, 435).Draw(screen)

	r := assets.Registry.DefaultTextRenderer(17)
	r.SetAlign(etxt.Top, etxt.Right)
	r.SetTarget(screen)

	r.SetColor(moneyGreen)
	r.Draw(fmt.Sprintf("$%.0f ", ss.State.ChosenStore.PlankPrice), 163, 55)
	r.Draw(fmt.Sprintf("$%.0f ", ss.State.ChosenStore.BoardPrice), 293, 55)
	r.Draw(fmt.Sprintf("$%.0f ", ss.State.ChosenStore.NailPrice),  423, 55)
	r.Draw(fmt.Sprintf("$%.0f ", ss.State.ChosenStore.ScrewPrice), 553, 55)
	
	r.SetAlign(etxt.Bottom, etxt.Right)
	r.SetColor(color.Black)
	r.Draw(fmt.Sprintf("%d", ss.State.ChosenStore.PlankStock), 160, 240)
	r.Draw(fmt.Sprintf("%d", ss.State.ChosenStore.BoardStock), 290, 240)
	r.Draw(fmt.Sprintf("%d", ss.State.ChosenStore.NailStock),  420, 240)
	r.Draw(fmt.Sprintf("%d", ss.State.ChosenStore.ScrewStock), 550, 240)

	var tool *state.CardState
	if ss.toolBought < 2 {
		tool = ss.State.ChosenStore.AvailableTools[ss.toolBought]
		components.NewCard(tool, 100, 250).Draw(screen)

		r.SetAlign(etxt.Top, etxt.Right)
		r.SetColor(moneyGreen)
		if tool.Quality == state.ThreeStar {
			r.Draw("$200", 220, 270)
		} else {
			r.Draw("$100", 220, 270)
		}

		r.SetAlign(etxt.Bottom, etxt.Right)
		r.SetColor(color.Black)
		r.Draw("1", 220, 440)
	}

	components.NewShopDeckIndicators(ss.State, 4, 235, 300).Draw(screen)

	if ss.continueButton.LeftClick {
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
					if ss.State.BuyCard(tool, 100, 1) == 0 {
						ss.toolBought++
					}
				case state.TwoStar:
					if ss.State.BuyCard(tool, 100, 1) == 0 {
						ss.toolBought++
					}
				case state.ThreeStar:
					if ss.State.BuyCard(tool, 200, 1) == 0 {
						ss.toolBought++
					}
			}
		}
	}

	switch {
	case ss.State.Controls.Key1:
		ss.State.ChosenStore.PlankStock = ss.State.BuyCard(ss.plankCard, ss.State.ChosenStore.PlankPrice, ss.State.ChosenStore.PlankStock)

	case ss.State.Controls.Key2:
		ss.State.ChosenStore.BoardStock = ss.State.BuyCard(ss.boardCard, ss.State.ChosenStore.BoardPrice, ss.State.ChosenStore.BoardStock)

	case ss.State.Controls.Key3:
		ss.State.ChosenStore.NailStock = ss.State.BuyCard(ss.nailCard, ss.State.ChosenStore.NailPrice, ss.State.ChosenStore.NailStock)

	case ss.State.Controls.Key4:
		ss.State.ChosenStore.ScrewStock = ss.State.BuyCard(ss.screwCard, ss.State.ChosenStore.ScrewPrice, ss.State.ChosenStore.ScrewStock)

	case ss.State.Controls.Key5:
		if tool != nil {
			switch tool.Quality {
			case state.OneStar:
				if ss.State.BuyCard(tool, 100, 1) == 0 {
					ss.toolBought++
				}
			case state.TwoStar:
				if ss.State.BuyCard(tool, 100, 1) == 0 {
					ss.toolBought++
				}
			case state.ThreeStar:
				if ss.State.BuyCard(tool, 200, 1) == 0 {
					ss.toolBought++
				}
			}
		}

	case ss.State.Controls.KeyTab, ss.State.Controls.KeySpace:
		ss.continueButton.OnClick()
	case ss.State.Controls.KeyEnter:
		ss.repairButton.OnClick()
	}
}
