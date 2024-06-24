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

	handStartIdx int
	topCardIdx int

	selectedCards []*CardState
}

func InitialState() *GlobalState {
	return &GlobalState{
		Deck : []*CardState{
			{ CardID : MaterialBoard, Quality: TwoStar  },
			{ CardID : MaterialBoard, Quality: TwoStar  },
			{ CardID : MaterialPlank, Quality: TwoStar  },
			{ CardID : MaterialPlank, Quality: TwoStar  },
			{ CardID : MaterialNail,  Quality: TwoStar },
			{ CardID : MaterialNail,  Quality: TwoStar },
			{ CardID : MaterialNail,  Quality: TwoStar },
			{ CardID : ToolHammer,    Quality: TwoStar, UsesLeft: 4 },
			{ CardID : ToolSaw,       Quality: TwoStar, UsesLeft: 4 },
			{ CardID : HelperGuru,    MoneyCost : 150, FavoriteTool: ToolHammer },
		},
		TimeLeft: 24,
		RequiredPlankParts: 5,
		RequiredBoardParts: 5,
		MoneyLeft: 1000,
		selectedCards : make([]*CardState, 0, 5),
	}
}
