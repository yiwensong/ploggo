package deck

import (
	testing "testing"

	"runtime"

	mathrand "math/rand"

	assert "github.com/stretchr/testify/assert"
	card "github.com/yiwensong/ploggo/game/deck/card"
	rng "github.com/yiwensong/ploggo/game/deck/rng"
)

var _ *runtime.Func = nil

func TestShuffle(t *testing.T) {
	deck := NewDeck()
	deck.RNG = &rng.SeededRNG{
		Rand: mathrand.New(mathrand.NewSource(12)),
	}

	deck.Shuffle()

	// This is the order that we expect this seed to come out as
	expectedOrder := []int{
		37, 26, 10, 6, 18, 29, 43, 45, 39, 40, 46, 23, 7,
		28, 9, 21, 49, 19, 25, 14, 33, 38, 32, 2, 27, 41,
		24, 22, 30, 42, 8, 3, 51, 20, 5, 15, 44, 50, 47,
		11, 13, 0, 35, 16, 31, 12, 34, 17, 36, 48, 1, 4,
	}
	assert.Equal(t, expectedOrder, deck.order)
}

func TestShuffle_Partial(t *testing.T) {
	deck := NewDeck()
	deck.RNG = &rng.SeededRNG{
		Rand: mathrand.New(mathrand.NewSource(12)),
	}

	drawn, err := deck.DrawN(26)
	assert.NotNil(t, drawn)
	assert.NoError(t, err)

	deck.Shuffle()
	expectedOrder := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
		13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25,
		37, 45, 38, 48, 46, 26, 35, 44, 43, 47, 29, 40, 49,
		41, 50, 34, 27, 33, 32, 51, 39, 30, 36, 42, 28, 31,
	}
	assert.Equal(t, expectedOrder, deck.order)
}

func TestDraw(t *testing.T) {
	deck := NewDeck()
	deck.RNG = &rng.SeededRNG{
		Rand: mathrand.New(mathrand.NewSource(12)),
	}

	deck.Shuffle()

	// This is the order that we expect this seed to come out as
	expectedOrder := []int{
		37, 26, 10, 6, 18, 29, 43, 45, 39, 40, 46, 23, 7,
		28, 9, 21, 49, 19, 25, 14, 33, 38, 32, 2, 27, 41,
		24, 22, 30, 42, 8, 3, 51, 20, 5, 15, 44, 50, 47,
		11, 13, 0, 35, 16, 31, 12, 34, 17, 36, 48, 1, 4,
	}

	for _, expectedCardId := range expectedOrder {
		expectedCard := card.NewCard(expectedCardId)
		drawn, err := deck.Draw()
		assert.NoError(t, err)
		assert.Equal(t, expectedCard, drawn)
	}
}

func TestDrawN(t *testing.T) {
	deck := NewDeck()
	deck.RNG = &rng.SeededRNG{
		Rand: mathrand.New(mathrand.NewSource(12)),
	}

	deck.Shuffle()

	// This is the order that we expect this seed to come out as
	expectedOrder := []int{
		37, 26, 10, 6, 18, 29, 43, 45, 39, 40, 46, 23, 7,
		28, 9, 21, 49, 19, 25, 14, 33, 38, 32, 2, 27, 41,
		24, 22, 30, 42, 8, 3, 51, 20, 5, 15, 44, 50, 47,
		11, 13, 0, 35, 16, 31, 12, 34, 17, 36, 48, 1, 4,
	}

	n := 5
	expectedCards := make([]card.Card, n)
	for i := range expectedCards {
		expectedCards[i] = card.NewCard(expectedOrder[i])
	}

	drawn, err := deck.DrawN(5)
	assert.NoError(t, err)

	for i, cd := range drawn {
		assert.Equal(t, expectedCards[i], cd)
	}

	expNextCard := card.NewCard(expectedOrder[n])
	nextCard, err := deck.Top()
	assert.NoError(t, err)
	assert.Equal(t, expNextCard, nextCard)
}

func TestTop(t *testing.T) {
	deck := NewDeck()

	expectedCard := card.NewCard(0)

	card, err := deck.Top()
	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)

	card, err = deck.Top()
	assert.NoError(t, err)
	assert.Equal(t, expectedCard, card)
}

func TestTop_TooManyDrawn(t *testing.T) {
	deck := NewDeck()
	deck.drawn = 52

	card, err := deck.Top()
	assert.Nil(t, card)
	assert.Error(t, err)
}

func TestDraw_TooManyDrawn(t *testing.T) {
	deck := NewDeck()
	deck.drawn = 52

	card, err := deck.Draw()
	assert.Nil(t, card)
	assert.Error(t, err)
}

func TestDrawN_InsufficientCardsTest(t *testing.T) {
	tests := []struct {
		testName    string
		cardsDrawn  int
		toDraw      int
		shouldError bool
	}{
		{
			testName:    "draw 5 has 5",
			cardsDrawn:  52 - 5,
			toDraw:      5,
			shouldError: false,
		},
		{
			testName:    "draw 5 has 6",
			cardsDrawn:  52 - 6,
			toDraw:      5,
			shouldError: false,
		},
		{
			testName:    "draw 5 has 4",
			cardsDrawn:  52 - 4,
			toDraw:      5,
			shouldError: true,
		},
		{
			testName:    "draw 0 has 52",
			cardsDrawn:  0,
			toDraw:      0,
			shouldError: false,
		},
		{
			testName:    "draw 0 has 0",
			cardsDrawn:  52,
			toDraw:      0,
			shouldError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			deck := NewDeck()
			deck.drawn = test.cardsDrawn

			_, err := deck.DrawN(test.toDraw)
			if test.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestDuplicate(t *testing.T) {
	deck := NewDeck()
	deck.RNG = &rng.SeededRNG{
		Rand: mathrand.New(mathrand.NewSource(12)),
	}

	// Draw the first 13 cards
	drawn, err := deck.DrawN(13)
	assert.NoError(t, err)
	assert.NotNil(t, drawn)

	newDeck, err := deck.Duplicate()
	assert.NoError(t, err)

	expectedOrder := []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12,
		37, 39, 46, 34, 31, 16, 43, 35, 22, 30, 48, 29, 21,
		28, 40, 33, 49, 42, 18, 17, 47, 20, 19, 50, 44, 13,
		27, 14, 45, 25, 41, 51, 23, 32, 24, 36, 15, 26, 38,
	}
	assert.Equal(t, expectedOrder, newDeck.(*DeckImpl).order)
	assert.Equal(t, deck.drawn, newDeck.(*DeckImpl).drawn)

	// Make sure our original order doesn't change
	for i, cardN := range deck.order {
		assert.Equal(t, i, cardN)
	}
}
