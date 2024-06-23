package state

import (
	"log"
	"math/rand"
	"slices"
)

type GamePhase byte

const (
	NotStarted GamePhase = iota
	StorePhase
	BuildPhase
)

type StoreData struct {
	BoardPrice float32
	PlankPrice float32
	NailPrice  float32
	ScrewPrice float32
	StoreQuality MaterialQuality
	AvailableTools []*CardState
	Encounters []EncounterRequirement
}

type GlobalState struct {
	Phase     GamePhase
	MoneyLeft float32
	Day       int
	TimeLeft  int
	
	ChosenStore StoreData

	ActiveHelpers []*HelperState

	EncounterNumber  int
	CurrentEncounter EncounterRequirement
	SelectedCardsPlayable bool

	DeckRequirements []CombineResult
	DeckProgress int

	Deck []*CardState

	handStartIdx int
	topCardIdx int

	selectedCards []*CardState
}

func InitialState() *GlobalState {
	return &GlobalState{
		Deck : []*CardState{
			{ CardID : MaterialBoard, MaterialQuality: OneStar, TimeCost: 3 },
			{ CardID : MaterialPlank, MaterialQuality: TwoStar, TimeCost: 2 },
			{ CardID : MaterialPlank, MaterialQuality: TwoStar, TimeCost: 2 },
			{ CardID : MaterialPlank, MaterialQuality: TwoStar, TimeCost: 2 },
			{ CardID : MaterialNail, MaterialQuality: OneStar, TimeCost: 3},
			{ CardID : MaterialNail, MaterialQuality: TwoStar, TimeCost: 2},
			{ CardID : MaterialNail, MaterialQuality: TwoStar, TimeCost: 2},
			{ CardID : ToolHammer,   UsesLeft: 4, TimeCost: 2 },
			{ CardID : ToolSaw,      UsesLeft: 4, TimeCost: 2 },
			{ CardID : HelperGuru,   MoneyCost : 300, FavoriteTool: ToolHammer },
		},
		selectedCards : make([]*CardState, 0, 5),
	}
}

// Plays the selected cards. Updates the deck/cards as needed. Returns true if the cards were played successfully
func (s *GlobalState) PlaySelected() bool {
	log.Printf("Playing %+v, %v", s.selectedCards, s.SelectedCardsPlayable)
	return s.SelectedCardsPlayable
}

func (s *GlobalState) ToggleSelectCard(handIdx int) {
	card := s.Deck[s.handStartIdx + handIdx]
	card.Selected = !card.Selected
	if card.Selected {
		s.selectedCards = append(s.selectedCards, card)
	} else {
		selIdx := -1
		for i := range s.selectedCards {
			if s.selectedCards[i] == card {
				selIdx = i
				break
			}
		}
		if selIdx >= 0 {
			s.selectedCards = slices.Delete(s.selectedCards, selIdx, selIdx + 1)
		}
	}

	if s.Phase == StorePhase {
		s.SelectedCardsPlayable = len(s.selectedCards) == 1 && s.selectedCards[0].Meets(s.CurrentEncounter)
	} else if s.Phase == BuildPhase {
		s.SelectedCardsPlayable = (len(s.selectedCards) == 1 && s.selectedCards[0].IsPlayableAlone()) ||
			(len(s.selectedCards) >= 2 && s.selectedCards[0].Combine(s.selectedCards[1:len(s.selectedCards)]...) != Trash)
	}
	log.Printf("Toggled Selected: %+v", s.selectedCards)
}

func (s *GlobalState) ClearSelectedCards() {
	for i := range s.selectedCards {
		s.selectedCards[i].Selected = false
	}
	s.selectedCards = s.selectedCards[0:0]
}

func (s *GlobalState) GetHand() []*CardState {
	return s.Deck[s.handStartIdx:s.topCardIdx]
}

func (s *GlobalState) DiscardHand() {
	log.Printf("Discarding hand: %d-%d", s.handStartIdx, s.topCardIdx)
	s.handStartIdx += 5
	s.topCardIdx = s.handStartIdx + 5
	if s.topCardIdx >= len(s.Deck) {
		log.Printf("Shuffling: %d-%d", s.handStartIdx, s.topCardIdx)
		s.ShuffleCards(len(s.Deck) - s.handStartIdx, 4)
	}
	log.Printf("New hand: %d-%d", s.handStartIdx, s.topCardIdx)
	s.ClearSelectedCards()
}

func (s *GlobalState) DiscardCard(i int) {
	// Move the card to after hand in the deck
	card := s.Deck[i]
	s.Deck[i] = s.Deck[s.handStartIdx]
	s.Deck[s.handStartIdx] = card
	s.DrawCard()
}

// Draws a card from the deck, shuffles if necessary
func (s *GlobalState) DrawCard() {
	// Draw a card
	s.handStartIdx++
	s.topCardIdx++
	if s.topCardIdx > len(s.Deck) {
		s.ShuffleCards(4, 4)
	}
}

// Shuffles all except for rem cards for numShuffles, preserves current hand
func (s *GlobalState) ShuffleCards(rem, numShuffles int) {
	log.Printf("Pre shuffle: %v", s.Deck)
	// Move remainder to the top
	for dst := range rem {
		src := len(s.Deck) - rem + dst
		hold := s.Deck[src]
		s.Deck[src] = s.Deck[dst]
		s.Deck[dst] = hold
	}
	log.Printf("Post Reorg: %v", s.Deck)
	// Shuffle the rest of the cards
	for range numShuffles {
		leftIdx := rem
		rightIdx := len(s.Deck) - 1
		for i := range len(s.Deck) {
			hold := s.Deck[i]
			if rand.Int() % 2 == 0 {
				s.Deck[i] = s.Deck[leftIdx]
				s.Deck[leftIdx] = hold
				if leftIdx < len(s.Deck) - 1 {
					leftIdx++
				}
			} else {
				s.Deck[i] = s.Deck[rightIdx]
				s.Deck[rightIdx] = hold
				if rightIdx > 0 {
					rightIdx--
				}
			}
		}
	}
	s.handStartIdx = 0
	s.topCardIdx = 5
	log.Printf("Post shuffle: %v", s.Deck)
}

type HelperState struct {
	FavoriteTool CardID
	IsBrokenHolder bool
	HeldCards []*CardState
}
