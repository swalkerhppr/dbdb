package state

type StoreData struct {
	BoardPrice float32
	PlankPrice float32
	NailPrice  float32
	ScrewPrice float32
	BoardStock int
	PlankStock int
	NailStock  int
	ScrewStock int
	StoreQuality MaterialOrToolQuality
	AvailableTools []*CardState
	Encounters []EncounterRequirement
}

type GlobalState struct {
	Phase     GamePhase
	MoneyLeft float32
	Day       int
	MaxDays   int
	TimeLeft  int
	GameWon   bool
	
	ChosenStore StoreData

	ShowAlert bool
	AlertText string

	ActiveHelpers []*HelperState

	EncounterNumber  int
	CurrentEncounter EncounterRequirement
	EncounterHelperCardID CardID
	SelectedCardsPlayable bool

	RequiredPlankParts int
	RequiredBoardParts int
	PlankPartsBuilt int
	BoardPartsBuilt int

	Deck []*CardState
	HeldCards []*CardState

	ActiveExpertise CardID
	Controls ControlHandler

	handStartIdx int
	topCardIdx int

	selectedCards []*CardState
}

func InitialState() *GlobalState {
	return &GlobalState{
		TimeLeft: 24,
		RequiredPlankParts: 5,
		RequiredBoardParts: 5,
		handStartIdx: 0,
		topCardIdx: 5,
		selectedCards : make([]*CardState, 0, 5),
	}
}

func InitialDeck(difficulty int) []*CardState {
	switch difficulty {
	case 0: //Easy
		return []*CardState{
			{ CardID : MaterialBoard, Quality: ThreeStar  },
			{ CardID : MaterialBoard, Quality: TwoStar  },
			{ CardID : MaterialPlank, Quality: ThreeStar  },
			{ CardID : MaterialNail,  Quality: ThreeStar },
			{ CardID : MaterialNail,  Quality: TwoStar },
			{ CardID : ToolHammer,    Quality: OneStar, UsesLeft: 4 },
			{ CardID : ToolSaw,       Quality: OneStar, UsesLeft: 4 },
			{ CardID : HelperGuru,    MoneyCost : 150, FavoriteTool: ToolHammer },
			{ CardID : ExpertiseTradesman },
		}
	case 1: // Normal
		return []*CardState{
			{ CardID : MaterialBoard, Quality: TwoStar  },
			{ CardID : MaterialPlank, Quality: TwoStar  },
			{ CardID : MaterialNail,  Quality: TwoStar },
			{ CardID : MaterialNail,  Quality: TwoStar },
			{ CardID : ToolHammer,    Quality: OneStar, UsesLeft: 4 },
			{ CardID : ToolSaw,       Quality: OneStar, UsesLeft: 4 },
			{ CardID : HelperGuru,    MoneyCost : 150, FavoriteTool: ToolHammer },
		}
	case 2: // Hard
		return []*CardState{
			{ CardID : MaterialBoard, Quality: TwoStar  },
			{ CardID : MaterialBoard, Quality: TwoStar  },
			{ CardID : MaterialPlank, Quality: TwoStar  },
			{ CardID : MaterialNail,  Quality: TwoStar },
			{ CardID : ToolHammer,    Quality: OneStar, UsesLeft: 4 },
			{ CardID : ToolSaw,       Quality: OneStar, UsesLeft: 4 },
		}
	}
	return nil
}
