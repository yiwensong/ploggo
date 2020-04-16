package deck

import (
	fmt "fmt"

	errors "github.com/pkg/errors"

	card "github.com/yiwensong/ploggo/game/deck/card"
	rng "github.com/yiwensong/ploggo/game/deck/rng"
)

var _ card.Card = nil

type Deck interface {
	Shuffle() error
	Draw() (card.Card, error)
	DrawN(int) ([]card.Card, error)
	Top() (card.Card, error)
}

// This deck is optimized for not getting shuffled in the middle
type DeckImpl struct {
	order []int
	drawn int

	// Used for shuffle function
	RNG rng.RNG
}

// Returns a completely new deck with all the cards in order
func NewDeck() *DeckImpl {
	order := make([]int, 52)
	for i := range order {
		order[i] = i
	}
	drawn := 0

	return &DeckImpl{
		order: order,
		drawn: drawn,
		RNG:   &rng.SecureRNG{},
	}
}

// Returns a shuffled deck
func NewShuffledDeck() *DeckImpl {
	deck := NewDeck()
	deck.Shuffle()
	return deck
}

// Shuffles the deck
func (d *DeckImpl) Shuffle() error {
	for i := range d.order {
		swap := d.RNG.RandInt(card.DECK_SIZE-i) + i
		currentValue := d.order[i]
		d.order[i] = d.order[swap]
		d.order[swap] = currentValue
	}
	return nil
}

// Draws one card
func (d *DeckImpl) Draw() (card.Card, error) {
	top, err := d.Top()
	if err != nil {
		return nil, errors.Wrapf(err, "Top")
	}
	d.drawn += 1
	return top, nil
}

// Draws N cards
func (d *DeckImpl) DrawN(n int) ([]card.Card, error) {
	if n < 0 {
		return nil, fmt.Errorf("You must draw a non-negative number of cards")
	}
	if d.drawn+n > card.DECK_SIZE {
		return nil, fmt.Errorf(
			"Not enough cards left in deck (needed: %d, has: %d)",
			n,
			card.DECK_SIZE-d.drawn,
		)
	}
	cards := make([]card.Card, n)
	for i := range cards {
		card_, err := d.Draw()
		if err != nil {
			return nil, errors.Wrapf(
				err,
				"Draw",
			)
		}
		cards[i] = card_
	}
	return cards, nil
}

// Returns the top card without drawing
func (d *DeckImpl) Top() (card.Card, error) {
	if d.drawn >= card.DECK_SIZE {
		return nil, fmt.Errorf("No cards left in deck")
	}
	card := card.NewCard(d.order[d.drawn])
	return card, nil
}

// Type checking
var _ Deck = (*DeckImpl)(nil)
