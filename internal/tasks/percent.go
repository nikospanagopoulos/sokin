package tasks

// DayPercent returns the average progress (0-100) of one day's tasks.
// ok is false when the day has no tasks: such days are ignored.
func DayPercent(progress []int) (percent int, ok bool) {
	if len(progress) == 0 {
		return 0, false
	}
	sum := 0
	for _, p := range progress {
		sum += p
	}
	return (sum + len(progress)/2) / len(progress), true
}
