package tasks

import "testing"

func TestPercent(t *testing.T) {
	tests := []struct {
		name string
		in   []Item
		want int
		ok   bool
	}{
		{"no tasks", nil, 0, false},
		{"one task", []Item{{70, 5}}, 70, true},
		{"equal weights", []Item{{100, 5}, {50, 5}, {0, 5}}, 50, true},
		{"hard task done", []Item{{100, 9}, {0, 1}}, 90, true},
		{"hard task not done", []Item{{100, 1}, {0, 9}}, 10, true},
		{"rounds", []Item{{100, 1}, {0, 2}}, 33, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := Percent(tc.in)
			if got != tc.want || ok != tc.ok {
				t.Errorf("Percent(%v) = %d, %v; want %d, %v", tc.in, got, ok, tc.want, tc.ok)
			}
		})
	}
}
