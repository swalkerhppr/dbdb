package scene

import (
	"dbdb/components"
	"dbdb/state"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type chooseStore struct {
	choice1 *components.PopoverArea
	choice2 *components.PopoverArea
	choice3 *components.PopoverArea
	textBox *components.TextBox
	*BaseScene
}

func CreateChooseStore(width, height int) stagehand.Scene[*State] {
	cs := &chooseStore{
		BaseScene: NewBaseWithBG(width, height, "store_choice.png"),
		textBox: components.NewTextBox("Choose a store to collect supplies", 17, 178, 94, 320, 48),
	}

	textOne := components.NewTextBox("Grade: C\nPrice: $\nTime: Low", 17, 112, 250, 96, 64)
	cs.choice1 = components.NewPopoverArea(128, 180, 64, 64, textOne, func() {
		cs.BaseScene.State.ChosenStore = state.StoreData{
			StoreQuality: state.OneStar,
			BoardPrice:   18.0,
			PlankPrice:   53.0,
			NailPrice:    8.0,
			ScrewPrice:   5.0,
		}
		cs.SceneManager.SwitchWithTransition(SceneMap[StorePhase], stagehand.NewDurationTimedSlideTransition[*State](stagehand.BottomToTop, time.Millisecond * 250))
	})

	textTwo := components.NewTextBox("Grade: B\nPrice: $$\nTime: Mid", 17, 272, 250, 96, 64)
	cs.choice2 = components.NewPopoverArea(288, 180, 64, 64, textTwo, func() {
		cs.BaseScene.State.ChosenStore = state.StoreData{
			StoreQuality: state.TwoStar,
			BoardPrice:   40.0,
			PlankPrice:   64.0,
			NailPrice:    10.0,
			ScrewPrice:   9.0,
		}
		cs.SceneManager.SwitchWithTransition(SceneMap[StorePhase], stagehand.NewDurationTimedSlideTransition[*State](stagehand.BottomToTop, time.Millisecond * 250))
	})

	textThree := components.NewTextBox("Grade: A\nPrice: $$$\nTime: High", 17, 432, 250, 96, 64)
	cs.choice3 = components.NewPopoverArea(448, 180, 64, 64, textThree, func() {
		cs.BaseScene.State.ChosenStore = state.StoreData{
			StoreQuality: state.ThreeStar,
			BoardPrice:   50.0,
			PlankPrice:   80.0,
			NailPrice:    12.0,
			ScrewPrice:   10.0,
		}
		cs.SceneManager.SwitchWithTransition(SceneMap[StorePhase], stagehand.NewDurationTimedSlideTransition[*State](stagehand.BottomToTop, time.Millisecond * 250))
	})


	return cs
}

func (s *chooseStore) Draw(screen *ebiten.Image) {
	s.BaseScene.DrawScene(screen)
	s.choice1.Draw(screen)
	s.choice2.Draw(screen)
	s.choice3.Draw(screen)
	s.textBox.Draw(screen)
}
