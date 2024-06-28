package scene

import (
	"dbdb/components"
	"dbdb/state"
	"fmt"
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
		BaseScene: NewBaseWithBG(width, height, "store_phase.png"),
	}
	sp.textBox = components.NewTextBox("", 19, 130, 32, 480, 96)
	return sp
}

func (s *storePhase) Draw(screen *ebiten.Image) {
	if s.State.CardsLeftInDeck() < 5 {
		s.discardButton.SetText("New Hand (-2)")
	} else {
		s.discardButton.SetText("New Hand")
	}
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
		s.textBox.SetText(encounterCounter + "There is a hole in the ground. Cover it?")
		s.textBox.Draw(screen)
		components.NewIndicator("plank-icon-grey", " OR", 328, 96).Draw(screen)
		components.NewIndicator("board-icon-grey", "", 385, 96).Draw(screen)
	case state.ShelfEncounter:
		// fastener or nail type
		s.textBox.SetText(encounterCounter + "A shelf is missing a part. Fix it?")
		s.textBox.Draw(screen)
		components.NewIndicator("nail-icon-grey", " OR", 328, 96).Draw(screen)
		components.NewIndicator("screw-icon-grey", "", 385, 96).Draw(screen)

	}
	s.hand.Draw(screen)
	s.discardButton.Draw(screen)
	s.playButton.Draw(screen)
	s.skipButton.Draw(screen)

	s.DrawIndicators(screen)

	if s.State.Controls.KeyEnter {
		s.skipButton.OnClick()
	}
	if s.State.Controls.KeyTab {
		s.discardButton.OnClick()
	}
	if s.State.Controls.KeySpace {
		s.playButton.OnClick()
	}
}

func (s *storePhase) nextEncounter() {
	if s.State.EncounterNumber == len(s.State.ChosenStore.Encounters) - 1 {
		s.SceneManager.SwitchWithTransition(SceneMap[StoreShop], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))

	} else {
		s.State.EncounterNumber++
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
	p.playButton = components.NewButton("Play", 96, 445, &s.Controls, func() {
		if s.PlaySelected() {
			p.nextEncounter()
		} 
	})
	p.discardButton = components.NewButton("New Hand", 235, 445, &s.Controls, func() {
		if s.DiscardHand() && s.TimeLeft > 0 {
			s.TimeLeft -= 2
		}
	})

	p.skipButton = components.NewButton("Skip", 532, 445, &s.Controls, func() {
		s.TimeLeft -= 1 + int(s.ChosenStore.StoreQuality)
		p.nextEncounter()
	})

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
	s.ShuffleCards(0, 4)
	p.BaseScene.Load(s, controller)
	p.skipButton.SetText(fmt.Sprintf("Skip (-%d)", 1 + s.ChosenStore.StoreQuality))
}

