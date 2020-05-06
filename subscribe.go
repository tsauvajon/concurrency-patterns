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

func (s *sub) Close() error {
	errc := make(chan error)
	s.closing <- errc
	return <-errc
}

func (s *sub) Updates() <-chan Item {
	return s.updates
}

type fetchResult struct {
	fetched []Item
	next    time.Time
	err     error
}

func (s *sub) loop() {
	var (
		err         error
		pending     []Item
		next        time.Time
		alreadySeen = make(map[string]bool)
		fetchDone   chan fetchResult // if chan isn't nil, Fetch is running
	)
	for {
		var fetchDelay time.Duration
		if now := time.Now(); next.After(now) {
			fetchDelay = next.Sub(now)
		}
		var startFetching <-chan time.Time
		if fetchDone == nil { // not currently fetching
			startFetching = time.After(fetchDelay)
		}

		var first Item
		var updates chan Item
		if len(pending) > 0 {
			first = pending[0]
			updates = s.updates
		}

		select {
		case <-startFetching:
			fetchDone = make(chan fetchResult, 1)
			go func() {
				fetched, next, err := s.fetcher.Fetch()
				fetchDone <- fetchResult{fetched, next, err}
			}()
		case r := <-fetchDone:
			fetchDone = nil
			err = r.err
			next = r.next
			if err != nil {
				fmt.Println(err)
				time.Sleep(10 * time.Second)
				break
			}

			for _, item := range r.fetched {
				if !alreadySeen[item.GUID] {
					pending = append(pending, item)
					alreadySeen[item.GUID] = true
				}
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
