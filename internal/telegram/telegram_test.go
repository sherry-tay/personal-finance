package telegram

import (
	"testing"
	"time"

	"github.com/sherry-tay/personal-finance/internal/firestore"
)

func TestGetProfit(t *testing.T) {
	type args struct {
		current, average float64
		volume           int
	}
	tests := []struct {
		input                                            args
		expectedDiff, expectedPercentage, expectedProfit float64
	}{
		{
			input: args{
				current: 24.67,
				average: 21.002,
				volume:  10000,
			},
			expectedDiff:       3.668,
			expectedPercentage: 17.465,
			expectedProfit:     36680.0,
		},
		{
			input: args{
				current: 65.209,
				average: 81.823902,
				volume:  89,
			},
			expectedDiff:       -16.615,
			expectedPercentage: -20.306,
			expectedProfit:     -1478.726,
		},
	}
	for _, test := range tests {
		t.Run("Test getProfit", func(t *testing.T) {
			if actualDiff, actualPercentage, actualProfit := getProfit(test.input.current, test.input.average, test.input.volume); actualDiff != test.expectedDiff || actualPercentage != test.expectedPercentage || actualProfit != test.expectedProfit {
				t.Errorf("getProfit() = %v %v %v, expected  %v %v %v", actualDiff, actualPercentage, actualProfit, test.expectedDiff, test.expectedPercentage, test.expectedProfit)
			}
		})
	}
}

func TestGetStatistics(t *testing.T) {
	tests := []struct {
		firestoreInput             map[string]firestore.Stock
		sgxInput                   map[string]float64
		expected, expectedDetailed string
	}{
		{
			firestoreInput: map[string]firestore.Stock{
				"ABC": {
					Code:   "ABC",
					Price:  1.234,
					Volume: 100,
				},
				"DEF": {
					Code:   "DEF",
					Price:  5.67,
					Volume: 87,
				},
			},
			sgxInput: map[string]float64{
				"ABC": 3.18,
				"DEF": 2.5,
			},
			expected: "```" +
				`
| CODE |  NOW  |  AVG  |  DIFF  |    %    | VOL |  TOTAL   |
|------|-------|-------|--------|---------|-----|----------|
| ABC  | 3.180 | 1.234 |  1.946 | 157.699 | 100 |  194.600 |
| DEF  | 2.500 | 5.670 | -3.170 | -55.908 |  87 | -275.790 |
` + "```" +
				`
Current portfolio: 535.500
Total invested: 616.690
Capital gains: -81.190 (-13.165%)`,
			expectedDetailed: "```" +
				`
+-----+-----+--------+----------+
|     |     |  PER   |  TOTAL   |
+-----+-----+--------+----------+
| ABC | Vol |      1 |      100 |
+     +-----+--------+----------+
|     | Now |  3.180 |  318.000 |
+     +-----+--------+----------+
|     | Avg |  1.234 |  123.400 |
+     +-----+--------+----------+
|     | Dif |  1.946 |  194.600 |
+     +-----+--------+----------+
|     | %   |        |  157.699 |
+-----+-----+--------+----------+
| DEF | Vol |      1 |       87 |
+     +-----+--------+----------+
|     | Now |  2.500 |  217.500 |
+     +-----+--------+----------+
|     | Avg |  5.670 |  493.290 |
+     +-----+--------+----------+
|     | Dif | -3.170 | -275.790 |
+     +-----+--------+----------+
|     | %   |        |  -55.908 |
+-----+-----+--------+----------+
` + "```" +
				`
Current portfolio: 535.500
Total invested: 616.690
Capital gains: -81.190 (-13.165%)`,
		},
		{
			firestoreInput: map[string]firestore.Stock{
				"A00": {
					Code:   "A00",
					Price:  42,
					Volume: 1000,
				},
				"ZYX": {
					Code:   "ZYX",
					Price:  185,
					Volume: 200,
				},
				"123": {
					Code:   "123",
					Price:  2.5,
					Volume: 23,
				},
				"JKLM": {
					Code:   "JKLM",
					Price:  0.5,
					Volume: 1,
				},
			},
			sgxInput: map[string]float64{
				"A00":  78,
				"123": 2.5,
				"ZYX":  220.2,
				"JKLM": 0.18,
			},
			expected: "```" +
				`
| CODE |   NOW   |   AVG   |  DIFF  |    %    | VOL  |   TOTAL   |
|------|---------|---------|--------|---------|------|-----------|
|  123 |   2.500 |   2.500 |  0.000 |   0.000 |   23 |     0.000 |
| A00  |  78.000 |  42.000 | 36.000 |  85.714 | 1000 | 36000.000 |
| JKLM |   0.180 |   0.500 | -0.320 | -64.000 |    1 |    -0.320 |
| ZYX  | 220.200 | 185.000 | 35.200 |  19.027 |  200 |  7040.000 |
` + "```" +
				`
Current portfolio: 122097.680
Total invested: 79058.000
Capital gains: 43039.680 (54.441%)`,
			expectedDetailed: "```" +
				`
+------+-----+---------+-----------+
|      |     |   PER   |   TOTAL   |
+------+-----+---------+-----------+
|  123 | Vol |       1 |        23 |
+      +-----+---------+-----------+
|      | Now |   2.500 |    57.500 |
+      +-----+         +           +
|      | Avg |         |           |
+      +-----+---------+-----------+
|      | Dif |   0.000 |     0.000 |
+      +-----+---------+           +
|      | %   |         |           |
+------+-----+---------+-----------+
| A00  | Vol |       1 |      1000 |
+      +-----+---------+-----------+
|      | Now |  78.000 | 78000.000 |
+      +-----+---------+-----------+
|      | Avg |  42.000 | 42000.000 |
+      +-----+---------+-----------+
|      | Dif |  36.000 | 36000.000 |
+      +-----+---------+-----------+
|      | %   |         |    85.714 |
+------+-----+---------+-----------+
| JKLM | Vol |       1 |         1 |
+      +-----+---------+-----------+
|      | Now |   0.180 |     0.180 |
+      +-----+---------+-----------+
|      | Avg |   0.500 |     0.500 |
+      +-----+---------+-----------+
|      | Dif |  -0.320 |    -0.320 |
+      +-----+---------+-----------+
|      | %   |         |   -64.000 |
+------+-----+---------+-----------+
| ZYX  | Vol |       1 |       200 |
+      +-----+---------+-----------+
|      | Now | 220.200 | 44040.000 |
+      +-----+---------+-----------+
|      | Avg | 185.000 | 37000.000 |
+      +-----+---------+-----------+
|      | Dif |  35.200 |  7040.000 |
+      +-----+---------+-----------+
|      | %   |         |    19.027 |
+------+-----+---------+-----------+
` + "```" +
				`
Current portfolio: 122097.680
Total invested: 79058.000
Capital gains: 43039.680 (54.441%)`,
		},
	}

	for _, test := range tests {
		fs := &customFirestore{
			read: func() (map[string]firestore.Stock, error) { return test.firestoreInput, nil },
			add:  func(id string, s firestore.Stock) error { 
				t.Error("Should not call add") 
				return nil 
			},
		}
		sgx := &customSgx{
			get: func(code string) (float64, error) { return test.sgxInput[code], nil },
		}
		t.Run("Test getStatistics", func(t *testing.T) {
			if actual := fs.getStatistics(sgx, formatTable); actual != test.expected {
				t.Errorf("getStatistics(formatTable) = %v, expected  %v", actual, test.expected)
			}
			if actual := fs.getStatistics(sgx, formatDetailedTable); actual != test.expectedDetailed {
				t.Errorf("getStatistics(formatDetailedTable) = %v, expected  %v", actual, test.expectedDetailed)
			}
		})
	}
}

