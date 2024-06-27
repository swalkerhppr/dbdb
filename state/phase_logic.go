package state

import (
	"dbdb/assets"
	"fmt"
	"log"
	"math/rand"
)

type GamePhase byte

const (
	NotStarted GamePhase = iota
	StorePhase
	BuildPhase
)

// Plays the selected cards. Updates the deck/cards as needed. Returns true if the cards were played successfully
func (s *GlobalState) PlaySelected() bool {
	log.Printf("Playing %+v, %v", s.selectedCards, s.SelectedCardsPlayable)
	defer s.ClearSelectedCards()
	if ! s.SelectedCardsPlayable {
		s.alert("You can't play that right now!")
		return false
	}
	if s.TimeLeft == 0 {
		s.alert("You do not have enough time!")
		return false
	}

	switch s.Phase {
	case StorePhase:
		return s.playStoreCard()
	case BuildPhase:
		return s.playBuildCard()
	}

	log.Println("Played cards out of playable phase!")
	return false
}

func (s *GlobalState) playStoreCard() bool {
	card := s.selectedCards[0]

	switch s.CurrentEncounter {
	case EmployeeEncounter:
		// Need to show a tool or expertise card
		isTool := card.CardID.IsTool()
		if isTool {
			if card.UsesLeft > 0 {
				card.UsesLeft--
			} else {
				s.alert("That tool is broken!")
				return false
			}
		}
		rollPrize := rand.Intn(100)
		if rollPrize >= 90 {
			// Get material or tool random prize
			prize := &CardState{
				CardID:   RandomToolOrMaterialCardID(),
				UsesLeft: 2 + int(s.ChosenStore.StoreQuality),
				Quality:  s.ChosenStore.StoreQuality,
			}
			s.AddCard(prize)
			s.alert(fmt.Sprintf("I think you dropped this.\n(Received %s)", prize.CardID.DisplayName()))
		} else if isTool && rollPrize >= 75 {
			// Get another tool of the same type
			prize := &CardState {
				CardID:   card.CardID,
				UsesLeft: 2 + int(s.ChosenStore.StoreQuality),
				Quality:  card.CardID.ToolQuality(),
			}
			s.AddCard(prize)
			s.alert(fmt.Sprintf("Maybe you'll like this brand\n(Received %s)", prize.CardID.DisplayName()))
		}
		s.DiscardCard(card)

	case NeighborEncounter:
		// Need to show a tool or helper card
		isTool := card.CardID.IsTool()
		if isTool {
			if card.UsesLeft > 0 {
				card.UsesLeft--
			} else {
				s.alert("That tool is broken!")
				return false
			}
		}
		rollPrize := rand.Intn(100)
		if rollPrize >= 80 {
			// Get a helper card
			favTool := card.CardID
			if !isTool {
				favTool = RandomToolID()
			}
			prize := &CardState{
				CardID:   s.EncounterHelperCardID,
				MoneyCost: float32(50)+ float32(100 * (2 - int(s.ChosenStore.StoreQuality))),
				FavoriteTool: favTool,
			}
			s.AddCard(prize)
			s.alert("Looks like you need some help.\n(Received HELPER)")
		} else if isTool && rollPrize >= 60 {
			// Get another tool of the same type
			prize := &CardState {
				CardID:   card.CardID,
				UsesLeft: rand.Intn(4),
				Quality:  card.CardID.ToolQuality(),
			}
			s.AddCard(prize)
			s.alert(fmt.Sprintf("You left this at my house.\n(Received %s)", prize.CardID.DisplayName()))
		}
		s.DiscardCard(card)

	case HoleEncounter:
		rollPrize := rand.Intn(100)
		if rollPrize > 35 { 
			prize := &CardState{
				CardID: RandomExpertiseID(),
			}
			s.AddCard(prize)
			s.alert(fmt.Sprintf("You became an expert!\n(Received %s)", prize.CardID.DisplayName()))
		}
		s.DestroyCard(card)

	case ShelfEncounter:
		rollPrize := rand.Intn(100)
		if rollPrize > 55 { 
			prize := &CardState{
				CardID: RandomExpertiseID(),
			}
			s.AddCard(prize)
			s.alert(fmt.Sprintf("You became an expert!\n(Received %s)", prize.CardID.DisplayName()))
		}
		s.DestroyCard(card)

	}
	return true
}

