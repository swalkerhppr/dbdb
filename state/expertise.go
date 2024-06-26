package state


func (s *GlobalState) ApplyExpertise(c CardID) {
	if !c.IsExpertise() {
		return
	}
	s.ActiveExpertise |= c
}

func (s *GlobalState) IsExpertiseActive(c CardID) bool {
	return c & s.ActiveExpertise !=0
}

func (s *GlobalState) DisableExpertise(c CardID) {
	s.ActiveExpertise ^= c
}
