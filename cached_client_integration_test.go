package cbrates

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_CachedClient_IntegrationTest(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	ctx := context.Background()

	client, err := NewCachedClient(log.Default())
	require.NoError(t, err)

	// That's a little hack to simplify things...
	// Overall we need to do a clean setup here, as
	// we still creating default caching path.
	// Evenutally this test should be in the cbrates_test package
	client.(*cachedClient).cachePath = t.TempDir()

	cases := []struct {
		Date    string
		RateUSD float64
	}{
		{
			"18.01.2022",
			76.0404,
		},
		{
			"24.01.2022",
			76.6903,
		},
		{
			"23.02.2022",
			80.4194,
		},
		{
			"24.02.2022",
			80.4194,
		},
		{
			"03.03.2022",
			103.2487,
		},
		{
			"23.03.2022",
			104.0741,
		},
		{
			"28.03.2022",
			95.6618,
		},
		{
			"25.04.2022",
			73.5050,
		},
		{
			"19.05.2022",
			63.5643,
		},
		{
			"03.06.2022",
			61.5750,
		},
		{
			"08.06.2022",
			60.9565,
		},
		{
			"24.06.2022",
			53.3578,
		},
		{
			"09.08.2022",
			60.3164,
		},
		{
			"12.08.2022",
			60.6229,
		},
		{
			"15.08.2022",
			60.8993,
		},
		{
			"16.08.2022",
			61.3747,
		},
		{
			"26.08.2022",
			59.7699,
		},
		{
			"29.08.2022",
			60.0924,
		},
		{
			"02.09.2022",
			60.2370,
		},
		{
			"26.09.2022",
			58.1006,
		},
		{
			"27.09.2022",
			57.9990,
		},
		{
			"27.10.2022",
			61.4277,
		},
		{
			"22.11.2022",
			60.7379,
		},
		{
			"13.12.2022",
			62.7674,
		},
		{
			"14.12.2022",
			63.2120,
		},
		{
			"16.12.2022",
			64.3015,
		},
		{
			"29.12.2022",
			71.3261,
		},
		{
			"17.01.2023",
			68.2892,
		},
		{
			"23.01.2023",
			68.6656,
		},
		{
			"27.01.2023",
			69.1263,
		},
		{
			"31.01.2023",
			69.5927,
		},
		{
			"03.02.2023",
			70.0414,
		},
		{
			"14.02.2023",
			73.6307,
		},
		{
			"27.02.2023",
			74.7087,
		},
		{
			"28.02.2023",
			75.4323,
		},
		{
			"31.03.2023",
			77.0863,
		},
	}

	getUSDValue := func(t *testing.T, vc *ValCurs) float64 {
		t.Helper()

		usd := vc.Filter([]string{"USD"})
		require.Len(t, usd.ValuteSeq, 1)

		return usd.ValuteSeq[0].MustParseFloat64()
	}

	for index, c := range cases {
		date, err := time.Parse(DateFormat, c.Date)
		require.NoError(t, err)

		t.Run("usd-on-"+c.Date, func(t *testing.T) {
			entries, err := os.ReadDir(client.(*cachedClient).cachePath)
			require.NoError(t, err)
			require.Len(t, entries, index)

			result, err := client.GetDailyRates(ctx, date.Year(), date.Month(), date.Day())
			require.NoError(t, err)
			_ = result

			require.Equal(t, c.RateUSD, getUSDValue(t, result))
		})

		t.Run("usd-on-"+c.Date+"-cached", func(t *testing.T) {
			entries, err := os.ReadDir(client.(*cachedClient).cachePath)
			require.NoError(t, err)
			require.Len(t, entries, index+1)

			result, err := client.GetDailyRates(ctx, date.Year(), date.Month(), date.Day())
			require.NoError(t, err)
			_ = result

			require.Equal(t, c.RateUSD, getUSDValue(t, result))
		})

	}
}
