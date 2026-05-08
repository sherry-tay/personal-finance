package telegram

import (
	"fmt"
	"testing"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sherry-tay/personal-finance/internal/firestore"
)

func TestGetPriceResponse(t *testing.T) {
	tests := []struct {
		sgx, yahoo *priceSource
		expected   string
	}{
		{
			sgx:      &priceSource{get: func(code string) (float64, string, error) { return 82, "USD", nil }},
			yahoo:    &priceSource{get: func(code string) (float64, string, error) { return 24.5, "SGD", nil }},
			expected: "USD82.000",
		},
		{
			sgx:      &priceSource{get: func(code string) (float64, string, error) { return 3.4, "GBP", nil }},
			yahoo:    &priceSource{get: func(code string) (float64, string, error) { return 0.0, "", fmt.Errorf("Yahoo error") }},
			expected: "GBP3.400",
		},
		{
			sgx:      &priceSource{get: func(code string) (float64, string, error) { return 0.0, "", fmt.Errorf("Sgx error") }},
			yahoo:    &priceSource{get: func(code string) (float64, string, error) { return 871.24122, "EUR", nil }},
			expected: "EUR871.241",
		},
		{
			sgx:      &priceSource{get: func(code string) (float64, string, error) { return 0.0, "", fmt.Errorf("Sgx error") }},
			yahoo:    &priceSource{get: func(code string) (float64, string, error) { return 0.0, "", fmt.Errorf("Yahoo error") }},
			expected: "Something went wrong while fetching price of FOO",
		},
	}
	for _, test := range tests {
		t.Run("Test getPriceResponse", func(t *testing.T) {
			if actualText, actualParseMode := routeCommands(mockTelegramMessage("/price", "FOO"), test.sgx, test.yahoo, nil, nil); actualText != test.expected || actualParseMode != "" {
				t.Errorf("getPriceResponse() = %v %v, expected  %v ''", actualText, actualParseMode, test.expected)
			}
		})
	}
}

