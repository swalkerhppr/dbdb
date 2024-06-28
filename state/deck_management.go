package state

import (
	"dbdb/assets"
	"fmt"
	"log"
	"math/rand"
	"slices"
)

func (s *GlobalState) CardsLeftInDeck() int {
	return len(s.Deck) - s.topCardIdx
}

func (s *GlobalState) SellCardWithID(cid CardID) {
	if s.TimeLeft == 0 {
		s.alert("Not enough time!")
		return
	}

	price := 5
	switch {
	case cid.IsExpertise():
		price = 100
	case cid.IsTool():
		price = 50
	case cid == MaterialPlank:
		price = 25
	case cid == MaterialBoard:
		price = 10
	}

	price *= 1 + int(s.ChosenStore.StoreQuality)

	for _, c := range s.Deck {
		if c.CardID & cid != 0 {
			log.Printf("Sold %+v", c)
			s.MoneyLeft += float32(price)
			s.DestroyCard(c)
			s.TimeLeft--
			return
		}
	}
}

func (s *GlobalState) BuyCard(card *CardState, cost float32, stock int) int {
	if s.TimeLeft == 0 {
		s.alert("Not enough time!")
		return stock
	}
	if len(s.Deck) >= 30 {
		s.alert("Can't buy anymore cards")
		return stock
	}
	if s.MoneyLeft < cost {
		s.alert("Not enough money!")
		return stock
	}
	if stock == 0 {
		s.alert("Out of stock!")
		return 0
	}
	s.MoneyLeft -= cost
	s.TimeLeft--
	s.AddCard(&CardState{
		CardID:       card.CardID,
		UsesLeft:     card.UsesLeft,
		Quality:      card.Quality,
		Selected:     false,
	})
	s.alert(fmt.Sprintf("Bought %v for $%.2f", card, cost))

	return stock-1
}

func (s *GlobalState) SetHint() {
	s.ShowHint = false
	for _, c := range s.Deck {
		c.Hint = false
	}
	if s.Phase == StorePhase {
		for _, c := range s.GetHand() {
			if ( s.CurrentEncounter == EmployeeEncounter && (c.CardID.IsTool() || c.CardID.IsExpertise()) ) || 
			   ( s.CurrentEncounter == NeighborEncounter && (c.CardID.IsTool() || c.CardID.IsHelper()) )    || 
			   ( s.CurrentEncounter == HoleEncounter && (c.CardID == MaterialBoard || c.CardID == MaterialPlank ) ) || 
			   ( s.CurrentEncounter == ShelfEncounter && (c.CardID == MaterialNail || c.CardID == MaterialScrew ) ) {
						c.Hint = true
						break
			}
		}
	}
	if s.Phase == BuildPhase {
		var result CombineResult = 0
		var find CombineResult
		hand := s.GetHand()
		for _, c := range hand {
			result |= CombineResult(c.CardID)
		}
		for _, match := range []CombineResult{ PlankNailHammer, PlankNailNailGun, PlankScrewDrill, BoardNailHammer, BoardNailNailGun, BoardScrewDrill, BoardSaw, BoardCircularSaw, PlankGlue } {
			if match & result == match {
				find = match
				break
			}
		}

		if find != 0 {
			s.NoPlay = false
			for _, c := range hand {
				if CombineResult(c.CardID) & find != 0 {
					c.Hint = true
					find ^= CombineResult(c.CardID)
				}
				if find == 0 {
					break
				}
			}
		} else {
			s.NoPlay = true
		}
	}
}

// Selects a card if it is not selected, deselects if it is not
// keeps track of all selected cards internally
func (s *GlobalState) ToggleSelectCard(handIdx int) {
	card := s.Deck[s.handStartIdx + handIdx]
	card.Selected = !card.Selected
	card.Hint = false
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
		s.UpdateSelectedCardsPlayable()
	}
	log.Printf("Toggled Selected: %+v", s.selectedCards)
}

func (s *GlobalState) UpdateSelectedCardsPlayable()  {
	allSelected := s.selectedCards
	if len(s.selectedCards) > 0 {
		for _, h := range s.HeldCards {
			if h.Selected {
				allSelected = append(allSelected, h)
			}
		}
	}
	s.SelectedCardsPlayable = (len(s.selectedCards) == 1 && s.selectedCards[0].IsPlayableAlone()) ||
		(len(allSelected) >= 2 && allSelected[0].Combine(allSelected[1:]...) != Trash)
}

