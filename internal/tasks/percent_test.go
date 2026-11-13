package tasks

import "testing"

func TestDayPercent(t *testing.T) {
	tests := []struct {
		name string
		in   []int
		want int
		ok   bool
	}{
		{"no tasks", nil, 0, false},
		{"one task", []int{70}, 70, true},
		{"average", []int{100, 50, 0}, 50, true},
		{"rounds", []int{100, 0, 0}, 33, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := DayPercent(tc.in)
			if got != tc.want || ok != tc.ok {
				t.Errorf("DayPercent(%v) = %d, %v; want %d, %v", tc.in, got, ok, tc.want, tc.ok)
			}
		})
	}
}
