//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

var suits = [...]Suit{Spade, Diamond, Club, Heart}

type Suit uint8

type Rank uint8

type Card struct {
	Suit
	Rank
}

const (
	minRank = Ace
	maxRank = King
)

const (
	Spade Suit = iota
	Diamond
	Club
	Heart
	Joker
)

const (
	_ Rank = iota
	Ace
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

// String prints a card object as string
func (c Card) String() string {
	if c.Suit == Joker {
		return c.Suit.String()
	}
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

// New returns a new instance of a Deck of cards
func New(opts ...func([]Card) []Card) []Card {
	var cards []Card
	for _, suit := range suits {
		for rank := minRank; rank <= maxRank; rank++ {
			cards = append(cards, Card{Suit: suit, Rank: rank})
		}
	}
	for _, opt := range opts {
		cards = opt(cards)
	}
	return cards
}

// Default sorting using the sort package
func DefaultSort(cards []Card) []Card {
	sort.Slice(cards, Less(cards))
	return cards
}

// Sort 
func Sort(less func(cards []Card) func(i, j int) bool) func([]Card) []Card {
	return func(c []Card) []Card {
		sort.Slice(c, less(c))
		return c
	}
}

// Less 
func Less(cards []Card) func(i, j int) bool {
	return func(i, j int) bool {
		return absRank(cards[i]) < absRank(cards[j])
	}
}

// Shuffle
func Shuffle(cards []Card) []Card {
	ret := make([]Card, len(cards))
	r := rand.New(rand.NewSource(time.Now().Unix()))
	perm := r.Perm(len(cards))
	for i, j := range perm {
		ret[i] = cards[j]
	}
	return ret
}

// Jokers
func Jokers(n int) func([]Card) []Card {
	return func(cards []Card) []Card {
		for i := 0; i < n; i++ {
			cards = append(cards, Card{
				Rank: Rank(i),
				Suit: Joker,
			})
		}
		return cards
	}
}

// Filter 
func Filter(f func(card Card) bool) func([]Card) []Card {
	return func(c []Card) []Card {
		var ret []Card
		for _, card := range c {
			if !f(card) {
				ret = append(ret, card)
			}
		}
		return ret
	}
}

// Deck
func Deck(n int) func([]Card) []Card {
	return func(c []Card) []Card {
		var ret []Card
		for i := 0; i < n; i++ {
			ret = append(ret, c...)
		}
		return ret
	}
}

// absRank
func absRank(c Card) int {
	return int(c.Suit)*int(maxRank) + int(c.Rank)
}
