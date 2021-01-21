package game

import (
	"image/color"
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	cardSize   = 80
	cardMargin = 4
)

var (
	mplusSmallFont  font.Face
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

var (
	cardImage = ebiten.NewImage(cardSize, cardSize)
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusSmallFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	cardImage.Fill(color.White)
}

type Card struct {
	value   int
	x       int
	y       int
	matched bool
	showing bool
}

func colorToScale(clr color.Color) (float64, float64, float64, float64) {
	r, g, b, a := clr.RGBA()
	rf := float64(r) / 0xffff
	gf := float64(g) / 0xffff
	bf := float64(b) / 0xffff
	af := float64(a) / 0xffff
	// Convert to non-premultiplied alpha components.
	if 0 < af {
		rf /= af
		gf /= af
		bf /= af
	}
	return rf, gf, bf, af
}

func newCardDeck(size int) ([]*Card, error) {
	rand.Seed(time.Now().Unix())
	var numbers []int
	for i := 0; i < size*size/2; i++ {
		num := 1 + rand.Intn(10)
		numbers = append(numbers, num)
		numbers = append(numbers, num)
	}
	rand.Shuffle(len(numbers), func(i, j int) { numbers[i], numbers[j] = numbers[j], numbers[i] })

	cards := make([]*Card, size*size)
	for i, n := range numbers {
		cards[i] = &Card{
			value:   n,
			x:       i % size,
			y:       i / size,
			matched: false,
			showing: false,
		}
	}

	return cards, nil
}

func (c Card) Draw(boardImage *ebiten.Image) {
	if !c.showing {
		return
	}

	i, j := c.x, c.y
	v := c.value

	op := &ebiten.DrawImageOptions{}
	x := i*cardSize + (i+1)*cardMargin
	y := j*cardSize + (j+1)*cardMargin

	op.GeoM.Translate(float64(x), float64(y))
	r, g, b, a := colorToScale(cardBackgroundColor(v))
	op.ColorM.Scale(r, g, b, a)
	boardImage.DrawImage(cardImage, op)
	str := strconv.Itoa(v)

	f := mplusBigFont
	switch {
	case 3 < len(str):
		f = mplusSmallFont
	case 2 < len(str):
		f = mplusNormalFont
	}

	bound, _ := font.BoundString(f, str)
	w := (bound.Max.X - bound.Min.X).Ceil()
	h := (bound.Max.Y - bound.Min.Y).Ceil()
	x = x + (cardSize-w)/2
	y = y + (cardSize-h)/2 + h
	text.Draw(boardImage, str, f, x, y, cardColor(v))
}

func (c Card) IsClicked(x, y int) bool {
	dx := c.x*cardSize + (c.x+1)*cardMargin + 40 + cardSize
	dy := c.y*cardSize + (c.y+1)*cardMargin + 140 + cardSize

	if (x > dx-cardSize && x < dx) && (y > dy-cardSize && y < dy) {
		return true
	}
	return false
}
