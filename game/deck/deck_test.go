package deck

import (
	testing "testing"

	mathrand "math/rand"

	assert "github.com/stretchr/testify/assert"
	card "github.com/yiwensong/ploggo/game/deck/card"
	rng "github.com/yiwensong/ploggo/game/deck/rng"
)

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
