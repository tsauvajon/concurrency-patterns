package main

import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

type Item struct {
	Channel, Title, GUID string // a subset of RSS fields
}

type Fetcher interface {
	Fetch() (items []Item, next time.Time, err error)
}

var exampleItems = []struct{ title, guid string }{
	{"Smoggy", "1"},
	{"Suppose", "2"},
	{"Yoke", "3"},
	{"Brainy", "4"},
	{"Manage", "5"},
	{"Optimal", "6"},
	{"Price", "7"},
	{"Giants", "8"},
	{"Value", "9"},
	{"Tiresome", "10"},
	{"Ill", "11"},
	{"Fix", "12"},
	{"Informed", "13"},
	{"Defective", "14"},
	{"White", "15"},
	{"Fearless", "16"},
}

type RSSFetcher struct {
	domain string
}

func Fetch(domain string) Fetcher {
	return RSSFetcher{
		domain: domain,
	}
}

func (r RSSFetcher) Fetch() (items []Item, next time.Time, err error) {
	nbItems := rand.Intn(2) + 1

	for i := 0; i < nbItems; i++ {
		i := exampleItems[rand.Intn(len(exampleItems)-1)]
		items = append(items, Item{r.domain, i.title, i.guid})
	}

	next = time.Now().Add(time.Duration(rand.Intn(3000)+1000) * time.Millisecond)

	return
}
