package watchlist_test

import (
	"fmt"
	"strings"

	"github.com/acarl005/stripansi"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	. "ticker/internal/position"
	. "ticker/internal/quote"
	. "ticker/internal/ui/component/watchlist"
)

func removeFormatting(text string) string {
	return stripansi.Strip(text)
}

var _ = Describe("Watchlist", func() {
	describe := func(desc string) func(bool, bool, float64, Position, string) string {
		return func(isActive bool, isRegularTradingSession bool, change float64, position Position, expected string) string {
			return fmt.Sprintf("%s expected:%s", desc, expected)
		}
	}

	DescribeTable("should render a watchlist",
		func(isActive bool, isRegularTradingSession bool, change float64, position Position, expected string) {

			var positionMap map[string]Position
			if (position == Position{}) {
				positionMap = map[string]Position{}
			} else {
				positionMap = map[string]Position{
					"AAPL": position,
				}
			}

			m := NewModel()
			m.Width = 80
			m.Positions = positionMap
			m.Quotes = []Quote{
				{
					ResponseQuote: ResponseQuote{
						Symbol:    "AAPL",
						ShortName: "Apple Inc.",
					},
					Price:                   1.00 + change,
					Change:                  change,
					ChangePercent:           change,
					IsActive:                isActive,
					IsRegularTradingSession: isRegularTradingSession,
				},
			}
			Expect(removeFormatting(m.View())).To(Equal(expected))
		},
		Entry(
			describe("gain"),
			true,
			true,
			0.05,
			Position{},
			strings.Join([]string{
				"",
				"AAPL                       ⦿                                                1.05",
				"Apple Inc.                                                       ↑ 0.05  (0.05%)",
			}, "\n"),
		),
		Entry(
			describe("loss"),
			true,
			true,
			-0.05,
			Position{},
			strings.Join([]string{
				"",
				"AAPL                       ⦿                                                0.95",
				"Apple Inc.                                                      ↓ -0.05 (-0.05%)",
			}, "\n"),
		),
		Entry(
			describe("gain, after hours"),
			true,
			false,
			0.05,
			Position{},
			strings.Join([]string{
				"",
				"AAPL                       ⦾                                                1.05",
				"Apple Inc.                                                       ↑ 0.05  (0.05%)",
			}, "\n"),
		),
		Entry(
			describe("position, gain"),
			true,
			true,
			0.05,
			Position{
				AggregatedLot: AggregatedLot{
					Symbol:   "AAPL",
					Quantity: 100.0,
					Cost:     100.0,
				},
				Value:            105.0,
				DayChange:        5.0,
				DayChangePercent: 5.0,
			},
			strings.Join([]string{
				"",
				"AAPL                       ⦿                     105.00                     1.05",
				"Apple Inc.                              ↑ 5.00  (5.00%)          ↑ 0.05  (0.05%)",
			}, "\n"),
		),
		Entry(
			describe("position, loss"),
			true,
			true,
			-0.05,
			Position{
				AggregatedLot: AggregatedLot{
					Symbol:   "AAPL",
					Quantity: 100.0,
					Cost:     100.0,
				},
				Value:            95.0,
				DayChange:        -5.0,
				DayChangePercent: -5.0,
			},
			strings.Join([]string{
				"",
				"AAPL                       ⦿                      95.00                     0.95",
				"Apple Inc.                             ↓ -5.00 (-5.00%)         ↓ -0.05 (-0.05%)",
			}, "\n"),
		),
		Entry(
			describe("position, closed market"),
			false,
			false,
			0.0,
			Position{
				AggregatedLot: AggregatedLot{
					Symbol:   "AAPL",
					Quantity: 100.0,
					Cost:     100.0,
				},
				Value:            95.0,
				DayChange:        0.0,
				DayChangePercent: 0.0,
			},
			strings.Join([]string{
				"",
				"AAPL                                              95.00                     1.00",
				"Apple Inc.                                                         0.00  (0.00%)",
			}, "\n"),
		),
	)

	Context("when there are more than one symbols on the watchlist", func() {
		It("should render a watchlist with each symbol", func() {

			m := NewModel()
			m.Width = 80
			m.Quotes = []Quote{
				{
					ResponseQuote: ResponseQuote{
						Symbol:    "AAPL",
						ShortName: "Apple Inc.",
					},
					Price:                   1.05,
					Change:                  0.00,
					ChangePercent:           0.00,
					IsActive:                false,
					IsRegularTradingSession: false,
				},
				{
					ResponseQuote: ResponseQuote{
						Symbol:    "TW",
						ShortName: "ThoughtWorks",
					},
					Price:                   109.04,
					Change:                  3.53,
					ChangePercent:           5.65,
					IsActive:                true,
					IsRegularTradingSession: false,
				},
				{
					ResponseQuote: ResponseQuote{
						Symbol:    "GOOG",
						ShortName: "Google Inc.",
					},
					Price:                   2523.53,
					Change:                  -32.02,
					ChangePercent:           -1.35,
					IsActive:                true,
					IsRegularTradingSession: false,
				},
				{
					ResponseQuote: ResponseQuote{
						Symbol:    "BTC-USD",
						ShortName: "Bitcoin",
					},
					Price:                   50000.0,
					Change:                  10000.0,
					ChangePercent:           20.0,
					IsActive:                true,
					IsRegularTradingSession: true,
				},
			}
			expected := strings.Join([]string{
				"",
				"BTC-USD                    ⦿                                            50000.00",
				"Bitcoin                                                     ↑ 10000.00  (20.00%)",
				"TW                         ⦾                                              109.04",
				"ThoughtWorks                                                     ↑ 3.53  (5.65%)",
				"GOOG                       ⦾                                             2523.53",
				"Google Inc.                                                    ↓ -32.02 (-1.35%)",
				"AAPL                                                                        1.05",
				"Apple Inc.                                                         0.00  (0.00%)",
			}, "\n")
			Expect(removeFormatting(m.View())).To(Equal(expected))
		})
	})

	Context("when no quotes are set", func() {
		It("should render an empty watchlist", func() {
			m := NewModel()
			Expect(m.View()).To(Equal(""))
		})

	})
})