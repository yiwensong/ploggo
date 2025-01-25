package hand

import (
	fmt "fmt"

	errors "github.com/pkg/errors"

	card "github.com/yiwensong/ploggo/game/deck/card"
)

type RankType int

const (
	Invalid RankType = iota
	HighCard
	Pair
	TwoPair
	Triple
	Straight
	Flush
	FullHouse
	Quad
	StraightFlush
)

type HandRank struct {
	// The type of the hand rank
	RankType RankType

	// A list of the ranks of the cards in the hands, in sorted order
	// of importance
	//
	// Only important ranks will be given. For example, straights will
	// only have the highest card, and quads will only give the quaded
	// card and the kicker.
	Ranks []card.Rank
}

// Returns -1 if the first hand is worse; 1 if better; 0 if tied.
func (r *HandRank) Compare(other *HandRank) int {
	// Nil checks
	if r == nil && other != nil {
		return -1
	}
	if r != nil && other == nil {
		return 1
	}
	if r == nil && other == nil {
		return 0
	}

	// Check if one hand outclasses the other
	if r.RankType < other.RankType {
		return -1
	}
	if r.RankType > other.RankType {
		return 1
	}

	// Check if one hand is better in the same class
	for i := range r.Ranks {
		if r.Ranks[i] < other.Ranks[i] {
			return -1
		}
		if r.Ranks[i] > other.Ranks[i] {
			return 1
		}
	}

	// Both hands are the same
	return 0
}

// Converts an array of card.Card to int
func CardsToInt(cards []card.Card) []int {
	ints := make([]int, len(cards))
	for i, card_ := range cards {
		ints[i] = card_.CardNum()
	}
	return ints
}

// Returns true if all the cards are suited
func suited(cards []card.Card) bool {
	suit := cards[0].Suit()
	for _, card_ := range cards[1:] {
		if card_.Suit() != suit {
			return false
		}
	}
	return true
}

// Sorting algorithm optimized for 5 or less int arrays
// Currently insertion sort
func smallSort_VOLATILE(intCards []int) []int {
	// Takes the max and puts at the beginning of each array
	for i := range intCards {
		max := -1
		maxIdx := -1

		left := intCards[i:]
		for j, val := range left {
			if val > max {
				max = val
				maxIdx = j
			}
		}

		holder := intCards[i]
		intCards[i] = intCards[maxIdx+i]
		intCards[maxIdx+i] = holder
	}

	return intCards
}

// Uses the above sort for ranks
// If we copy the implementation, this will be faster as we will not copy
func sortRanks_VOLATILE(ranks []card.Rank) []card.Rank {
	ints := make([]int, len(ranks))
	for i, rank := range ranks {
		ints[i] = int(rank)
	}
	ints = smallSort_VOLATILE(ints)
	for i, val := range ints {
		ranks[i] = card.Rank(val)
	}
	return ranks
}

// Gets a list of the card ranks for a list of cards
func getCardRanks(cards []card.Card) []card.Rank {
	ranks := make([]card.Rank, len(cards))
	for i, card_ := range cards {
		ranks[i] = card_.Rank()
	}
	return ranks
}

// Returns true if you have a wheel straight
func checkWheel(cardRanks []card.Rank) bool {
	wheel := []card.Rank{
		card.Rank(12),
		card.Rank(3),
		card.Rank(2),
		card.Rank(1),
		card.Rank(0),
	}
	for i, card_ := range cardRanks {
		if wheel[i] != card_ {
			return false
		}
	}
	return true
}

// Returns true if the cards form a straight
// Input must be sorted
func straighted(cardRanks []card.Rank) bool {
	prev := cardRanks[0]
	for _, curr := range cardRanks[1:] {
		if curr+1 != prev {
			// Check special case A,5,4,3,2
			return checkWheel(cardRanks)
		}
		prev = curr
	}
	return true
}

// Returns a map of duplicate ranks
func rankDupes(cardRanks []card.Rank) map[card.Rank]int {
	rankMap := map[card.Rank]int{}
	for _, rank := range cardRanks {
		_, ok := rankMap[rank]
		if !ok {
			rankMap[rank] = 1
		} else {
			rankMap[rank] += 1
		}
	}

	// Remove ranks that are not duped
	for rank, repeats := range rankMap {
		if repeats == 1 {
			delete(rankMap, rank)
		}
	}
	return rankMap
}

// Finds unpaired cards, and returns in order
func findUnpaired(cardRanks []card.Rank, dupeRanks map[card.Rank]int) []card.Rank {
	unpaired := []card.Rank{}
	for _, rank := range cardRanks {
		if _, ok := dupeRanks[rank]; !ok {
			unpaired = append(unpaired, rank)
		}
	}
	return unpaired
}

