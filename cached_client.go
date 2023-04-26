package cbrates

import (
	"context"
	"encoding/gob"
	"log"
	"os"
	"path"
	"time"
)

type cachedClient struct {
	logger *log.Logger
	client Client

	cachePath string
}

func (c *cachedClient) GetDailyRates(ctx context.Context, year int, month time.Month, day int) (*ValCurs, error) {
	key := path.Join(c.cachePath, "GetDailyRates-"+formatDate(year, month, day))

	if result := func() *ValCurs {
		f, err := os.Open(key)
		if err != nil {
			if !os.IsNotExist(err) {
				c.logger.Println(err)
			}
			return nil
		}
		defer f.Close()
		var result ValCurs
		decoder := gob.NewDecoder(f)
		if err := decoder.Decode(&result); err != nil {
			c.logger.Println(err)
			return nil
		}
		return &result
	}(); result != nil {
		return result, nil
	}

	result, err := c.client.GetDailyRates(ctx, year, month, day)
	if err != nil {
		return nil, err
	}

	func() {
		f, err := os.CreateTemp(c.cachePath, "tmpcache")
		if err != nil {
			c.logger.Println(err)
			return
		}
		name := f.Name()
		defer f.Close()

		encoder := gob.NewEncoder(f)
		if err := encoder.Encode(result); err != nil {
			c.logger.Println(err)
			return
		}

		if err := f.Close(); err != nil {
			c.logger.Println(err)
			return
		}
		f = nil

		if err := os.Rename(name, key); err != nil {
			c.logger.Println(err)
			return
		}
	}()

	return result, nil
}
