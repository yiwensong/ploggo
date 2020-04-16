package deck

import (
	fmt "fmt"

	errors "github.com/pkg/errors"

	card "github.com/yiwensong/ploggo/game/deck/card"
	rng "github.com/yiwensong/ploggo/game/deck/rng"
)

var _ card.Card = nil

type Deck interface {
	// Shuffles the deck
	Shuffle() error

	// Draws one card from the deck
	Draw() (card.Card, error)

	// Draws N cards from the deck
	DrawN(int) ([]card.Card, error)

	// Shows the top card
	Top() (card.Card, error)

	// Duplicates and randomizes the duplication
	Duplicate() (Deck, error)
}

// This is the default implementaiton of the Deck interface
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
		// Don't shuffle cards that are already drawn
		if i < d.drawn {
			continue
		}
		swap := d.RNG.RandInt(len(d.order)-i) + i
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
	if d.drawn+n > len(d.order) {
		return nil, fmt.Errorf(
			"Not enough cards left in deck (needed: %d, has: %d)",
			n,
			len(d.order)-d.drawn,
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
	if d.drawn >= len(d.order) {
		return nil, fmt.Errorf("No cards left in deck")
	}
	card := card.NewCard(d.order[d.drawn])
	return card, nil
}

// Duplicates the deck, and shuffles the duplicated deck to avoid
// deterministic results from the duplicates
func (d *DeckImpl) Duplicate() (Deck, error) {
	newOrder := make([]int, len(d.order))
	for i, value := range d.order {
		newOrder[i] = value
	}
	newDeck := &DeckImpl{
		order: newOrder,
		drawn: d.drawn,
		RNG:   d.RNG,
	}
	newDeck.Shuffle()
	return newDeck, nil
}

// Type checking
var _ Deck = (*DeckImpl)(nil)
