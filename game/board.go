package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

var (
	firstCard  *Card
	secondCard *Card
)

type Board struct {
	cards []*Card
	size  int
}

func NewBoard(size int) (*Board, error) {
	b := &Board{
		size:  size,
		cards: make([]*Card, size*size),
	}

	var err error
	b.cards, err = newCardDeck(size)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (b Board) Size() (int, int) {
	x := b.size*cardSize + (b.size+1)*cardMargin
	y := x
	return x, y
}

func cardClicked(x int, y int, cards []*Card) *Card {
	for _, c := range cards {
		if !c.showing && c.IsClicked(x, y) {
			return c
		}
	}
	return nil
}

func (b Board) Draw(boardImage *ebiten.Image) {
	boardImage.Fill(frameColor)
	for j := 0; j < b.size; j++ {
		for i := 0; i < b.size; i++ {
			op := &ebiten.DrawImageOptions{}
			x := i*cardSize + (i+1)*cardMargin
			y := j*cardSize + (j+1)*cardMargin
			op.GeoM.Translate(float64(x), float64(y))
			r, g, b, a := colorToScale(color.NRGBA{0xee, 0xe4, 0xda, 0x59})
			op.ColorM.Scale(r, g, b, a)
			boardImage.DrawImage(cardImage, op)
		}
	}
	for _, c := range b.cards {
		c.Draw(boardImage)
	}
}
