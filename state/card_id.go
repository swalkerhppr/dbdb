package state

import "math/rand"

type CardType uint32

const (
	UnknownType CardType = 0x0000000
	MaterialType         = 0x000000F
	HelperType           = 0x00007F0
	ToolType             = 0x001F800
	ExpertiseType        = 0x1FE0000
)

type CardID uint32

const (
	MaterialPlank CardID = 1 << iota
	MaterialBoard
	MaterialScrew
	MaterialNail

	HelperColo
	HelperShip
	HelperHood
	HelperSass
	HelperOarm
	HelperGuru
	HelperCurt

	ToolHammer
	ToolSaw
	ToolDrill
	ToolGlue
	ToolCircularSaw
	ToolNailGun

	ExpertiseTradesman
	ExpertiseBlacksmith
	ExpertiseRoofer
	ExpertiseLumberjack
	ExpertiseGrifter
	ExpertiseDuplicator
	ExpertiseWoodsman
	ExpertiseOptimizer
)

var cardIDNames = map[CardID]string{
	MaterialPlank : "plank",
	MaterialBoard : "board",
	MaterialScrew : "screw",
	MaterialNail  : "nail",

	HelperColo    : "colo",
	HelperShip    : "shop",
	HelperHood    : "hood",
	HelperSass    : "sass",
	HelperOarm    : "oarm",
	HelperGuru    : "guru",
	HelperCurt    : "curt",

	ToolHammer      : "hammer",
	ToolSaw         : "saw",
	ToolDrill       : "drill",
	ToolGlue        : "glue",
	ToolCircularSaw : "circular-saw",
	ToolNailGun     : "nail-gun",

	ExpertiseTradesman  : "tradesman",
	ExpertiseBlacksmith : "blacksmith",
	ExpertiseRoofer     : "roofer",
	ExpertiseLumberjack : "lumberjack",
	ExpertiseGrifter    : "grifter",
	ExpertiseDuplicator : "duplicator",
	ExpertiseWoodsman   : "woodsman",
	ExpertiseOptimizer  : "optimizer",
}

func RandomCardID() CardID {
	t := rand.Intn(25)
	i := 0
	for card := range cardIDNames {
		if i == t {
			return card
		}
		i++
	}
	return 0
}

func RandomToolID() CardID {
	return []CardID{
		ToolHammer,
		ToolSaw,
		ToolDrill,
		ToolGlue,
		ToolCircularSaw,
		ToolNailGun,
	}[rand.Intn(6)]
}

func (c CardID) AssetName() string {
	return cardIDNames[c] + "-card"
}

func (c CardID) IconName() string {
	if c.IsMaterial() || c.IsTool() {
		return cardIDNames[c] + "-icon"
	}
	if c.IsExpertise() {
		return "expertise-icon"
	}
	return "helper-icon"
}

func (c CardID) CardType() CardType {
	if c.IsMaterial() {
		return MaterialType
	}
	if c.IsHelper() {
		return HelperType
	}
	if c.IsTool() {
		return ToolType
	}
	if c.IsExpertise() {
		return ExpertiseType
	}
	return UnknownType
}

func (c CardID) IsMaterial() bool {
	return c & MaterialType > 0
}

func (c CardID) IsHelper() bool {
	return c & HelperType > 0
}

func (c CardID) IsTool() bool {
	return c & ToolType > 0
}

func (c CardID) IsExpertise() bool {
	return c & ExpertiseType > 0
}
