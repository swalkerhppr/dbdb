package state

import (
	"log"
	"math/rand"
	"strings"
)

// Card types
const (
	UnknownType       = 0x0000000
	MaterialType      = 0x000000F
	HelperType        = 0x00007F0
	ToolType          = 0x001F800
	ExpertiseType     = 0x1FE0000
	EmptyCardSlotType = 0x2000000
)

type CardID uint32

const (
	MaterialPlank CardID = 1 << iota
	MaterialBoard
	MaterialScrew
	MaterialNail

	HelperColo
	HelperShop
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

	EmptyCardSlot // Used to signify an empty hold location
	UnknownCard
)

var cardIDNames = map[CardID]string{
	MaterialPlank : "plank",
	MaterialBoard : "board",
	MaterialScrew : "screw",
	MaterialNail  : "nail",

	HelperColo    : "colo",
	HelperShop    : "shop",
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

func RandomMaterialID() CardID {
	return []CardID{
		MaterialBoard,
		MaterialPlank,
		MaterialNail,
		MaterialScrew,
	}[rand.Intn(4)]
}

func RandomExpertiseID() CardID {
	return []CardID{
		ExpertiseTradesman,
		ExpertiseBlacksmith,
		ExpertiseRoofer,
		ExpertiseLumberjack,
		ExpertiseGrifter,
		ExpertiseDuplicator,
		ExpertiseWoodsman,
		ExpertiseOptimizer,
	}[rand.Intn(4)]
}

func RandomToolOrMaterialCardID() CardID {
	return []CardID{
		ToolHammer,
		ToolSaw,
		ToolDrill,
		ToolGlue,
		ToolCircularSaw,
		ToolNailGun,
		MaterialBoard,
		MaterialPlank,
		MaterialNail,
		MaterialScrew,
	}[rand.Intn(10)]
}

// num must be less than 8 or a panic occurs
// The returned cards are unique
func RandomExpertiseIDs(num int) []CardID {
	if num >= 8 {
		panic("Can't generate more than 8 random expertisk cards")
	}

	exps := []CardID{
		ExpertiseTradesman,
		ExpertiseBlacksmith,
		ExpertiseRoofer,
		ExpertiseLumberjack,
		ExpertiseGrifter,
		ExpertiseDuplicator,
		ExpertiseWoodsman,
		ExpertiseOptimizer,
	}

	picks := make([]CardID, num)
	for i := range picks {
		if num * (i + 1) > 8 {
			picks[i] = exps[( i * num ) + rand.Intn(8 % num)]
		} else {
			picks[i] = exps[( i * num ) + rand.Intn(num)]
		}
	}

	return picks
}

func (c CardID) AssetName() string {
	if _, ok := cardIDNames[c]; !ok {
		log.Printf("Could not determine card's asset name! %v", c)
		return "unknown-card"
	}
	return cardIDNames[c] + "-card"
}

func (c CardID) IconName() string {
	if c.IsEmptyCardSlot() {
		c = c ^ EmptyCardSlot
		if c == MaterialType {
			return "material-icon"
		}
	}
	if c.IsMaterial() || c.IsTool() {
		return cardIDNames[c] + "-icon"
	}
	if c.IsExpertise() {
		return "expertise-icon"
	}
	return "helper-icon"
}

func (c CardID) DisplayName() string {
	if c.IsHelper() {
		return "HELPER"
	}
	name := cardIDNames[c]
	name = strings.ToTitle(name)
	return strings.Replace(name, "-", " ", 1)
}

func (c CardID) IsEmptyCardSlot() bool {
	return c & EmptyCardSlot != 0
}

func (c CardID) CanSlot(o CardID) bool {
	return (c & EmptyCardSlot) != 0 && ( (c ^ EmptyCardSlot) & o) != 0
}

func (c CardID) GetCompatibleSlotCard() CardID {
	return c ^ EmptyCardSlot
}

func (c CardID) CardType() uint32 {
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
