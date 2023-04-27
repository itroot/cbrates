package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/itroot/cbrates"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	logger := log.Default()
	if err := run(ctx, logger, os.Args); err != nil {
		logger.Fatal("failed to run main", err)
	}
}

func run(ctx context.Context, logger *log.Logger, args []string) error {
	var (
		dateString   string = time.Now().Format(cbrates.DateFormat)
		filterString string = "CNY,USD,EUR"
		noCache      bool
	)
	fs := flag.NewFlagSet(args[0], flag.ExitOnError)
	fs.StringVar(&dateString, "date", dateString, "Date to get rate values")
	fs.StringVar(&filterString, "filter", filterString, "Currencies codes to filter, empty string for all currencies")
	fs.BoolVar(&noCache, "no-cache", noCache, "Whether or not to cache requests")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	date, err := time.Parse(cbrates.DateFormat, dateString)
	if err != nil {
		return err
	}

	codes := make(map[string]bool)
	if len(filterString) != 0 {
		for _, code := range strings.Split(filterString, ",") {
			codes[code] = true
		}
	}

	newClient := cbrates.NewCachedClient
	if noCache {
		newClient = cbrates.NewClient
	}

	client, err := newClient(logger)
	if err != nil {
		logger.Fatal(err)
	}

	rates, err := client.GetDailyRates(context.Background(), date.Year(), date.Month(), date.Day())
	if err != nil {
		logger.Fatal(err)
	}

	tw := table.NewWriter()
	tw.SetTitle("Курсы валют ЦБ РФ на %s", date.Format(cbrates.DateFormat))
	tw.AppendHeader(table.Row{
		"Code",
		"Num",
		"Amount",
		"Name",
		"Value",
	})
	for _, valute := range rates.ValuteSeq {
		if len(codes) == 0 || codes[valute.CharCode] {
			tw.AppendRow(table.Row{
				valute.CharCode,
				valute.NumCode,
				valute.Nominal,
				valute.Name,
				valute.Value,
			})
		}
	}
	fmt.Println(tw.Render())
	return nil
}
