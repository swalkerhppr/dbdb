package components

import (
	"dbdb/state"
	"fmt"
	"image/color"
	"slices"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

type DeckIndicators struct{
	globalState *state.GlobalState
	left int
	top  int
	cols int
	sellable bool
}

func NewDeckIndicators(s *state.GlobalState, cols, left, top int) *DeckIndicators {
	return &DeckIndicators{
		globalState : s,
		left        : left,
		top         : top,
		cols        : cols,
		sellable   : false,
	}
}

func NewShopDeckIndicators(s *state.GlobalState, cols, left, top int) *DeckIndicators {
	return &DeckIndicators{
		globalState : s,
		left        : left,
		top         : top,
		cols        : cols,
		sellable    : true,
	}
}

func (d *DeckIndicators) Draw(screen *ebiten.Image) {
	rclicked := inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonRight)
	var textColor color.Color = color.Black

	iconCountMap := make(map[state.CardID]int)
	for _, c := range d.globalState.Deck {
		if c.CardID.IsHelper() {
			iconCountMap[state.HelperType]++
		} else if c.CardID.IsExpertise() {
			iconCountMap[state.ExpertiseType]++
		} else {
			iconCountMap[c.CardID]++
		}
	}

	keys := make([]state.CardID, len(iconCountMap))
	col := 0
	for k := range iconCountMap {
		keys[col] = k
		col++
	}

	slices.SortFunc(keys, func(x, y state.CardID) int { return int(x) - int(y) })

	col = 0
	row := 0
	for _, key := range keys {
		l := d.left + (40 * col)
		t := d.top + (40 * row)
		NewIndicatorWithColor(key.IconName(), fmt.Sprint(" ", iconCountMap[key]), l, t, textColor).Draw(screen)
		col++
		if col == d.cols {
			col = 0
			row++
		}
	}

	col = 0
	row = 0
	for _, key := range keys {
		l := d.left + (40 * col)
		t := d.top + (40 * row)
		if IsMouseover(l, t, 32, 32, false) && d.sellable {
			if rclicked {
				d.globalState.SellCardWithID(key)
			} else {
				NewTextBox("RClick: Sell", 16, l + 32, t - 13, 96, 50).Draw(screen)
			}
		}
		col++
		if col == d.cols {
			col = 0
			row++
		}
	}
}
