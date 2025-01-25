package hand

import (
	fmt "fmt"
	testing "testing"

	assert "github.com/stretchr/testify/assert"

	card "github.com/yiwensong/ploggo/game/deck/card"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		testName string
		rank1    *HandRank
		rank2    *HandRank
		expected int
	}{
		{
			testName: "SF > Quad",
			rank1: &HandRank{
				RankType: StraightFlush,
				Ranks:    []card.Rank{card.Rank(6)},
			},
			rank2: &HandRank{
				RankType: Quad,
				Ranks:    []card.Rank{card.Rank(12), card.Rank(2)},
			},
			expected: 1,
		},
		{
			testName: "Pair < FullHouse",
			rank1: &HandRank{
				RankType: Pair,
				Ranks:    []card.Rank{card.Rank(6), card.Rank(12), card.Rank(11), card.Rank(10)},
			},
			rank2: &HandRank{
				RankType: FullHouse,
				Ranks:    []card.Rank{card.Rank(12), card.Rank(6)},
			},
			expected: -1,
		},
		{
			testName: "ST = ST",
			rank1: &HandRank{
				RankType: Straight,
				Ranks:    []card.Rank{card.Rank(8)},
			},
			rank2: &HandRank{
				RankType: Straight,
				Ranks:    []card.Rank{card.Rank(8)},
			},
			expected: 0,
		},
		{
			testName: "Pair > Pair on kicker",
			rank1: &HandRank{
				RankType: Pair,
				Ranks:    []card.Rank{card.Rank(12), card.Rank(11), card.Rank(10), card.Rank(9)},
			},
			rank2: &HandRank{
				RankType: Pair,
				Ranks:    []card.Rank{card.Rank(12), card.Rank(11), card.Rank(10), card.Rank(8)},
			},
			expected: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			result := test.rank1.Compare(test.rank2)
			assert.Equal(t, test.expected, result)
		})
	}
}

func TestSuited_Success(t *testing.T) {
	cases := [][]int{
		{0, 4, 16, 20, 44},
		{1, 9, 13, 25, 49},
		{2, 6, 10, 14, 18},
		{3, 15, 19, 43, 51},
	}

	for _, handInt := range cases {
		hand := make([]card.Card, len(handInt))
		for i, cardNum := range handInt {
			hand[i] = card.NewCard(cardNum)
		}
		isSuited := suited(hand)
		assert.True(t, isSuited)
	}
}

func TestSuited_Failure(t *testing.T) {
	cases := [][]int{
		{0, 4, 17, 20, 44},
		{1, 9, 11, 25, 49},
		{2, 4, 9, 15, 19},
		{4, 15, 19, 43, 51},
	}

	for _, handInt := range cases {
		hand := make([]card.Card, len(handInt))
		for i, cardNum := range handInt {
			hand[i] = card.NewCard(cardNum)
		}
		isSuited := suited(hand)
		assert.False(t, isSuited)
	}
}

func TestSmallSort(t *testing.T) {
	tests := []struct {
		unsorted []int
		sorted   []int
	}{
		{
			unsorted: []int{5, 4, 3, 2, 1},
			sorted:   []int{5, 4, 3, 2, 1},
		},
		{
			unsorted: []int{25, 4, 33, 2, 1},
			sorted:   []int{33, 25, 4, 2, 1},
		},
		{
			unsorted: []int{1, 2, 3, 4},
			sorted:   []int{4, 3, 2, 1},
		},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%+v", test.unsorted), func(t *testing.T) {
			sorted := smallSort_VOLATILE(test.unsorted)
			assert.Equal(t, test.sorted, sorted)
		})
	}
}

func TestStraighted_Success(t *testing.T) {
	cases := [][]card.Rank{
		{
			card.Rank(9),
			card.Rank(8),
			card.Rank(7),
			card.Rank(6),
			card.Rank(5),
		},
		{
			card.Rank(12),
			card.Rank(11),
			card.Rank(10),
			card.Rank(9),
			card.Rank(8),
		},
		{
			card.Rank(12),
			card.Rank(3),
			card.Rank(2),
			card.Rank(1),
			card.Rank(0),
		},
	}

	for _, hand := range cases {
		t.Run(fmt.Sprintf("%+v", hand), func(t *testing.T) {
			isStraight := straighted(hand)
			assert.True(t, isStraight)
		})
	}
}