func TestGetProfit(t *testing.T) {
	type args struct {
		current, average, volume float64
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
		{
			input: args{
				current: 50.22,
				average: 34.72,
				volume:  1.5,
			},
			expectedDiff:       15.5,
			expectedPercentage: 44.643,
			expectedProfit:     23.25,
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
		firestoreInput map[string]firestore.Stock
		sgxInput       map[string]struct {
			price    float64
			currency string
		}
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
					Volume: 87.34,
				},
			},
			sgxInput: map[string]struct {
				price    float64
				currency string
			}{
				"ABC": {
					price: 3.18,
					currency: "SGD",
				},
				"DEF": {
					price: 2.5,
					currency: "",
				},
			},
			expected: "```" +
				`
|     | NOW  |   %    |   +/-   |
|-----|------|--------|---------|
| ABC | 3.18 | 157.70 |  194.60 |
| DEF | 2.50 | -55.91 | -276.87 |
` + "```" +
				`
Current portfolio: 536.350
Total invested: 618.618
Capital gains: -82.268 (-13.299%)`,
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
| DEF | Vol |      1 |    87.34 |
+     +-----+--------+----------+
|     | Now |  2.500 |  218.350 |
+     +-----+--------+----------+
|     | Avg |  5.670 |  495.218 |
+     +-----+--------+----------+
|     | Dif | -3.170 | -276.868 |
+     +-----+--------+----------+
|     | %   |        |  -55.908 |
+-----+-----+--------+----------+
` + "```" +
				`
Current portfolio: 536.350
Total invested: 618.618
Capital gains: -82.268 (-13.299%)`,
		},
		{
			firestoreInput: map[string]firestore.Stock{
				"QWE": {
					Code:   "QWE",
					Price:  1.234,
					Volume: 100,
				},
				"RTY": {
					Code:   "RTY",
					Price:  5.67,
					Volume: 87,
				},
			},
			sgxInput: map[string]struct {
				price    float64
				currency string
			}{
				"QWE": {
					price: 3.18,
					currency: "SGD",
				},
				"RTY": {
					price: 2.5,
					currency: "USD",
				},
			},
			expected: "```" +
				`
|     | NOW  |   %    |   +/-   |
|-----|------|--------|---------|
| QWE | 3.18 | 157.70 |  194.60 |
| RTY | 3.33 | -41.21 | -203.29 |
` + "```" +
				`
Current portfolio: 608.000
Total invested: 616.690
Capital gains: -8.690 (-1.409%)`,
			expectedDetailed: "```" +
				`
+-----+-----+--------+----------+
|     |     |  PER   |  TOTAL   |
+-----+-----+--------+----------+
| QWE | Vol |      1 |      100 |
+     +-----+--------+----------+
|     | Now |  3.180 |  318.000 |
+     +-----+--------+----------+
|     | Avg |  1.234 |  123.400 |
+     +-----+--------+----------+
|     | Dif |  1.946 |  194.600 |
+     +-----+--------+----------+
|     | %   |        |  157.699 |
+-----+-----+--------+----------+
| RTY | Vol |      1 |       87 |
+     +-----+--------+----------+
|     | Now |  3.333 |  290.000 |
+     +-----+--------+----------+
|     | Avg |  5.670 |  493.290 |
+     +-----+--------+----------+
|     | Dif | -2.337 | -203.290 |
+     +-----+--------+----------+
|     | %   |        |  -41.211 |
+-----+-----+--------+----------+
` + "```" +
				`
Current portfolio: 608.000
Total invested: 616.690
Capital gains: -8.690 (-1.409%)`,
		},
		{
			firestoreInput: map[string]firestore.Stock{
				"INVALID": {
					Code:   "INVALID",
					Price:  2000,
					Volume: 1000,
				},
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
				"WRONG": {
					Code:   "WRONG",
					Price:  68,
					Volume: 13000,
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
				"OOPS": {
					Code:   "OOPS",
					Price:  240,
					Volume: 10000,
				},
			},
			sgxInput: map[string]struct {
				price    float64
				currency string
			}{
				"A00": {
					price:    78,
					currency: "",
				},
				"123": {
					price:    2.5,
					currency: "SGD",
				},
				"ZYX": {
					price:    220.2,
					currency: "",
				},
				"JKLM": {
					price:    0.18,
					currency: "SGD",
				},
			},
			expected: "```" +
				`
|      |  NOW   |   %    |   +/-    |
|------|--------|--------|----------|
|  123 |   2.50 |   0.00 |     0.00 |
| A00  |  78.00 |  85.71 | 36000.00 |
| JKLM |   0.18 | -64.00 |    -0.32 |
| ZYX  | 220.20 |  19.03 |  7040.00 |
` + "```" +
				`
List of unknown codes: INVALID, OOPS, WRONG

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
List of unknown codes: INVALID, OOPS, WRONG

Current portfolio: 122097.680
Total invested: 79058.000
Capital gains: 43039.680 (54.441%)`,
		},
	}

	for _, test := range tests {
		fs := &customFirestore{
			read: func() (map[string]firestore.Stock, error) { return test.firestoreInput, nil },
			add: func(id string, s firestore.Stock) error {
				t.Error("Should not call add")
				return nil
			},
		}
		sgx := &priceSource{
			get: func(code string) (float64, string, error) {
				if stock, ok := test.sgxInput[code]; ok {
					return stock.price, stock.currency, nil
				}
				return 0.0, "", fmt.Errorf("No input price")
			},
		}
		yahoo := &priceSource{
			get: func(code string) (float64, string, error) {
				return 0.0, "", fmt.Errorf("No input price")
			},
		}
		currencyConverter := func(currency string) (float64, error) {
			return 0.75, nil
		}
		t.Run("Test getStatistics", func(t *testing.T) {
			if actualText, actualParseMode := routeCommands(mockTelegramMessage("/summary", ""), sgx, yahoo, currencyConverter, fs); actualText != test.expected || actualParseMode != "Markdown" {
				t.Errorf("getStatistics(formatTable) = %v %v, expected  %v Markdown", actualText, actualParseMode, test.expected)
			}
			if actualText, actualParseMode := routeCommands(mockTelegramMessage("/detailed", ""), sgx, yahoo, currencyConverter, fs); actualText != test.expectedDetailed || actualParseMode != "Markdown" {
				t.Errorf("getStatistics(formatDetailedTable) = %v %v, expected  %v Markdown", actualText, actualParseMode, test.expectedDetailed)
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
			testName:     "correct format with volume as float for fractional shares",
			input:        "ABC 1.23 100.5 20200131 myBroker",
			expected:     "Successfully added holdings!",
			shouldRunAdd: true,
			expectedID:   "20200131-ABC",
			expectedStock: firestore.Stock{
				Code:     "ABC",
				Price:    1.23,
				Volume:   100.5,
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
			testName:     "wrong price format",
			input:        "ABC S1.23 100 20200131 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "wrong price format",
			input:        "ABC 1.23N 100 20200131 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "wrong price format",
			input:        "ABC $1.23 100 20200131 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "wrong price format",
			input:        "ABC SGD1.23 100 20200131 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "wrong volume format",
			input:        "ABC 1.23 10A0 20200131 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "wrong date format",
			input:        "ABC 1.23 100 D20200131 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "wrong date format",
			input:        "ABC 1.23 100 20200131D myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "invalid date",
			input:        "ABC 1.23 100 20200132 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "invalid date",
			input:        "ABC 1.23 100 20201331 myBroker",
			expected:     "Something went wrong while parsing the arguments... Please enter using the correct format.",
			shouldRunAdd: false,
		},
		{
			testName:     "not enough arguments",
			input:        "ABC 1.23 100 20200131",
			expected:     "Not enough arguments supplied! Aborting...",
			shouldRunAdd: false,
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
			if actualText, actualParseMode := routeCommands(mockTelegramMessage("/add", test.input), nil, nil, nil, fs); actualText != test.expected || actualParseMode != "" || argID != test.expectedID || argStock != test.expectedStock {
				t.Errorf("addHoldings() = %v %v %v %v, expected  %v '' %v %v", actualText, actualParseMode, argID, argStock, test.expected, test.expectedID, test.expectedStock)
			}
		})
	}
}
func TestDefault(t *testing.T) {
	commands := []string{"/help", "/default", "/anything", "/command", "/any", "/invalid", "/"}
	expectedText := `Use one of the commands available. 

To check the price, follow the format: "/price <ticker>" where ticker is the ticker found in either SGX or Yahoo e.g. /price AAPL.

To add, follow the format: "/add <code> <price> <volume> <yyyymmdd> <stored location>" e.g. /add ABC 1.23 100 20210101 mybroker.`

	for _, command := range commands {
		t.Run("Test default", func(t *testing.T) {
			if actualText, actualParseMode := routeCommands(mockTelegramMessage(command, ""), nil, nil, nil, nil); actualText != expectedText || actualParseMode != "" {
				t.Errorf("default = %v %v, expected  %v ''", actualText, actualParseMode, expectedText)
			}
		})
		t.Run("Test default with args", func(t *testing.T) {
			if actualText, actualParseMode := routeCommands(mockTelegramMessage(command, "args"), nil, nil, nil, nil); actualText != expectedText || actualParseMode != "" {
				t.Errorf("default = %v %v, expected  %v ''", actualText, actualParseMode, expectedText)
			}
		})
	}
}

func mockTelegramMessage(command, commandArgs string) *tgbotapi.Message {
	return &tgbotapi.Message{
		Entities: &[]tgbotapi.MessageEntity{
			{
				Type:   "bot_command",
				Offset: 0,
				Length: len(command),
			},
		},
		Text: command + " " + commandArgs,
	}
}
