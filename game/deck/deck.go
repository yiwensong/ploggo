package deck

import "fmt"

type Suit int

const (
	Club Suit = iota
	Diamond
	Heart
	Spade
)
const NUM_SUITS = 4

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
	if 2 <= r && r <= 9 {
		return fmt.Sprintf("%d", r)
	}
	switch r {
	case 1:
		return "A"
	case 10:
		return "T"
	case 11:
		return "J"
	case 12:
		return "Q"
	default: // King is canonically zero
		return "K"
	}
}

type Card interface {
	GetSuit() Suit
	GetRank() Rank
	String() string
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

func (c *CardImpl) Suit() Suit {
	return Suit(c.cardNum % NUM_SUITS)
}

func (c *CardImpl) Rank() Rank {
	return Rank(c.cardNum / NUM_SUITS)
}
