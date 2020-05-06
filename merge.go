package main

type merge struct {
	subs    []Subscription
	udpates chan Item
}

func Merge(subs ...Subscription) Subscription {
	m := merge{
		subs:    subs,
		udpates: make(chan Item),
	}
	return &m
}

func (m *merge) Close() error {
	closing := make(chan error)
	close(m.udpates)
	for _, s := range m.subs {
		s := s
		go func() {
			closing <- s.Close()
		}()
	}
	var err error
	for range m.subs {
		if e := <-closing; e != nil {
			err = e
		}
	}
	return err
}

func (m *merge) Updates() <-chan Item {
	go func() {
		for _, s := range m.subs {
			s := s
			go func() {
				uc := s.Updates()
				for {
					u := <-uc
					m.udpates <- u
				}
			}()
		}
	}()

	return m.udpates
}
