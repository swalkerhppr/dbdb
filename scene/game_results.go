package scene

import (
	"dbdb/components"
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/joelschutz/stagehand"
)

type gameResults struct {
	mainMenuButton *components.Button
	won bool
	scoreText string
	score int
	*BaseScene
}

func CreateGameResults(width, height int) stagehand.Scene[*State] {
	// If day > 5 you lose. If all requirements are met you win!
	gr := &gameResults{
		BaseScene: NewBaseWithFill(width, height, color.RGBA{R: 0x73, G: 0x7b, B: 0x9a, A: 0xFF}), 
		won : true,
	}

	return gr
}

func (gr *gameResults) Load(s *State, controller stagehand.SceneController[*State]) {
	gr.mainMenuButton = components.NewButton("Main Menu", 320, 436, &s.Controls, func() {
		gr.SceneManager.SwitchWithTransition(SceneMap[MainMenu], stagehand.NewDurationTimedFadeTransition[*State](time.Millisecond * 100))
	})
	gr.won = s.RequiredBoardParts == s.BoardPartsBuilt && s.RequiredPlankParts == s.PlankPartsBuilt
	gr.scoreText = createScoreText(s)
	gr.score = (s.PlankPartsBuilt + s.BoardPartsBuilt) * 1000 - (s.Day * 500) + int(s.MoneyLeft) - (len(s.Deck) * 10)
	gr.BaseScene.Load(s, controller)
}

func (gr *gameResults) Draw(screen *ebiten.Image) {
	gr.DrawScene(screen)
	text := "You ran out of time."
	if gr.won {
		text = "You Won!\nCongratulations,\nEnjoy your new deck!\n\n\n\n\n\n\n\n\n"
	} else {
		text += "\nDangit!\nNo deck this time.\n\n\n\n\n\n\n\n\n\n"
	}
	text += "Thanks for playing!"
	components.NewTextBox(text, 24, 18, 16, 608, 448).Draw(screen)
	components.NewDeckIndicators(gr.State, 8, 164, 163).Draw(screen)
	components.NewTextBox(gr.scoreText, 16, 191, 240, 256, 96).Draw(screen)
	components.NewTextBox(fmt.Sprintf("Final Score: %d", gr.score), 23, 207, 338, 224, 40).Draw(screen)
	gr.mainMenuButton.Draw(screen)

	if gr.State.Controls.KeySpace {
		gr.mainMenuButton.OnClick()
	}
	if gr.State.Controls.KeyTab || gr.State.Controls.KeyEnter {
		takeScreenshot(fmt.Sprintf("dbdb_score_%d.png", gr.score), screen)
	}
}

func createScoreText(s *State) string {
	plankText := fmt.Sprintf("%d/%d", s.PlankPartsBuilt, s.RequiredPlankParts)
	boardText := fmt.Sprintf("%d/%d", s.BoardPartsBuilt, s.RequiredBoardParts)
	daysText := fmt.Sprintf("%d/%d", s.Day, s.MaxDays)
	moneyText := fmt.Sprintf("$%.0f", s.MoneyLeft)
	cardsLeftText := fmt.Sprintf("%d", len(s.Deck))
	return fmt.Sprintf("Plank Parts Built:%7s\nBoard Parts Built:%7s\nDays Taken:%14s\nMoney Left:%14s\nCards In Deck:%11s", plankText, boardText, daysText, moneyText, cardsLeftText)
}
