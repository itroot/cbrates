package cbrates

import (
	"context"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"golang.org/x/net/html/charset"
)

type client struct {
	httpClient http.Client

	logger *log.Logger
}

func (c *client) GetDailyRates(ctx context.Context, year int, month time.Month, day int) (*ValCurs, error) {
	const (
		UserAgent = "cbrates/v0 (+https://github.com/itroot/cbrates)" // by some reason default go ua was getting blocked
	)

	date := formatDate(year, month, day)
	url := fmt.Sprintf("http://www.cbr.ru/scripts/XML_daily.asp?date_req=%s", date)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", UserAgent)

	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result ValCurs
	decoder := xml.NewDecoder(resp.Body)
	decoder.CharsetReader = charset.NewReaderLabel
	err = decoder.Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.ValuteSeq) == 0 {
		return nil, errors.New("no data for this date")
	}

	if result.Date != date {
		return nil, errors.New("response does not match request date")
	}

	return &result, nil
}