func TestAddHoldings(t *testing.T) {
	tests := []struct {
		testName, input, expected, expectedID string
		shouldRunAdd                          bool
		expectedStock                         firestore.Stock
	}{
		{
			testName:     "correct format",
			input:        "ABC 1.23 100 20200131 myBroker",
			expected:     "Successfully added holdings!",
			shouldRunAdd: true,
			expectedID:   "20200131-ABC",
			expectedStock: firestore.Stock{
				Code:     "ABC",
				Price:    1.23,
				Volume:   100,
				Date:     time.Date(2020, time.January, 31, 0, 0, 0, 0, time.UTC),
				StoredIn: "myBroker",
			},
		},
		{
			testName:     "correct format with more than one word as stored location",
			input:        "ABC 1.23 100 20200131 my awesome broker",
			expected:     "Successfully added holdings!",
			shouldRunAdd: true,
			expectedID:   "20200131-ABC",
			expectedStock: firestore.Stock{
				Code:     "ABC",
				Price:    1.23,
				Volume:   100,
				Date:     time.Date(2020, time.January, 31, 0, 0, 0, 0, time.UTC),
				StoredIn: "my awesome broker",
			},
		},
		{
			testName:      "wrong price format",
			input:         "ABC S1.23 100 20200131 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "wrong price format",
			input:         "ABC 1.23N 100 20200131 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "wrong price format",
			input:         "ABC $1.23 100 20200131 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "wrong price format",
			input:         "ABC SGD1.23 100 20200131 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "wrong volume format",
			input:         "ABC 1.23 10A0 20200131 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "wrong date format",
			input:         "ABC 1.23 100 D20200131 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "wrong date format",
			input:         "ABC 1.23 100 20200131D myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "invalid date",
			input:         "ABC 1.23 100 20200132 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "invalid date",
			input:         "ABC 1.23 100 20201331 myBroker",
			expected:      "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd:  false,
		},
		{
			testName:      "not enough arguments",
			input:         "ABC 1.23 100 20200131",
			expected:      "Not enough arguments supplied! Aborting...",
			shouldRunAdd:  false,
		},
	}

	for _, test := range tests {
		var argID string
		var argStock firestore.Stock
		fs := &customFirestore{
			read: func() (map[string]firestore.Stock, error) {
				t.Error("Should not call read")
				return map[string]firestore.Stock{}, nil
			},
			add: func(id string, s firestore.Stock) error {
				argID = id
				argStock = s
				return nil
			},
		}
		t.Run("Test addHoldings - "+test.testName, func(t *testing.T) {
			if actual := fs.addHoldings(test.input); actual != test.expected || argID != test.expectedID || argStock != test.expectedStock {
				t.Errorf("addHoldings() = %v %v %v, expected  %v %v %v", actual, argID, argStock, test.expected, test.expectedID, test.expectedStock)
			}
		})
	}
}
