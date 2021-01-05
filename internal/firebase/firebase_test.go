package firebase

import (
	"reflect"
	"testing"
)

func TestGetAveragePrice(t *testing.T) {
	tests := []struct {
		input []Stock
		expected map[string]Stock
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
			expected: map[string]Stock {
				"ABC": Stock {
					Code: "ABC",
					Price: 2.0,
					Volume: 100,
				},
				"DEF": Stock {
					Code: "DEF",
					Price: 5.0,
					Volume: 100,
				},
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
			expected: map[string]Stock {
				"ABC": Stock {
					Code: "ABC",
					Price: 7.2,
					Volume: 500,
				},
				"DEF": Stock {
					Code: "DEF",
					Price: 5.0,
					Volume: 100,
				},
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