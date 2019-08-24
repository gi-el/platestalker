package alpr

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"golang.org/x/net/html"
)

type Scraper struct {
	Address string
	PollInterval time.Duration

	lastTimestampNano int64

	plateChan chan *ALPRGroup
	errChan chan error
}

func NewScraper(address string, interval time.Duration) (*Scraper, error) {
	s := &Scraper{
		Address: address,
		PollInterval: interval,
		plateChan: make(chan *ALPRGroup, CACHE_SIZE),
		errChan: make(chan error, 1),
	}

	return s, nil
}

// Blocks until a new plate is retrieved
// Internally polls until we find a new plate.
func (s *Scraper) Next() (*ALPRGroup, error) {
	select {
	case plate := <-s.plateChan:
		return plate, nil
	case err := <-s.errChan:
		return nil, err
	}
}

func (s *Scraper) Run(ctx context.Context) {
	timer := time.NewTimer(0) // Fire immediately, set poll interval after
	defer timer.Stop()

	log.Info().Msg("starting scraper")

	running := true

	for running {
		select {
		case <-timer.C:
			err := s.scrape()
			if err != nil {
				s.errChan <- err
			}
			timer.Reset(s.PollInterval)
		case <-ctx.Done():
			running = false
			s.errChan <- io.EOF
			break
		}
	}

	log.Info().Msg("scraper stopped")
}

func (s *Scraper) scrape() error {
	resp, err := http.Get(s.Address)
	if err != nil {
		return fmt.Errorf("failed to get index page: %s", err)
	}
	defer resp.Body.Close()

	z := html.NewTokenizer(resp.Body)
	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				break
			}
			return fmt.Errorf("html parse error: %s", z.Err())
		}

		if tt != html.StartTagToken {
			continue
		}

		tn, _ := z.TagName()
		if len(tn) != 1 || tn[0] != 'a' {
			continue
		}

		_, val, _ := z.TagAttr()
		valStr := string(val)
		if !strings.HasPrefix(valStr, "/meta/") {
			continue
		}

		meta := strings.TrimPrefix(valStr, "/meta/")
		elems := strings.Split(meta, "-")
		ts, err := strconv.ParseInt(elems[2], 10, 64)
		if err != nil {
			return fmt.Errorf("failed parsing timestamp: %s", err)
		}

		if ts <= s.lastTimestampNano {
			continue
		}
		s.lastTimestampNano = ts

		group, err := s.getGroup(meta)
		if err != nil {
			return err
		}
		s.plateChan <- group
	}

	return nil
}

func (s *Scraper) getGroup(id string) (*ALPRGroup, error) {
	resp, err := http.Get(s.groupAddress(id))
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve plate meta data: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read plate meta data: %s", err)
	}
	group := &ALPRGroup{}
	err = json.Unmarshal(body, group)
	if err != nil {
		return nil, fmt.Errorf("failed to parse plate meta data: %s", err)
	}

	if group.DataType != ALPR_GROUP_DATA_TYPE || group.Version != ALPR_GROUP_VERSION {
		return nil, fmt.Errorf("invalid meta data version %s/%d", group.DataType, group.Version)
	}

	return group, nil
}

func (s *Scraper) groupAddress(id string) string {
	return fmt.Sprintf("%s/meta/%s", s.Address, id)
}
