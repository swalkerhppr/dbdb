package state

import (
	"log"
)

type HelperState struct {
	Card *CardState
	IsBrokenHolder bool
	MaterialCard *CardState
	ToolCard *CardState
}

func (s *GlobalState) DeployHelper(c *CardState) bool {
	if s.MoneyLeft < c.MoneyCost {
		s.alert("Not enough money!")
		return false
	}
	s.MoneyLeft -= c.MoneyCost
	s.ActiveHelpers = append(s.ActiveHelpers, &HelperState{
		Card: c,
		IsBrokenHolder: false, // TODO
	})
	s.HeldCards = append(s.HeldCards, []*CardState{
		{ CardID : EmptyCardSlot | CardID(MaterialType) },
		{ CardID: EmptyCardSlot | c.FavoriteTool },
	}...)
	return true
}

// Returns true if a single card is selected and it can fit in the given slot
func (s *GlobalState) CanSlot(c *CardState) bool {
	if !c.CardID.IsEmptyCardSlot() || len(s.selectedCards) != 1 {
		return false
	}
	if s.selectedCards[0].CardID.IsTool() {
		return c.CardID ^ s.selectedCards[0].CardID == EmptyCardSlot
	}
	return s.selectedCards[0].CardID.IsMaterial() && c.CardID ^ MaterialType == EmptyCardSlot
}

func (s *GlobalState) HoldCard(c *CardState) {
	log.Printf("Holding %+v", s.selectedCards[0])
	for i, held := range s.HeldCards {
		if held == c {
			s.HeldCards[i] = s.selectedCards[0]
			s.HeldCards[i].Selected = false
			s.DestroyCard(s.selectedCards[0])
			s.selectedCards = s.selectedCards[0:0]
			return
		}
	}
	log.Println("Could not find card to hold")
}

func (s *GlobalState) ReleaseCard(c *CardState) {
	log.Printf("Releasing %+v", c)
	for i, held := range s.HeldCards {
		if held == c {
			s.HeldCards[i].Selected = false
			s.AddCard(s.HeldCards[i])
			if held.CardID.IsMaterial() {
				s.HeldCards[i] = &CardState{ CardID: MaterialType | EmptyCardSlot }
			} else {
				s.HeldCards[i] = &CardState{ CardID: held.CardID | EmptyCardSlot }
			}
			return
		}
	}
	log.Println("Could not find card to release")
}