// Given a hand of cards, get the rank of the hand
func RankHand(cards []card.Card) (*HandRank, error) {
	// Get an integer array representation of the cards
	isSuited := suited(cards)

	cardRanks := getCardRanks(cards)
	cardRanks = sortRanks_VOLATILE(cardRanks)

	isStraight := straighted(cardRanks)

	// Check for straight flush
	if isSuited && isStraight {
		straightTo := cardRanks[0]
		if checkWheel(cardRanks) {
			straightTo = card.Rank(3)
		}
		return &HandRank{
			RankType: StraightFlush,
			Ranks:    []card.Rank{straightTo},
		}, nil
	}
	// Check for flush
	if isSuited {
		return &HandRank{
			RankType: Flush,
			Ranks:    cardRanks,
		}, nil
	}
	// Check for straight
	if isStraight {
		straightTo := cardRanks[0]
		if checkWheel(cardRanks) {
			straightTo = card.Rank(3)
		}
		return &HandRank{
			RankType: Straight,
			Ranks:    []card.Rank{straightTo},
		}, nil
	}
	// We can check these out of order because a straight/flush
	// necessarily does not contain pair hands

	dupeRanks := rankDupes(cardRanks)

	// If there's more than one duplicated rank, it's a two pair or house
	if len(dupeRanks) == 2 {
		houseRank := card.Rank(-1)
		ranks := []card.Rank{}
		for rank, dupes := range dupeRanks {
			if dupes == 3 {
				houseRank = rank
			} else {
				ranks = append(ranks, rank)
			}
		}
		// Check for full house
		if houseRank >= 0 {
			return &HandRank{
				RankType: FullHouse,
				Ranks:    []card.Rank{houseRank, ranks[0]},
			}, nil
		}

		// Must be two pair
		ranks = sortRanks_VOLATILE(ranks)
		unpaired := findUnpaired(cardRanks, dupeRanks)
		return &HandRank{
			RankType: TwoPair,
			Ranks:    append(ranks, unpaired...),
		}, nil
	}
	for rank, dupes := range dupeRanks {
		unpaired := findUnpaired(cardRanks, dupeRanks)

		// Check for quads
		if dupes == 4 {
			return &HandRank{
				RankType: Quad,
				Ranks:    append([]card.Rank{rank}, unpaired...),
			}, nil
		}

		// Check for trips
		if dupes == 3 {
			return &HandRank{
				RankType: Triple,
				Ranks:    append([]card.Rank{rank}, unpaired...),
			}, nil
		}

		// Must be pair
		return &HandRank{
			RankType: Pair,
			Ranks:    append([]card.Rank{rank}, unpaired...),
		}, nil
	}
	// Return high card because the hand sucks
	return &HandRank{
		RankType: HighCard,
		Ranks:    cardRanks,
	}, nil
}

// A hand consists of a number of hole cards and a number of board cards.
type Hand interface {
	// Returns the best possible HandRank given a hand & table combination
	BestHand() (*HandRank, error)
}

// A rudimentary way of calculating the best hand on the board
type SimpleHandImpl struct {
	// The cards in hand (Play exactly 2/4 or 2/5)
	hole []card.Card

	// The cards on the table (Play exactly 3)
	community []card.Card
}

// Creates a new simple hand evaluator
func NewSimpleHandEvaluator(hole []card.Card, community []card.Card) *SimpleHandImpl {
	return &SimpleHandImpl{
		hole:      hole,
		community: community,
	}
}

func (h *SimpleHandImpl) BestHand() (*HandRank, error) {
	if len(h.hole) < 4 || len(h.hole) > 5 {
		return nil, fmt.Errorf("expected 4 or 5 hole cards, instead got: %d", len(h.hole))
	}
	if len(h.community) != 5 {
		return nil, fmt.Errorf("expected 5 community cards, instead got: %d", len(h.community))
	}

	var best *HandRank = nil

	// Iterate through all combinations of 2 hole cards and 3 community cards
	// and find the best hand from these combos
	for i, hole1 := range h.hole {
		for _, hole2 := range h.hole[i+1:] {
			for j, com1 := range h.community {
				for k, com2 := range h.community[j+1:] {
					for _, com3 := range h.community[k+j+2:] {
						combo := []card.Card{
							hole1,
							hole2,
							com1,
							com2,
							com3,
						}
						rank, err := RankHand(combo)
						if err != nil {
							return nil, errors.Wrapf(err, "RankHand")
						}

						if rank.Compare(best) > 0 {
							best = rank
						}
					}
				}
			}
		}
	}

	return best, nil
}

var _ Hand = (*SimpleHandImpl)(nil)