func TestStraighted_Failure(t *testing.T) {
	cases := [][]card.Rank{
		{
			card.Rank(10),
			card.Rank(8),
			card.Rank(7),
			card.Rank(6),
			card.Rank(5),
		},
		{
			card.Rank(12),
			card.Rank(11),
			card.Rank(9),
			card.Rank(9),
			card.Rank(8),
		},
		{
			card.Rank(12),
			card.Rank(3),
			card.Rank(2),
			card.Rank(1),
			card.Rank(1),
		},
	}

	for _, hand := range cases {
		t.Run(fmt.Sprintf("%+v", hand), func(t *testing.T) {
			isStraight := straighted(hand)
			assert.False(t, isStraight)
		})
	}
}

func TestRankDupes(t *testing.T) {
	tests := []struct {
		testName string
		hand     []card.Rank
		expected map[card.Rank]int
	}{
		{
			testName: "pair",
			hand: []card.Rank{
				card.Rank(12),
				card.Rank(12),
				card.Rank(11),
				card.Rank(10),
				card.Rank(9),
			},
			expected: map[card.Rank]int{
				card.Rank(12): 2,
			},
		},
		{
			testName: "two pair",
			hand: []card.Rank{
				card.Rank(12),
				card.Rank(12),
				card.Rank(9),
				card.Rank(9),
				card.Rank(4),
			},
			expected: map[card.Rank]int{
				card.Rank(12): 2,
				card.Rank(9):  2,
			},
		},
		{
			testName: "trips",
			hand: []card.Rank{
				card.Rank(12),
				card.Rank(11),
				card.Rank(9),
				card.Rank(9),
				card.Rank(9),
			},
			expected: map[card.Rank]int{
				card.Rank(9): 3,
			},
		},
		{
			testName: "quads",
			hand: []card.Rank{
				card.Rank(0),
				card.Rank(3),
				card.Rank(3),
				card.Rank(3),
				card.Rank(3),
			},
			expected: map[card.Rank]int{
				card.Rank(3): 4,
			},
		},
		{
			testName: "full house",
			hand: []card.Rank{
				card.Rank(1),
				card.Rank(1),
				card.Rank(6),
				card.Rank(6),
				card.Rank(6),
			},
			expected: map[card.Rank]int{
				card.Rank(6): 3,
				card.Rank(1): 2,
			},
		},
		{
			testName: "no pair",
			hand: []card.Rank{
				card.Rank(1),
				card.Rank(2),
				card.Rank(8),
				card.Rank(9),
				card.Rank(12),
			},
			expected: map[card.Rank]int{},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			dupes := rankDupes(test.hand)
			assert.Equal(t, test.expected, dupes)
		})
	}
}

