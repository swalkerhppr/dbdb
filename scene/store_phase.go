package scene

import (
	"dbdb/components"
	"dbdb/state"
	"fmt"
	"image/color"
	"log"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
		BaseScene: NewBaseWithFill(width, height, color.Black),
	}
	sp.discardButton = components.NewButton("Discard", 500, 445, func() {
		sp.BaseScene.State.DiscardHand()
	})
	sp.playButton = components.NewButton("Play", 128, 445, sp.playSelected)
	sp.skipButton = components.NewButton("Skip (-1 Time)", 280, 445, func() {
		sp.BaseScene.State.TimeLeft--
		sp.nextEncounter()
	})
	sp.textBox = components.NewTextBox("", 19, 130, 32, 480, 96)
	return sp
}

func (s *storePhase) playSelected() {
	if s.BaseScene.State.PlaySelected() {
		s.nextEncounter()
	} 
	s.BaseScene.State.ClearSelectedCards()
}

func (s *storePhase) Draw(screen *ebiten.Image) {
	s.DrawScene(screen)
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButton1) {
		s.nextEncounter()
	}
	encounterCounter := fmt.Sprintf("Encounter %d/%d\n", s.BaseScene.State.EncounterNumber+1, len(s.BaseScene.State.ChosenStore.Encounters))
	switch s.BaseScene.State.CurrentEncounter {
	case state.EmployeeEncounter:
		// tool type or expertise type
		s.textBox.SetText(encounterCounter + "Employee: How can I help you?")
		s.encounterNPC.Draw(screen)
	case state.NeighborEncounter:
		// tool type or helper type
		s.textBox.SetText(encounterCounter + "Neighbor: Did you see the game last night?")
		s.encounterNPC.Draw(screen)
	case state.HoleEncounter:
		// plank or board type
		s.textBox.SetText(encounterCounter + "There is a big hole in the ground. Fix it?")
	case state.ShelfEncounter:
		// fastener or nail type
		s.textBox.SetText(encounterCounter + "A shelf is falling apart. Fix it?")

	}
	s.textBox.Draw(screen)
	s.hand.Draw(screen)
	s.discardButton.Draw(screen)
	s.playButton.Draw(screen)
	s.skipButton.Draw(screen)
}

func (s *storePhase) nextEncounter() {
	s.BaseScene.State.EncounterNumber++
	if s.BaseScene.State.EncounterNumber == len(s.BaseScene.State.ChosenStore.Encounters) {
		// TODO Change to results
		s.SceneManager.SwitchTo(SceneMap[BuildingPhase])
	} else {
		s.BaseScene.State.CurrentEncounter = s.BaseScene.State.ChosenStore.Encounters[s.BaseScene.State.EncounterNumber]
		s.encounterNPC = components.NewRandomFigure(100, 100, 1.2)
		if s.BaseScene.State.CurrentEncounter == state.EmployeeEncounter {
			s.encounterNPC = components.NewRandomEmployeeFigure(90, 120, 1.2)
		} else {
			s.encounterNPC = components.NewRandomNeighborFigure(90, 120, 1.2)
		}
	}
}

func (p *storePhase) Load(s *State, controller stagehand.SceneController[*State]) {
	// Generate random encounters and available items
	numEncounters := 0
	toolTimeCost := 0
	switch s.ChosenStore.StoreQuality {
	case state.OneStar:
		numEncounters = 2 + rand.Intn(4)
		toolTimeCost = 2
	case state.TwoStar:
		numEncounters = 4 + rand.Intn(4)
		toolTimeCost = 1
	case state.ThreeStar:
		numEncounters = 6 + rand.Intn(4)
		toolTimeCost = 1
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
	s.ChosenStore.AvailableTools = make([]*state.CardState, 2)
	s.ChosenStore.AvailableTools[0] = &state.CardState{
		CardID:   state.RandomToolID(),
		TimeCost:	toolTimeCost,
		UsesLeft: 4,
	}
	s.ChosenStore.AvailableTools[1] = &state.CardState{
		CardID:   state.RandomToolID(),
		TimeCost:	toolTimeCost,
		UsesLeft: 4,
	}
	s.Phase = state.StorePhase
	s.EncounterNumber = 0
	s.CurrentEncounter = s.ChosenStore.Encounters[0]
	log.Printf("Populated Store: %+v", s.ChosenStore)
	p.hand = components.NewCardHand(s)

	if s.CurrentEncounter == state.EmployeeEncounter {
		p.encounterNPC = components.NewRandomEmployeeFigure(90, 120, 1.2)
	} else {
		p.encounterNPC = components.NewRandomNeighborFigure(90, 120, 1.2)
	}

	s.ClearSelectedCards()
	p.BaseScene.Load(s, controller)
}

