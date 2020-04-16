package deck

import (
	testing "testing"

	assert "github.com/stretchr/testify/assert"
)

func TestSuits(t *testing.T) {
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
