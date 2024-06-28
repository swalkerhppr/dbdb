package scene

import (
	"dbdb/components"
	"dbdb/state"
	"math/rand"
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
		textBox: components.NewTextBox("Choose a store to buy supplies", 17, 178, 94, 320, 48),
	}
	return cs
}

func (s *chooseStore) Draw(screen *ebiten.Image) {
	s.BaseScene.DrawScene(screen)
	s.choice1.Draw(screen)
	s.choice2.Draw(screen)
	s.choice3.Draw(screen)
	s.textBox.Draw(screen)

	switch {
	case s.State.Controls.Key1:
		s.choice1.OnClick()
	case s.State.Controls.Key2:
		s.choice2.OnClick()
	case s.State.Controls.Key3:
		s.choice3.OnClick()
	}
}

func (cs *chooseStore) Load(s *State, controller stagehand.SceneController[*State]) {
	textOne := components.NewTextBox("Grade: C\nPrice: $\nTime: Low", 17, 112, 250, 96, 64)
	cs.choice1 = components.NewPopoverArea(128, 180, 64, 64, textOne, &s.Controls, func() {
		cs.BaseScene.State.ChosenStore = state.StoreData{
			StoreQuality: state.OneStar,
			BoardPrice:   50.0,
			PlankPrice:   75.0,
			NailPrice:    12.0,
			ScrewPrice:   10.0,
			PlankStock:   1 + rand.Intn(1),
			BoardStock:   2 + rand.Intn(2),
			NailStock:    2 + rand.Intn(10),
			ScrewStock:   2 + rand.Intn(10),
		}
		cs.SceneManager.SwitchWithTransition(SceneMap[StorePhase], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})

	textTwo := components.NewTextBox("Grade: B\nPrice: $$\nTime: Mid", 17, 272, 250, 96, 64)
	cs.choice2 = components.NewPopoverArea(288, 180, 64, 64, textTwo, &s.Controls, func() {
		cs.BaseScene.State.ChosenStore = state.StoreData{
			StoreQuality: state.TwoStar,
			BoardPrice:   75.0,
			PlankPrice:   100.0,
			NailPrice:    15.0,
			ScrewPrice:   12.0,
			PlankStock:   1 + rand.Intn(1),
			BoardStock:   2 + rand.Intn(2),
			NailStock:    2 + rand.Intn(12),
			ScrewStock:   2 + rand.Intn(12),
		}
		cs.SceneManager.SwitchWithTransition(SceneMap[StorePhase], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})

	textThree := components.NewTextBox("Grade: A\nPrice: $$$\nTime: High", 17, 432, 250, 96, 64)
	cs.choice3 = components.NewPopoverArea(448, 180, 64, 64, textThree, &s.Controls, func() {
		cs.BaseScene.State.ChosenStore = state.StoreData{
			StoreQuality: state.ThreeStar,
			BoardPrice:   100.0,
			PlankPrice:   120.0,
			NailPrice:    20.0,
			ScrewPrice:   15.0,
			PlankStock:   1 + rand.Intn(3),
			BoardStock:   2 + rand.Intn(3),
			NailStock:    2 + rand.Intn(15),
			ScrewStock:   2 + rand.Intn(15),
		}
		cs.SceneManager.SwitchWithTransition(SceneMap[StorePhase], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	s.TimeLeft = 24
	cs.BaseScene.Load(s, controller)
}