// Clears the list of selected cards. Updates cards internal state to not be seleted
func (s *GlobalState) ClearSelectedCards() {
	for _, c := range s.GetHand() {
		c.Hint = false
	}
	for i := range s.selectedCards {
		s.selectedCards[i].Selected = false
	}
	for i := range s.HeldCards {
		s.HeldCards[i].Selected = false
	}
	s.selectedCards = s.selectedCards[0:0]
	s.SelectedCardsPlayable = false
}

// Gets the current hand
func (s *GlobalState) GetHand() []*CardState {
	if len(s.Deck) <= 5{
		s.topCardIdx = len(s.Deck)
	}
	return s.Deck[s.handStartIdx:s.topCardIdx]
}

// Discards the current hand and draws a new hand. Shuffles if necessary
// Returns true if we had to shuffle
func (s *GlobalState) DiscardHand() bool {
	log.Printf("Discarding hand: %d-%d", s.handStartIdx, s.topCardIdx)

	shuffled := false
	s.handStartIdx += 5
	s.topCardIdx = s.handStartIdx + 5

	if s.topCardIdx > len(s.Deck) && len(s.Deck) > 5 {
		log.Printf("Shuffling: %d-%d", s.handStartIdx, s.topCardIdx)
		if s.handStartIdx >= len(s.Deck) {
			s.handStartIdx = len(s.Deck)
		}
		s.ShuffleCards(len(s.Deck) - s.handStartIdx, 4)
		shuffled = true
	}

	if len(s.Deck) <= 5 {
		s.handStartIdx = 0
		s.topCardIdx = len(s.Deck)
	}

	log.Printf("New hand: %d-%d", s.handStartIdx, s.topCardIdx)
	s.ClearSelectedCards()
	assets.Registry.Sound("dealcards.ogg").Play()
	return shuffled
}

// Adds a card to the deck
func (s *GlobalState) AddCard(card *CardState) {
	s.Deck = slices.Insert(s.Deck, 0, card)
	s.handStartIdx++
	s.topCardIdx++
}

// Removes a card from the deck. does not remove the card if it doesn't exist
func (s *GlobalState) DestroyCard(card *CardState) {
	log.Printf("Destroying: %+v from %+v", card, s.Deck)
	for i, c := range s.Deck {
		if c == card {
			s.topCardIdx--
			if s.topCardIdx <= s.handStartIdx {
				s.topCardIdx = s.handStartIdx
			}
			s.Deck = slices.Delete(s.Deck, i, i+1)
			return
		}
	}
}

// Discards a card from hand.
// returns true if we had to shuffle
func (s *GlobalState) DiscardCard(card *CardState) bool {
	// Move the card to after hand in the deck
	i := 0
	for j, c := range s.Deck {
		if c == card {
			i = j
			break
		}
	}

	s.Deck[i] = s.Deck[s.handStartIdx]
	s.Deck[s.handStartIdx] = card
	return s.DrawCard()
}

// Draws a card from the deck, shuffles if necessary
func (s *GlobalState) DrawCard() bool {
	// Draw a card
	s.handStartIdx++
	if s.handStartIdx == s.topCardIdx {
		return s.DiscardHand()
	} 
	assets.Registry.Sound("singlecard.ogg").Play()
	return false
}

// Shuffles all except for rem cards for numShuffles, preserves current hand
func (s *GlobalState) ShuffleCards(rem, numShuffles int) {
	if len(s.Deck) < 5{
		return
	}
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
		for i := range len(s.Deck) - rem {
			if leftIdx >= rightIdx {
				break
			}

			hold := s.Deck[i + rem]
			if rand.Int() % 2 == 0 {
				s.Deck[i + rem] = s.Deck[leftIdx]
				s.Deck[leftIdx] = hold
				if leftIdx < len(s.Deck) - 1 {
					leftIdx++
				}
			} else {
				s.Deck[i + rem] = s.Deck[rightIdx]
				s.Deck[rightIdx] = hold
				if rightIdx > 0 {
					rightIdx--
				}
			}
		}
	}
	s.handStartIdx = 0
	s.topCardIdx = 5
	assets.Registry.Sound("shuffle.ogg").Play()
	log.Printf("Post shuffle: %v", s.Deck)
}
