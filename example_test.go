package cbrates_test

import (
	"context"
	"fmt"
	"log"

	"github.com/itroot/cbrates"
)

func Example() {
	client, err := cbrates.NewCachedClient(log.Default())
	if err != nil {
		log.Fatal(err)
	}

	rates, err := client.GetDailyRates(context.Background(), 2022, 5, 20)
	if err != nil {
		log.Fatal(err)
	}

	for _, valute := range rates.ValuteSeq {
		if valute.CharCode == "USD" {
			fmt.Println(valute.MustParseFloat64())
			break
		}
	}
	// Output: 62.4031
}
