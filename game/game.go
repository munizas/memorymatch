package game

import (
	"fmt"
	"image/color"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	ScreenWidth  = 420
	ScreenHeight = 600
	boardSize    = 4
	xOffset      = (ScreenWidth - (boardSize*cardSize + (boardSize+1)*cardMargin)) / 2
	yOffset      = (ScreenHeight - (boardSize*cardSize + (boardSize+1)*cardMargin)) / 2
	NormalGame   = 0
	HardGame     = 1
)

var (
	matchAttempts   int
	gameMessage     string
	gameModeMessage string
)

type GameState int

const (
	FirstSelection GameState = iota
	SecondSelection
	CheckMatch
	Finished
)

type Game struct {
	state      GameState
	board      *Board
	boardImage *ebiten.Image
	mode       int
}

func NewGame() (*Game, error) {
	g := &Game{state: FirstSelection, mode: NormalGame}
	gameModeMessage = "Normal"
	var err error
	g.board, err = NewBoard(boardSize)
	if err != nil {
		return nil, err
	}
	return g, nil
}

func (g *Game) resetGame(gameMode int) error {
	g.state = FirstSelection
	g.mode = gameMode
	matchAttempts = 0
	gameMessage = ""
	if gameMode == HardGame {
		gameModeMessage = "Hard"
	} else {
		gameModeMessage = "Normal"
	}
	var err error
	g.board, err = NewBoard(boardSize)
	if err != nil {
		return err
	}
	return nil
}

func (g *Game) Update() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyN) {
		if err := g.resetGame(NormalGame); err != nil {
			panic(err)
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyH) {
		if err := g.resetGame(HardGame); err != nil {
			panic(err)
		}
	}

	switch g.state {
	case FirstSelection:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			firstCard = cardClicked(x, y, g.board.cards)
			if firstCard != nil {
				firstCard.showing = true
				g.state = SecondSelection
			}
		}
	case SecondSelection:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			x, y := ebiten.CursorPosition()
			secondCard = cardClicked(x, y, g.board.cards)
			if secondCard != nil {
				secondCard.showing = true
				g.state = CheckMatch
			}
		}
	case CheckMatch:
		time.Sleep(300 * time.Millisecond)
		matchAttempts++
		if firstCard.value != secondCard.value {
			g.resetCards()
		}
		if isGameFinished(g.board) {
			g.state = Finished
		} else {
			g.state = FirstSelection
		}
	case Finished:
		gameMessage = "Nice! You did it!"
	}
	return nil
}

func (g *Game) resetCards() {
	if g.mode == NormalGame {
		firstCard.showing = false
		secondCard.showing = false
	} else {
		for _, c := range g.board.cards {
			c.showing = false
		}
	}
}

func isGameFinished(board *Board) bool {
	for _, c := range board.cards {
		if !c.showing {
			return false
		}
	}
	return true
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (g Game) Draw(screen *ebiten.Image) {
	if g.boardImage == nil {
		w, h := g.board.Size()
		g.boardImage = ebiten.NewImage(w, h)
	}
	screen.Fill(backgroundColor)
	g.board.Draw(g.boardImage)
	op := &ebiten.DrawImageOptions{}
	sw, sh := screen.Size()
	bw, bh := g.boardImage.Size()
	x := (sw - bw) / 2
	y := (sh - bh) / 2
	op.GeoM.Translate(float64(x), float64(y))
	screen.DrawImage(g.boardImage, op)

	text.Draw(screen, fmt.Sprintf("Game Mode: %v", gameModeMessage), mplusMiniFont, xOffset, yOffset-10, color.RGBA{0x77, 0x6e, 0x65, 0xff})
	text.Draw(screen, fmt.Sprintf("Match Attempts: %d", matchAttempts), mplusSmallFont, xOffset, sh-xOffset*2, color.RGBA{0x77, 0x6e, 0x65, 0xff})
	text.Draw(screen, fmt.Sprintf("Press N for new game"), mplusMiniFont, xOffset, sh-xOffset, color.RGBA{0x77, 0x6e, 0x65, 0xff})
	text.Draw(screen, fmt.Sprintf("Press H for hard mode game"), mplusMiniFont, xOffset, sh-xOffset/2, color.RGBA{0x77, 0x6e, 0x65, 0xff})
	text.Draw(screen, gameMessage, mplusNormalFont, xOffset, xOffset*2, color.RGBA{0x77, 0x6e, 0x65, 0xff})
}
