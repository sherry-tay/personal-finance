package main

import (
	"testing"
)

func TestGetProfit(t *testing.T) {
	type args struct {
		current, average float64
		volume int
	}
	tests := []struct {
		input args
		expectedDiff, expectedProfit float64
	}{
		{
			input: args {
				current: 24.67,
				average: 21.002,
				volume: 10000,
			},
			expectedDiff: 3.668,
			expectedProfit: 36680.0,
		},
		{
			input: args {
				current: 65.209,
				average: 81.823902,
				volume: 89,
			},
			expectedDiff: -16.615,
			expectedProfit: -1478.726,
		},
	}
	for _, test := range tests {
		t.Run("Test getProfit", func(t *testing.T) {
			if actualDiff, actualProfit := getProfit(test.input.current, test.input.average, test.input.volume); actualDiff != test.expectedDiff || actualProfit != test.expectedProfit {
				t.Errorf("getProfit() = %v %v, expected  %v %v", actualDiff, actualProfit, test.expectedDiff, test.expectedProfit)
			}
		})
	}
}