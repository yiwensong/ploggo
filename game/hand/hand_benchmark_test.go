package hand

import (
	testing "testing"

	card "github.com/yiwensong/ploggo/game/deck/card"
)

func BenchmarkSmallSort(b *testing.B) {
	for i := 0; i < b.N; i++ {
		smallSort_VOLATILE([]int{25, 4, 33, 1, 2})
	}
}

func BenchmarkSortRanks(b *testing.B) {
	for i := 0; i < b.N; i++ {
		sortRanks_VOLATILE([]card.Rank{
			card.Rank(25),
			card.Rank(4),
			card.Rank(33),
			card.Rank(1),
			card.Rank(2),
		})
	}
}

func BenchmarkRankHand(b *testing.B) {
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
		b.Run(test.testName, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				RankHand(test.cards)
			}
		})
	}
}

func BenchmarkSimpleHandImpl_BestHand(b *testing.B) {
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
		b.Run(test.testName, func(b *testing.B) {
			evaluator := NewSimpleHandEvaluator(test.hole, test.community)
			for i := 0; i < b.N; i++ {
				evaluator.BestHand()
			}
		})
	}
}
