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

// Selects a card if it is not selected, deselects if it is not
// keeps track of all selected cards internally
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
	shuffled := false
	log.Printf("Discarding hand: %d-%d", s.handStartIdx, s.topCardIdx)
	s.handStartIdx += 5
	s.topCardIdx = s.handStartIdx + 5
	if s.topCardIdx > len(s.Deck) && len(s.Deck) > 5 {
		log.Printf("Shuffling: %d-%d", s.handStartIdx, s.topCardIdx)
		if s.handStartIdx >= len(s.Deck) {
			s.handStartIdx = len(s.Deck)
			s.topCardIdx = len(s.Deck)
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
	s.Deck = append(s.Deck, card)
}

// Removes a card from the deck. does not remove the card if it doesn't exist
func (s *GlobalState) DestroyCard(card *CardState) {
	log.Printf("Deck pre-destroy: %+v", s.Deck)
	for i, c := range s.Deck {
		if c == card {
			s.Deck = slices.Delete(s.Deck, i, i+1)
			if len(s.Deck) > s.topCardIdx {
				s.ShuffleCards(0, 4)
			}
			log.Printf("Deck post-destroy: %+v", s.Deck)
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
	s.topCardIdx++
	assets.Registry.Sound("singlecard.ogg").Play()
	if s.topCardIdx > len(s.Deck) {
		s.ShuffleCards(4, 4)
		return true
	}
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
		if leftIdx >= rightIdx {
			break
		}
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
	assets.Registry.Sound("shuffle.ogg").Play()
	log.Printf("Post shuffle: %v", s.Deck)
}
