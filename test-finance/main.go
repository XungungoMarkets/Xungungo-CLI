package main

import (
	"fmt"
	"log"

	"github.com/piquette/finance-go/chart"
	"github.com/piquette/finance-go/datetime"
	"github.com/piquette/finance-go/quote"
)

func main() {
	// Test 1: Quote actual
	fmt.Println("=== Test Quote: NVDA ===")
	q, err := quote.Get("NVDA")
	if err != nil {
		log.Printf("Error quote: %v", err)
	} else {
		fmt.Printf("Symbol:  %s\n", q.Symbol)
		fmt.Printf("Price:   %.2f\n", q.RegularMarketPrice)
		fmt.Printf("Change:  %.2f (%.2f%%)\n", q.RegularMarketChange, q.RegularMarketChangePercent)
		fmt.Printf("Volume:  %d\n", q.RegularMarketVolume)
	}

	// Test 2: Quote AAPL
	fmt.Println("\n=== Test Quote: AAPL ===")
	q2, err := quote.Get("AAPL")
	if err != nil {
		log.Printf("Error quote: %v", err)
	} else {
		fmt.Printf("Symbol:  %s\n", q2.Symbol)
		fmt.Printf("Price:   %.2f\n", q2.RegularMarketPrice)
	}

	// Test 3: Histórico
	fmt.Println("\n=== Test History: NVDA (5 días) ===")
	params := &chart.Params{
		Symbol:   "NVDA",
		Interval: datetime.OneDay,
	}
	params.Start = datetime.FromUnix(1741392000)
	params.End = datetime.FromUnix(1741824000)

	iter := chart.Get(params)
	for iter.Next() {
		b := iter.Bar()
		open, _ := b.Open.Float64()
		close_, _ := b.Close.Float64()
		fmt.Printf("  Open: %.2f  Close: %.2f  Vol: %d\n", open, close_, b.Volume)
	}
	if iter.Err() != nil {
		log.Printf("Error history: %v", iter.Err())
	}

	fmt.Println("\n=== Done ===")
}