func (s *GlobalState) playBuildCard() bool {
	allSelected := s.selectedCards
	for _, h := range s.HeldCards {
		if h.Selected {
			allSelected = append(allSelected, h)
		}
	}

	card1 := allSelected[0]
	if len(allSelected) == 1 {
		switch card1.CardID.CardType() {
		case HelperType:
			if s.DeployHelper(card1) {
				s.DestroyCard(card1)
			}
		case ExpertiseType:
			s.ApplyExpertise(card1.CardID)
			s.DiscardCard(card1)
		default:
			log.Println("Single card was unplayable")
			return false
		}
		return true
	}

	totalTime := 0
	for _, c := range allSelected {
		switch c.CardID.CardType() {
		case MaterialType:
			totalTime += 3 - int(c.Quality)
		case ToolType:
			totalTime += 2 - int(c.Quality)
			if c.UsesLeft == 0 {
				s.alert("That tool is broken!")
			}
		}
	}

	result := card1.Combine(allSelected[1:]...)

	if s.IsExpertiseActive(ExpertiseWoodsman) {
		// All combinations involve wood
		totalTime--
		s.DisableExpertise(ExpertiseWoodsman)
	}

	if s.IsExpertiseActive(ExpertiseTradesman) &&
		( (CardID(result) & MaterialNail != 0) ||
		  (CardID(result) & MaterialScrew != 0) ) {
		// saw/glue don't use nail or screw
		totalTime--
		s.DisableExpertise(ExpertiseTradesman)
	}

	if s.IsExpertiseActive(ExpertiseRoofer) &&
		( (CardID(result) & ToolHammer != 0) ||
		  (CardID(result) & ToolNailGun != 0) ){
		totalTime--
		s.DisableExpertise(ExpertiseRoofer)
	}

	if s.IsExpertiseActive(ExpertiseLumberjack) &&
		( (CardID(result) & ToolSaw != 0) ||
			(CardID(result) & ToolCircularSaw != 0) ) {
		totalTime--
		s.DisableExpertise(ExpertiseRoofer)
	}

	if totalTime > s.TimeLeft {
		s.alert("You don't have enough time!")
		return false
	}

	switch result {
	case PlankNailHammer, PlankNailNailGun, PlankScrewDrill:
		if s.PlankPartsBuilt >= s.RequiredPlankParts {
			s.alert("You don't need anymore PLANK parts!")
			return false
		}
		s.PlankPartsBuilt++
		s.alert(fmt.Sprintf("Deck PLANK part built!(%d/%d)", s.PlankPartsBuilt, s.RequiredPlankParts))
	case BoardNailHammer, BoardNailNailGun, BoardScrewDrill:
		if s.BoardPartsBuilt >= s.RequiredBoardParts {
			s.alert("You don't need anymore BOARD parts!")
			return false
		}
		if s.PlankPartsBuilt < s.BoardPartsBuilt {
			s.alert("You need to build a PLANK part first!")
			return false
		}
		s.BoardPartsBuilt++
		s.alert(fmt.Sprintf("Deck BOARD part built!(%d/%d)", s.BoardPartsBuilt, s.RequiredBoardParts))
	case BoardSaw, BoardCircularSaw:
		// Add two planks
		lowestQual := MaterialOrToolQuality(100)
		for _, c := range s.selectedCards {
			switch c.CardID.CardType() {
			case MaterialType, ToolType:
				if lowestQual > c.Quality {
					lowestQual = c.Quality
				}
			}
		}
		s.AddCard(&CardState{
			CardID:       MaterialPlank,
			Quality:      lowestQual,
		})
		s.AddCard(&CardState{
			CardID:       MaterialPlank,
			Quality:      lowestQual,
		})
		s.alert("Got two PLANKs")
	case PlankGlue:
		// Add one board
		s.AddCard(&CardState{
			CardID:       MaterialBoard,
			Quality:      TwoStar,
		})
		s.alert("Got a BOARD")
	}

	switch result {
	case PlankNailHammer, BoardNailHammer:
		assets.Registry.Sound("hammer.ogg").Play()
	case PlankNailNailGun, BoardNailNailGun:
		assets.Registry.Sound("nailgun.ogg").Play()
	case PlankScrewDrill, BoardScrewDrill:
		assets.Registry.Sound("drill.ogg").Play()
	case BoardSaw:
		assets.Registry.Sound("saw.ogg").Play()
	case BoardCircularSaw:
		assets.Registry.Sound("circular-saw.ogg").Play()
	case PlankGlue:
		assets.Registry.Sound("glue.ogg").Play()
	}

	// All used materials are destroyed and all tools take a use
	for _, c := range allSelected {
		switch c.CardID.CardType() {
		case MaterialType:
			s.ReleaseCard(c)
			s.DestroyCard(c)
		case ToolType:
			c.UsesLeft--
			if s.IsExpertiseActive(ExpertiseBlacksmith) {
				c.UsesLeft++
				s.DisableExpertise(ExpertiseBlacksmith)
			}
			s.ReleaseCard(c)
			s.DiscardCard(c)
		}
	}

	s.TimeLeft -= totalTime
	if s.IsExpertiseActive(ExpertiseOptimizer) {
		s.TimeLeft += totalTime
		s.DisableExpertise(ExpertiseOptimizer)
	}

	if s.RequiredBoardParts == s.BoardPartsBuilt && s.RequiredPlankParts == s.PlankPartsBuilt {
		s.GameWon = true
	}

	return true
}

func (s *GlobalState) alert(text string) {
	s.ShowAlert = true
	s.AlertText = text
}
