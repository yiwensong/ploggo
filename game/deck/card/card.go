package deck

import (
	fmt "fmt"
)

const (
	// Number of suits
	NUM_SUITS = 4

	// Number of cards in the deck
	DECK_SIZE = 52
)

type Suit int

const (
	Club Suit = iota
	Diamond
	Heart
	Spade
)

func (s Suit) Name() string {
	switch s {
	case Club:
		return "Club"
	case Diamond:
		return "Diamond"
	case Heart:
		return "Heart"
	default:
		return "Spade"
	}
}

func (s Suit) ShortName() string {
	switch s {
	case Club:
		return "c"
	case Diamond:
		return "d"
	case Heart:
		return "h"
	default:
		return "s"
	}
}

type Rank int

func (r Rank) Name() string {
	if 0 <= r && r <= 7 {
		return fmt.Sprintf("%d", r+2)
	}
	switch r {
	case 8:
		return "T"
	case 9:
		return "J"
	case 10:
		return "Q"
	case 11:
		return "K"
	default:
		return "A"
	}
}

type Card interface {
	Suit() Suit
	Rank() Rank
	String() string
	CardNum() int
}

type CardImpl struct {
	cardNum int
}

// Takes a card number and turns it into a real card
func NewCard(cardNum int) *CardImpl {
	return &CardImpl{
		cardNum: cardNum,
	}
}

// Takes a rank and suit and returns the card
func FromRankAndSuit(rank Rank, suit Suit) *CardImpl {
	return &CardImpl{
		cardNum: int(rank)*NUM_SUITS + int(suit),
	}
}

func (c *CardImpl) Suit() Suit {
	return Suit(c.cardNum % NUM_SUITS)
}

func (c *CardImpl) Rank() Rank {
	return Rank(c.cardNum / NUM_SUITS)
}

func (c *CardImpl) String() string {
	return fmt.Sprintf("%s%s", c.Rank().Name(), c.Suit().ShortName())
}

func (c *CardImpl) CardNum() int {
	return c.cardNum
}
