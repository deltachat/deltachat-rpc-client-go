package deltachat

type Reactions struct {
	ReactionsByContact map[uint64][]string
	Reactions          map[string]int // Unique reactions and their count
}
