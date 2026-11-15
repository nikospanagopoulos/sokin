package tasks

// Item is one task's contribution to a percentage.
type Item struct {
	Progress int // 0-100
	Weight   int // task difficulty, 1-10
}

// Percent returns the weighted average progress (0-100) of the given tasks.
// ok is false when there is nothing to measure: such days are ignored.
func Percent(items []Item) (percent int, ok bool) {
	sum, totalWeight := 0, 0
	for _, it := range items {
		sum += it.Progress * it.Weight
		totalWeight += it.Weight
	}
	if totalWeight == 0 {
		return 0, false
	}
	return (sum + totalWeight/2) / totalWeight, true
}
