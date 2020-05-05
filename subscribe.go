package main

import (
	"fmt"
	"time"
)

type Subscription interface {
	Close() error
	Updates() <-chan Item
}

type sub struct {
	fetcher Fetcher
	updates chan Item
	closing chan chan error
}

func Subscribe(f Fetcher) Subscription {
	s := sub{
		fetcher: f,
		updates: make(chan Item),
		closing: make(chan chan error),
	}
	go s.loop()
	return &s
}

func Merge(subs ...Subscription) Subscription {
	return nil
}

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

func (s *sub) loop() {
	var (
		err     error
		pending []Item
		next    time.Time
	)
	for {
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		startFetching := time.After(fetchDelay)

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates
		}

		select {
		case <-startFetching:
			var fetched []Item
			fetched, next, err = s.fetcher.Fetch()
			pending = append(pending, fetched...)

			if err != nil {
				fmt.Println(err)
				time.Sleep(10 * time.Second)
			}
		case updates <- first:
			pending = pending[1:] // remove first
		case errc := <-s.closing:
			errc <- err
			close(s.updates)
			return
		}
	}
}
