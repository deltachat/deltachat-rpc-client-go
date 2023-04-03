package deltachat

type Reactions struct {
	ReactionsByContact map[ContactId][]string
	Reactions          map[string]int // Unique reactions and their count
}
