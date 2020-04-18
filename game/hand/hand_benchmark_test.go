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