func TestRankHand(t *testing.T) {
	tests := []struct {
		testName string
		cards    []card.Card
		expected *HandRank
	}{
		{
			testName: "high card",
			cards: []card.Card{
				card.NewCard(51),
				card.NewCard(44),
				card.NewCard(30),
				card.NewCard(21),
				card.NewCard(11),
			},
			expected: &HandRank{
				RankType: HighCard,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(7),
					card.Rank(5),
					card.Rank(2),
				},
			},
		},
		{
			testName: "pair",
			cards: []card.Card{
				card.NewCard(51),
				card.NewCard(49),
				card.NewCard(30),
				card.NewCard(21),
				card.NewCard(11),
			},
			expected: &HandRank{
				RankType: Pair,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(7),
					card.Rank(5),
					card.Rank(2),
				},
			},
		},
		{
			testName: "two pair",
			cards: []card.Card{
				card.NewCard(51),
				card.NewCard(49),
				card.NewCard(30),
				card.NewCard(31),
				card.NewCard(11),
			},
			expected: &HandRank{
				RankType: TwoPair,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(7),
					card.Rank(2),
				},
			},
		},
		{
			testName: "triples",
			cards: []card.Card{
				card.NewCard(31),
				card.NewCard(30),
				card.NewCard(29),
				card.NewCard(21),
				card.NewCard(11),
			},
			expected: &HandRank{
				RankType: Triple,
				Ranks: []card.Rank{
					card.Rank(7),
					card.Rank(5),
					card.Rank(2),
				},
			},
		},
		{
			testName: "straight",
			cards: []card.Card{
				card.NewCard(31),
				card.NewCard(25),
				card.NewCard(20),
				card.NewCard(16),
				card.NewCard(13),
			},
			expected: &HandRank{
				RankType: Straight,
				Ranks: []card.Rank{
					card.Rank(7),
				},
			},
		},
		{
			testName: "wheel",
			cards: []card.Card{
				card.NewCard(51),
				card.NewCard(12),
				card.NewCard(8),
				card.NewCard(4),
				card.NewCard(1),
			},
			expected: &HandRank{
				RankType: Straight,
				Ranks: []card.Rank{
					card.Rank(3),
				},
			},
		},
		{
			testName: "flush",
			cards: []card.Card{
				card.NewCard(48),
				card.NewCard(40),
				card.NewCard(16),
				card.NewCard(12),
				card.NewCard(8),
			},
			expected: &HandRank{
				RankType: Flush,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(10),
					card.Rank(4),
					card.Rank(3),
					card.Rank(2),
				},
			},
		},
		{
			testName: "full house",
			cards: []card.Card{
				card.NewCard(48),
				card.NewCard(49),
				card.NewCard(50),
				card.NewCard(12),
				card.NewCard(13),
			},
			expected: &HandRank{
				RankType: FullHouse,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(3),
				},
			},
		},
		{
			testName: "full house small over big",
			cards: []card.Card{
				card.NewCard(48),
				card.NewCard(49),
				card.NewCard(14),
				card.NewCard(12),
				card.NewCard(13),
			},
			expected: &HandRank{
				RankType: FullHouse,
				Ranks: []card.Rank{
					card.Rank(3),
					card.Rank(12),
				},
			},
		},
		{
			testName: "quads",
			cards: []card.Card{
				card.NewCard(48),
				card.NewCard(15),
				card.NewCard(14),
				card.NewCard(12),
				card.NewCard(13),
			},
			expected: &HandRank{
				RankType: Quad,
				Ranks: []card.Rank{
					card.Rank(3),
					card.Rank(12),
				},
			},
		},
		{
			testName: "straight flush",
			cards: []card.Card{
				card.NewCard(27),
				card.NewCard(23),
				card.NewCard(19),
				card.NewCard(15),
				card.NewCard(11),
			},
			expected: &HandRank{
				RankType: StraightFlush,
				Ranks: []card.Rank{
					card.Rank(6),
				},
			},
		},
		{
			testName: "straight flush - wheel",
			cards: []card.Card{
				card.NewCard(48),
				card.NewCard(12),
				card.NewCard(8),
				card.NewCard(4),
				card.NewCard(0),
			},
			expected: &HandRank{
				RankType: StraightFlush,
				Ranks: []card.Rank{
					card.Rank(3),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			rank, err := RankHand(test.cards)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, rank)
		})
	}
}

