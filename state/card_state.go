package state

type CombineResult uint32

const (
	Trash CombineResult = 0
	PlankNailHammer  = CombineResult(MaterialPlank | MaterialNail | ToolHammer)
	PlankNailNailGun = CombineResult(MaterialPlank | MaterialNail | ToolNailGun)
	BoardNailHammer  = CombineResult(MaterialBoard | MaterialNail | ToolHammer)
	BoardNailNailGun = CombineResult(MaterialBoard | MaterialNail | ToolNailGun)
	PlankScrewDrill  = CombineResult(MaterialPlank | MaterialScrew | ToolDrill)
	BoardScrewDrill  = CombineResult(MaterialBoard | MaterialScrew | ToolDrill)
	BoardSaw         = CombineResult(MaterialBoard | ToolSaw)
	BoardCircularSaw = CombineResult(MaterialBoard | ToolCircularSaw)
	PlankGlue        = CombineResult(MaterialPlank | ToolGlue)
)

type EncounterRequirement uint32

const (
	NotMet EncounterRequirement = 0
	EmployeeEncounter = EncounterRequirement(ToolType | ExpertiseType)
	NeighborEncounter = EncounterRequirement(ToolType | HelperType)
	HoleEncounter     = EncounterRequirement(MaterialBoard | MaterialPlank)
	ShelfEncounter    = EncounterRequirement(MaterialNail | MaterialScrew)
)

type MaterialOrToolQuality byte

const (
	OneStar MaterialOrToolQuality = iota
	TwoStar
	ThreeStar
)

type CardState struct {
	CardID  CardID
	MoneyCost float32
	UsesLeft  int
	Quality MaterialOrToolQuality
	FavoriteTool CardID
	Selected bool
}

func (c *CardState) IsPlayableAlone() bool {
	return c.CardID.IsHelper() || c.CardID.IsExpertise()
}

func (c *CardState) Combine(others ...*CardState) CombineResult {
	result := CombineResult(c.CardID)
	for _, other := range others {
		result |= CombineResult(other.CardID)
	}
	if !result.IsValid() {
		return Trash
	}
	return result
}

func (c *CardState) String() string {
	return c.CardID.AssetName()
}

func (c *CardState) Meets(e EncounterRequirement) bool {
	return e & EncounterRequirement(c.CardID) != 0
}

func (e EncounterRequirement) String() string {
	switch e {
	case EmployeeEncounter:
		return "Employee"
	case NeighborEncounter:
		return "Neighbor"
	case HoleEncounter:
		return "HoleTrap"
	case ShelfEncounter:
		return "ShelfTrap"
	}
	return "NotMet"
}

func (e EncounterRequirement) IsValid() bool {
	return e == EmployeeEncounter ||
		e == NeighborEncounter ||
		e == HoleEncounter     ||
		e == ShelfEncounter
}

func (c CombineResult) IsValid() bool {
	return c == PlankNailHammer  ||
		c == PlankNailNailGun ||
		c == BoardNailHammer  ||
		c == BoardNailNailGun ||
		c == PlankScrewDrill  ||
		c == BoardScrewDrill  ||
		c == BoardSaw         ||
		c == BoardCircularSaw ||
		c == PlankGlue
}
