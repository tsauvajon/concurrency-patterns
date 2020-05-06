package main

type merge struct {
	subs    []Subscription
	udpates chan Item
	closing chan struct{}
	errs    chan error
}

func Merge(subs ...Subscription) Subscription {
	m := merge{
		subs:    subs,
		udpates: make(chan Item),
		closing: make(chan struct{}),
		errs:    make(chan error),
	}

	for _, sub := range subs {
		sub := sub
		go func() {
			uc := sub.Updates()
			for {
				var item Item
				select {
				case <-m.closing:
					m.errs <- sub.Close()
					return
				case item = <-uc:
				}
				select {
				case <-m.closing:
					m.errs <- sub.Close()
					return
				case m.udpates <- item:

				}
			}
		}()
	}

	return &m
}

func (m *merge) Close() error {
	for range m.subs {
		m.closing <- struct{}{}
	}
	close(m.closing)
	var err error
	for range m.subs {
		if e := <-m.errs; e != nil {
			err = e
		}
	}
	close(m.udpates)
	return err
}

func (m *merge) Updates() <-chan Item {
	return m.udpates
}
