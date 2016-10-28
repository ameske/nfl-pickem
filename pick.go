package nflpickem

import "errors"

// PickRetriever is the interface implemented by types that can retrieve Pick information
type PickRetriever interface {
	Picks(username string, year int, week int) (PickSet, error)
}

// Picker is the interface implemented by a type that can make/update picks
type Picker interface {
	MakePicks(PickSet) error
}

// A Pick represents a user's selection for a given game.
//
// The Pick can stand on its own since it contains embedded game information
type Pick struct {
	Game      Game `json:"game"`
	User      User `json:"user"`
	Selection Team `json:"selection"`
	Points    int  `json:"points"`
}

func (p Pick) Equal(other Pick) bool {
	return p.Game.Equal(other.Game) && p.User.Equal(other.User)
}

const (
	maxSevens = 1
	maxFives  = 2
	maxThrees = 5
)

// A PickSet represents the set of all picks for a user for a given Week
type PickSet []Pick

// IsLegal returns whether or not the current set of picks is considered legal.
//
// A PickSet is legal if it doesn't contain too many "special" point values,
// and contains picks for the same year, week, and user.
func (picks PickSet) IsLegal() bool {
	threes := 0
	fives := 0
	sevens := 0

	for _, p := range picks {
		switch p.Points {
		case 3:
			threes++
		case 5:
			fives++
		case 7:
			sevens++

		}

	}

	return threes <= maxThrees && fives <= maxFives && sevens <= maxSevens
}

var (
	ErrGameLocked       = errors.New("game has already started - pick locked")
	ErrUnknownSelection = errors.New("selection does not match a game in the given pick set")
)

type PickFilterFunc func(p Pick) bool

func (picks PickSet) Filter(f PickFilterFunc) PickSet {
	matches := make(PickSet, 0, len(picks))

	for _, p := range picks {
		if f(p) {
			matches = append(matches, p)
		}
	}

	return matches
}

// Merge compares the original PickSet to the new PickSet, updating any picks
// that differ.
func (picks PickSet) Merge(other PickSet) error {
	for i, o := range other {
		originalFound := false
		for _, p := range picks {
			if o.Equal(p) {
				picks[i] = p
			}
		}
		if !originalFound {
			return ErrUnknownSelection
		}
	}

	return nil
}
