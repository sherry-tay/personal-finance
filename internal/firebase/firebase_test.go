package firebase

import (
	"reflect"
	"testing"
)

func TestGetAveragePrice(t *testing.T) {
	tests := []struct {
		input []Stock
		expected map[string]float64
	}{
		{
			input: []Stock {
				{
					Code: "ABC",
					Price: 2.0,
					Volume: 100,
				},
				{
					Code: "DEF",
					Price: 5.0,
					Volume: 100,
				},
			},
			expected: map[string]float64 {
				"ABC": 2.0,
				"DEF": 5.0,
			},
		},
		{
			input: []Stock {
				{
					Code: "ABC",
					Price: 2.0,
					Volume: 100,
				},
				{
					Code: "ABC",
					Price: 4.0,
					Volume: 100,
				},
				{
					Code: "DEF",
					Price: 5.0,
					Volume: 100,
				},
				{
					Code: "ABC",
					Price: 10.0,
					Volume: 300,
				},
			},
			expected: map[string]float64 {
				"ABC": 7.2,
				"DEF": 5.0,
			},
		},
	}
	for _, test := range tests {
		t.Run("Test getAveragePrice", func(t *testing.T) {
			if actual := getAveragePrice(test.input); !reflect.DeepEqual(actual, test.expected) {
				t.Errorf("getAveragePrice() = %v, expected %v", actual, test.expected)
			}
		})
	}
}