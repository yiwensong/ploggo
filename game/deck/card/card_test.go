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
		0:  "K",
		1:  "A",
		2:  "2",
		3:  "3",
		4:  "4",
		5:  "5",
		6:  "6",
		7:  "7",
		8:  "8",
		9:  "9",
		10: "T",
		11: "J",
		12: "Q",
	}

	for rank, name := range ranksAndNames {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, rank.Name(), name)
		})
	}
}

func CardTest(t *testing.T) {
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
