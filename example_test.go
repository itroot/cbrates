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

	fmt.Println(rates.Filter([]string{"USD"}).ValuteSeq[0].MustParseFloat64())
	// Output: 62.4031
}