func TestSimpleHandImpl_BestHand(t *testing.T) {
	tests := []struct {
		testName  string
		hole      []card.Card
		community []card.Card
		expected  *HandRank
	}{
		{
			testName: "high card",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(11), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(10), card.Diamond),
				card.FromRankAndSuit(card.Rank(0), card.Diamond),
				card.FromRankAndSuit(card.Rank(1), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: HighCard,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(10),
					card.Rank(7),
					card.Rank(2),
				},
			},
		},
		{
			testName: "pair with board",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(11), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(0), card.Diamond),
				card.FromRankAndSuit(card.Rank(1), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: Pair,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(7),
					card.Rank(2),
				},
			},
		},
		{
			testName: "pair in hand",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(0), card.Diamond),
				card.FromRankAndSuit(card.Rank(1), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: Pair,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(7),
					card.Rank(2),
				},
			},
		},
		{
			testName: "pair on board",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(0), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(1), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: Pair,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(9),
					card.Rank(7),
				},
			},
		},
		{
			testName: "two pair with board",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(0), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(1), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: TwoPair,
				Ranks: []card.Rank{
					card.Rank(11),
					card.Rank(9),
					card.Rank(7),
				},
			},
		},
		{
			testName: "two pair hand and board",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(11), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(2), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(1), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: TwoPair,
				Ranks: []card.Rank{
					card.Rank(11),
					card.Rank(2),
					card.Rank(9),
				},
			},
		},
		{
			testName: "set",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(11), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(0), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(11), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: Triple,
				Ranks: []card.Rank{
					card.Rank(11),
					card.Rank(9),
					card.Rank(7),
				},
			},
		},
		{
			testName: "paired board trip",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(10), card.Diamond),
				card.FromRankAndSuit(card.Rank(6), card.Club),
				card.FromRankAndSuit(card.Rank(1), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(11), card.Club),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Spade),
			},
			expected: &HandRank{
				RankType: Triple,
				Ranks: []card.Rank{
					card.Rank(11),
					card.Rank(10),
					card.Rank(9),
				},
			},
		},
		{
			testName: "straight",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(10), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Spade),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: Straight,
				Ranks: []card.Rank{
					card.Rank(12),
				},
			},
		},
		{
			testName: "not a straight (too many cards in hand)",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(10), card.Diamond),
				card.FromRankAndSuit(card.Rank(9), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(4), card.Spade),
				card.FromRankAndSuit(card.Rank(0), card.Spade),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: HighCard,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(10),
					card.Rank(7),
					card.Rank(4),
				},
			},
		},
		{
			testName: "not a straight (too many cards on board)",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(4), card.Diamond),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(0), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(10), card.Spade),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: HighCard,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(10),
					card.Rank(9),
					card.Rank(4),
				},
			},
		},
		{
			testName: "wheel straight",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(3), card.Diamond),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(1), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(0), card.Diamond),
				card.FromRankAndSuit(card.Rank(2), card.Spade),
				card.FromRankAndSuit(card.Rank(1), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: Straight,
				Ranks: []card.Rank{
					card.Rank(3),
				},
			},
		},
		{
			testName: "flush",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(1), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(0), card.Spade),
				card.FromRankAndSuit(card.Rank(2), card.Spade),
				card.FromRankAndSuit(card.Rank(1), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: Flush,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(11),
					card.Rank(2),
					card.Rank(1),
					card.Rank(0),
				},
			},
		},
		{
			testName: "full house",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(1), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Club),
				card.FromRankAndSuit(card.Rank(8), card.Spade),
				card.FromRankAndSuit(card.Rank(1), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: FullHouse,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(8),
				},
			},
		},
		{
			testName: "quads",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(12), card.Diamond),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(1), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Club),
				card.FromRankAndSuit(card.Rank(12), card.Heart),
				card.FromRankAndSuit(card.Rank(1), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: Quad,
				Ranks: []card.Rank{
					card.Rank(12),
					card.Rank(8),
				},
			},
		},
		{
			testName: "straight flush",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(11), card.Spade),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(1), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(10), card.Spade),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: StraightFlush,
				Ranks: []card.Rank{
					card.Rank(12),
				},
			},
		},
		{
			testName: "straight flush wheel",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(12), card.Spade),
				card.FromRankAndSuit(card.Rank(0), card.Spade),
				card.FromRankAndSuit(card.Rank(2), card.Club),
				card.FromRankAndSuit(card.Rank(1), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(1), card.Spade),
				card.FromRankAndSuit(card.Rank(2), card.Spade),
				card.FromRankAndSuit(card.Rank(3), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
			},
			expected: &HandRank{
				RankType: StraightFlush,
				Ranks: []card.Rank{
					card.Rank(3),
				},
			},
		},
		{
			testName: "Full house with straight",
			hole: []card.Card{
				card.FromRankAndSuit(card.Rank(7), card.Spade),
				card.FromRankAndSuit(card.Rank(7), card.Diamond),
				card.FromRankAndSuit(card.Rank(6), card.Club),
				card.FromRankAndSuit(card.Rank(5), card.Club),
			},
			community: []card.Card{
				card.FromRankAndSuit(card.Rank(2), card.Heart),
				card.FromRankAndSuit(card.Rank(2), card.Spade),
				card.FromRankAndSuit(card.Rank(9), card.Spade),
				card.FromRankAndSuit(card.Rank(8), card.Club),
				card.FromRankAndSuit(card.Rank(7), card.Club),
			},
			expected: &HandRank{
				RankType: FullHouse,
				Ranks: []card.Rank{
					card.Rank(7),
					card.Rank(2),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			evaluator := NewSimpleHandEvaluator(test.hole, test.community)
			rank, err := evaluator.BestHand()
			assert.NoError(t, err)
			assert.Equal(t, test.expected, rank)
		})
	}
}
