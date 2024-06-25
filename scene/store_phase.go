package scene

import (
	"dbdb/components"
	"dbdb/state"
	"fmt"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type storePhase struct {
	hand *components.CardHand
	discardButton *components.Button
	playButton *components.Button
	skipButton *components.Button
	textBox *components.TextBox
	encounterNPC *components.Figure
	*BaseScene
}

func CreateStorePhase(width, height int) stagehand.Scene[*State] {
	sp := &storePhase{
		BaseScene: NewBaseWithFill(width, height, color.White),
	}
	sp.playButton = components.NewButton("Play", 96, 445, func() {
		if sp.State.PlaySelected() {
			sp.nextEncounter()
		} 
	})
	sp.discardButton = components.NewButton("New Hand", 235, 445, func() {
		if sp.State.DiscardHand() && sp.State.TimeLeft > 0 {
			sp.State.TimeLeft--
		}
	})

	sp.skipButton = components.NewButton("Skip (-2 Time)", 532, 445, func() {
		sp.State.TimeLeft -= 2
		sp.nextEncounter()
	})
	sp.textBox = components.NewTextBox("", 19, 130, 32, 480, 96)
	return sp
}

func (s *storePhase) Draw(screen *ebiten.Image) {
	s.DrawScene(screen)
	encounterCounter := fmt.Sprintf("Encounter %d/%d\n", s.State.EncounterNumber+1, len(s.State.ChosenStore.Encounters))
	switch s.State.CurrentEncounter {
	case state.EmployeeEncounter:
		// tool type or expertise type
		s.textBox.SetText(encounterCounter + "Employee: How can I help you?")
		s.textBox.Draw(screen)
		s.encounterNPC.Draw(screen)
		components.NewIndicator("tool-icon-grey", " OR", 328, 96).Draw(screen)
		components.NewIndicator("expertise-icon-grey", "", 379, 98).Draw(screen)
	case state.NeighborEncounter:
		// tool type or helper type
		s.textBox.SetText(encounterCounter + "Neighbor: Did you see the game last night?")
		s.textBox.Draw(screen)
		s.encounterNPC.Draw(screen)
		components.NewIndicator("tool-icon-grey", " OR", 328, 96).Draw(screen)
		components.NewIndicator("helper-icon-grey", "", 379, 98).Draw(screen)
	case state.HoleEncounter:
		// plank or board type
		s.textBox.SetText(encounterCounter + "You see a hole in the ground. Cover it?")
		s.textBox.Draw(screen)
		components.NewIndicator("plank-icon-grey", "  OR", 328, 96).Draw(screen)
		components.NewIndicator("board-icon-grey", "", 385, 96).Draw(screen)
	case state.ShelfEncounter:
		// fastener or nail type
		s.textBox.SetText(encounterCounter + "A shelf is falling apart. Fix it with a nail or screw?")
		s.textBox.Draw(screen)
		components.NewIndicator("nail-icon-grey", "  OR", 328, 96).Draw(screen)
		components.NewIndicator("screw-icon-grey", "", 385, 96).Draw(screen)

	}
	s.hand.Draw(screen)
	s.discardButton.Draw(screen)
	s.playButton.Draw(screen)
	s.skipButton.Draw(screen)

	s.DrawIndicators(screen)
}

func (s *storePhase) nextEncounter() {
	s.State.EncounterNumber++
	if s.State.EncounterNumber == len(s.State.ChosenStore.Encounters) {
		s.SceneManager.SwitchWithTransition(SceneMap[StoreShop], stagehand.NewDurationTimedSlideTransition[*State](stagehand.TopToBottom, time.Millisecond * 500))

	} else {
		s.State.CurrentEncounter = s.State.ChosenStore.Encounters[s.State.EncounterNumber]
		s.encounterNPC = components.NewRandomFigure(100, 100, 1.2)
		if s.State.CurrentEncounter == state.EmployeeEncounter {
			s.encounterNPC = components.NewRandomEmployeeFigure(90, 120, 1.2)
			s.State.EncounterHelperCardID = s.encounterNPC.CardID()

		} else {
			s.encounterNPC = components.NewRandomNeighborFigure(90, 120, 1.2)
			s.State.EncounterHelperCardID = s.encounterNPC.CardID()
		}
	}
}

func (p *storePhase) Load(s *State, controller stagehand.SceneController[*State]) {
	// Generate random encounters
	numEncounters := 0
	switch s.ChosenStore.StoreQuality {
	case state.OneStar:
		numEncounters = 2 + rand.Intn(4)
	case state.TwoStar:
		numEncounters = 4 + rand.Intn(4)
	case state.ThreeStar:
		numEncounters = 6 + rand.Intn(4)
	}
	s.ChosenStore.Encounters = make([]state.EncounterRequirement, numEncounters)
	for i := range numEncounters {
		 n := rand.Intn(100)
		 if n % 23 == 0 {
		 	 // 4 in 100
			 s.ChosenStore.Encounters[i] = state.HoleEncounter
		 } else if n % 13 == 0 {
			 // 7 in 100
			 s.ChosenStore.Encounters[i] = state.ShelfEncounter
		 } else if n % 3 == 0 {
			 // 33 in 100
			 s.ChosenStore.Encounters[i] = state.NeighborEncounter
		 } else {
			 s.ChosenStore.Encounters[i] = state.EmployeeEncounter
		 }
	}
	s.Phase = state.StorePhase
	s.EncounterNumber = 0
	s.CurrentEncounter = s.ChosenStore.Encounters[0]
	log.Printf("Populated Store: %+v", s.ChosenStore)
	p.hand = components.NewCardHand(s)

	if s.CurrentEncounter == state.EmployeeEncounter {
		p.encounterNPC = components.NewRandomEmployeeFigure(90, 120, 1.2)
		s.EncounterHelperCardID = p.encounterNPC.CardID()
	} else {
		p.encounterNPC = components.NewRandomNeighborFigure(90, 120, 1.2)
		s.EncounterHelperCardID = p.encounterNPC.CardID()
	}

	s.ClearSelectedCards()
	p.BaseScene.Load(s, controller)
}

