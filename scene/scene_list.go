package scene

import "github.com/joelschutz/stagehand"

type SceneID int

const (
	MainMenu SceneID = iota
	ChooseStore
	StorePhase
	StoreShop
	StoreResults
	BuildingPhase
	DayResults
	GameResults
)

var SceneMap map[SceneID]stagehand.Scene[*State]

func InitializeScenes(width, height int) {
	SceneMap = map[SceneID]stagehand.Scene[*State]{
		MainMenu       : CreateMainMenu(width, height),
		ChooseStore    : CreateChooseStore(width, height),
		StorePhase     : CreateStorePhase(width, height),
		StoreShop      : CreateStoreShop(width, height),
		StoreResults   : CreateStoreResults(width, height),
		BuildingPhase  : CreateBuildingPhase(width, height),
		DayResults     : CreateDayResults(width, height),
		GameResults    : CreateGameResults(width, height),
	}
}
