// Package cbrates implements client for the API of Central Bank of Russia, described here: http://www.cbr.ru/development/sxml/ .
package cbrates

import (
	"context"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kirsle/configdir"
)

const (
	DateFormat = "02.01.2006"
)

type ValCurs struct {
	Date      string    `xml:"Date,attr"`
	ValuteSeq []*Valute `xml:"Valute"`
}

// Filter returns only requested currencies' info
func (vc *ValCurs) Filter(codes []string) *ValCurs {
	if len(codes) == 0 {
		return vc
	}

	result := &ValCurs{
		Date:      vc.Date,
		ValuteSeq: make([]*Valute, 0, len(codes)),
	}
	for _, valute := range vc.ValuteSeq {
		for _, code := range codes {
			if code == valute.CharCode {
				result.ValuteSeq = append(result.ValuteSeq, valute)
				break
			}
		}
	}
	return result
}

type Valute struct {
	ID       string  `xml:"ID,attr"`
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Nominal  float64 `xml:"Nominal"`
	Name     string  `xml:"Name"`
	Value    string  `xml:"Value"`
}

// MustParseFloat64 returns currency value in RUB in float64.
func (v *Valute) MustParseFloat64() float64 {
	value, err := strconv.ParseFloat(strings.Join(strings.Split(v.Value, ","), "."), 64)
	if err != nil {
		panic(err)
	}
	return value
}

// Client is a interface for accessing the CBR API.
type Client interface {
	// GetDailyRates implements a client to the http://www.cbr.ru/development/sxml/ for a single day request
	GetDailyRates(ctx context.Context, year int, month time.Month, day int) (*ValCurs, error)
}

// NewClient creates a new client for the CBR API.
func NewClient(logger *log.Logger) (Client, error) {
	return &client{
		logger: logger,
	}, nil
}

// NewClient creates cache wrapper and a new client for the CBR API.
func NewCachedClient(logger *log.Logger) (Client, error) {
	cachePath := configdir.LocalCache("cbrates")
	if err := configdir.MakePath(cachePath); err != nil {
		return nil, err
	}
	client, err := NewClient(logger)
	if err != nil {
		return nil, err
	}
	return &cachedClient{
		logger: logger,
		client: client,

		cachePath: cachePath,
	}, err
}
