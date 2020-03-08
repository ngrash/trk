package main

import (
	"testing"
	"time"
)

func d(s string) *time.Duration {
	dur, err := time.ParseDuration(s)
	if err != nil {
		panic(err)
	}
	return &dur
}

func TestParseEntry(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in     string
		out    *time.Duration
		hasErr bool
	}{
		{"20-03-02 5h55m", d("5h55m"), false},
		{"20-03-02 -8h", d("-8h"), false},
		{"20-03-02 1000-1600", d("6h"), false},
		{"20-03-02 1000-", nil, false},
		{"20-03-02 -1000", nil, true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.in, func(t *testing.T) {
			t.Parallel()
			e, err := parseEntry(tt.in)
			if tt.hasErr != (err != nil) {
				t.Errorf("err wanted: %v, got: %v\n", tt.hasErr, err)
			}

			if err == nil {
				if tt.out == nil && e.duration != nil {
					t.Errorf("no duration wanted, got: %v\n", e.duration)
				} else if e.duration != nil && *e.duration != *tt.out {
					t.Errorf("duration wanted: %v, got: %v\n", tt.out, e.duration)
				}
			}
		})
	}
}
