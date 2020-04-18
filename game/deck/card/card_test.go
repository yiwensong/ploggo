package deck

import (
	fmt "fmt"
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

func TestSuitNames(t *testing.T) {
	suitsAndNames := map[Suit]string{
		Club:    "Club",
		Diamond: "Diamond",
		Heart:   "Heart",
		Spade:   "Spade",
	}

	for suit, name := range suitsAndNames {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, suit.Name(), name)
		})
	}
}

func TestSuitShortNames(t *testing.T) {
	suitsAndNames := map[Suit]string{
		Club:    "c",
		Diamond: "d",
		Heart:   "h",
		Spade:   "s",
	}

	for suit, name := range suitsAndNames {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, suit.ShortName(), name)
		})
	}
}

func TestRankName(t *testing.T) {
	ranksAndNames := map[Rank]string{
		0:  "2",
		1:  "3",
		2:  "4",
		3:  "5",
		4:  "6",
		5:  "7",
		6:  "8",
		7:  "9",
		8:  "T",
		9:  "J",
		10: "Q",
		11: "K",
		12: "A",
	}

	for rank, name := range ranksAndNames {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, rank.Name(), name)
		})
	}
}

func TestCard(t *testing.T) {
	tests := []struct {
		cardNum int
		rank    Rank
		suit    Suit
	}{
		{
			cardNum: 0,
			rank:    0,
			suit:    Club,
		},
		{
			cardNum: 15,
			rank:    3,
			suit:    Spade,
		},
		{
			cardNum: 16,
			rank:    4,
			suit:    Club,
		},
		{
			cardNum: 44,
			rank:    11,
			suit:    Club,
		},
		{
			cardNum: 45,
			rank:    11,
			suit:    Diamond,
		},
		{
			cardNum: 51,
			rank:    12,
			suit:    Spade,
		},
		{
			cardNum: 6,
			rank:    1,
			suit:    Heart,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s%s", test.rank.Name(), test.suit.ShortName()), func(t *testing.T) {
			card := NewCard(test.cardNum)
			assert.Equal(t, test.suit, card.Suit())
			assert.Equal(t, test.rank, card.Rank())
		})
	}
}

func TestFromRankAndSuit(t *testing.T) {
	tests := []struct {
		rank Rank
		suit Suit
	}{
		{
			rank: 0,
			suit: Diamond,
		},
		{
			rank: 1,
			suit: Club,
		},
		{
			rank: 3,
			suit: Heart,
		},
		{
			rank: 7,
			suit: Spade,
		},
		{
			rank: 12,
			suit: Heart,
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%s%s", test.rank.Name(), test.suit.ShortName()), func(t *testing.T) {
			card := FromRankAndSuit(test.rank, test.suit)
			assert.Equal(t, test.suit, card.Suit())
			assert.Equal(t, test.rank, card.Rank())
		})
	}
}
